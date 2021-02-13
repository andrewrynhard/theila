package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	caps "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metal "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"golang.org/x/net/websocket"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/scale/scheme"
	toolscache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"github.com/talos-systems/theila/pkg/client"
	"github.com/talos-systems/theila/pkg/render"
)

type contextKeyServer struct{}
type contextKeyServerBinding struct{}

var (
	// ContextKeyServer is a type safe representation of the key "server" inside of a context.Context
	ContextKeyServer = contextKeyServer{}
	// ContextKeyServerBinding is a type safe representation of the key "serverbinding" inside of a context.Context
	ContextKeyServerBinding = contextKeyServerBinding{}
)

type ServerEvent struct {
	Server *metal.Server
	Action string
	Target string
}

// ServersController implements Controller functionality for the Server model
type ServersController struct {
	sync.Mutex
	*render.Renderer

	Client *client.Client
	Events chan ServerEvent

	servers []*metal.Server
}

// Context is a middleware that parses the Server ID from a request, loads the
// corresponding Server model, and makes it available as part of the request's
// context
func (s *ServersController) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx           context.Context
			server        *metal.Server
			serverbinding *caps.ServerBinding
			err           error
		)

		if id := chi.URLParam(r, "id"); id != "" {
			server, err = s.serverFromIDString(id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Fatal error: %v", err), http.StatusInternalServerError)

				return
			}

			ctx = context.WithValue(r.Context(), ContextKeyServer, server)

			if server.Status.InUse {
				serverbinding, err = s.serverbindingForServer(id)

				if err != nil {
					http.Error(w, fmt.Sprintf("Fatal error: %v", err), http.StatusInternalServerError)

					return
				}

				ctx = context.WithValue(ctx, ContextKeyServerBinding, serverbinding)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *ServersController) serverFromIDString(id string) (*metal.Server, error) {
	log.Println("Getting server")

	var server *metal.Server

	for _, srvr := range s.Servers() {
		if srvr.Name == id {
			server = srvr

			break
		}
	}

	if server == nil {
		return nil, fmt.Errorf("server not found")
	}

	log.Printf("Found server %q", id)

	return server, nil
}

func (s *ServersController) serverbindingForServer(id string) (*caps.ServerBinding, error) {
	log.Println("Getting serverbinding")

	serverbinding := &caps.ServerBinding{}

	if err := s.Client.Kubernetes.Get(context.TODO(), types.NamespacedName{Name: id}, serverbinding); err != nil {
		return nil, err
	}

	log.Printf("Found serverbinding %q", id)

	return serverbinding, nil
}

// Index shows a list of all Servers.
func (s *ServersController) Index(w http.ResponseWriter, r *http.Request) {
	log.Println("Listing servers")

	content := map[string]interface{}{"Servers": s.Servers(), "Current": "Servers"}

	s.RenderFull(w, http.StatusOK, "list servers", content)
}

// Get retrieves a single Server by its ID.
func (s *ServersController) Get(w http.ResponseWriter, r *http.Request) {
	var (
		server        *metal.Server
		serverbinding *caps.ServerBinding
	)

	switch v := r.Context().Value(ContextKeyServer).(type) {
	case nil:
	case *metal.Server:
		server = v
	default:
		log.Println("type unknown")
		return
	}

	switch v := r.Context().Value(ContextKeyServerBinding).(type) {
	case nil:
	case *caps.ServerBinding:
		serverbinding = v
	default:
		log.Println("type unknown")
		return
	}

	content := map[string]interface{}{
		"Server":        server,
		"ServerBinding": serverbinding,
		"Current":       "Servers",
	}

	s.RenderFull(w, http.StatusOK, "get server", content)
}

// Summary shows a summary of all Servers.
func (s *ServersController) Summary(w http.ResponseWriter, r *http.Request) {
	count := 0

	servers := s.Servers()

	for _, server := range servers {
		if server.Status.InUse {
			count++
		}
	}

	capacity := 100.0 * float64(count) / float64(len(servers))

	content := map[string]interface{}{
		"TotalServers":          len(servers),
		"TotalServersAllocated": count,
		"ServerCapacity":        fmt.Sprintf("%.2f%%", capacity),
	}

	s.Renderer.RenderPartial(w, http.StatusOK, "servers summary", content)
}

func (s *ServersController) Servers() []*metal.Server {
	s.Lock()
	defer s.Unlock()

	return s.servers
}

func (s *ServersController) Watch() error {
	skheme := runtime.NewScheme()
	_ = scheme.AddToScheme(skheme)
	_ = metal.AddToScheme(skheme)

	config, err := s.Client.KubernetesClientConfig.ClientConfig()
	if err != nil {
		return err
	}
	cache, err := cache.New(config, cache.Options{Scheme: skheme})
	if err != nil {
		return err
	}

	informer, err := cache.GetInformer(context.TODO(), &metal.Server{})
	if err != nil {
		return err
	}

	informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s.Lock()
			defer s.Unlock()

			server := obj.(*metal.Server)

			s.servers = append(s.servers, server)

			log.Printf("Added server %q", server.Name)

			s.Events <- ServerEvent{Server: server, Action: "append", Target: "server-list"}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			s.Lock()
			defer s.Unlock()

			server := newObj.(*metal.Server)

			for i, old := range s.servers {
				if old.UID == server.UID {
					s.servers[i] = server

					log.Printf("Updated server %q", server.Name)

					s.Events <- ServerEvent{Server: server, Action: "replace", Target: fmt.Sprintf("server-%s", server.Name)}

					break
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			s.Lock()
			defer s.Unlock()

			server := obj.(*metal.Server)

			for i, old := range s.servers {
				if old.UID == server.UID {
					s.servers[i] = s.servers[len(s.servers)-1]
					s.servers[len(s.servers)-1] = nil
					s.servers = s.servers[:len(s.servers)-1]

					log.Printf("Deleted server %q", server.Name)

					s.Events <- ServerEvent{Server: server, Action: "remove", Target: fmt.Sprintf("server-%s", server.Name)}

					break
				}
			}
		},
	})

	ctx := context.Background()

	go cache.Start(ctx)

	if ok := cache.WaitForCacheSync(ctx); ok {
		log.Println("Server cache synced.")
	}

	return nil
}

func (s *ServersController) Socket(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		log.Println("Streaming server updates")

		for {
			var (
				content map[string]interface{}
				err     error
			)

			select {
			case event := <-s.Events:
				content = map[string]interface{}{
					"Server": event.Server,
					"Action": event.Action,
					"Target": event.Target,
				}

				log.Printf("Performing %q action on %q target", event.Action, event.Target)
			}

			msg, err := s.Render("server list entry stream", content)
			if err != nil {
				log.Println(err)
			}

			websocket.Message.Send(ws, msg)
		}
	}).ServeHTTP(w, r)
}
