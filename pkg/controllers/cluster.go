package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"golang.org/x/net/websocket"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/scale/scheme"
	toolscache "k8s.io/client-go/tools/cache"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"github.com/talos-systems/theila/pkg/client"
	"github.com/talos-systems/theila/pkg/render"
)

type contextKeyCluster struct{}
type contextKeyClusterBinding struct{}

var (
	// ContextKeyCluster is a type safe representation of the key "cluster" inside of a context.Context
	ContextKeyCluster = contextKeyCluster{}
	// ContextKeyClusterBinding is a type safe representation of the key "clusterbinding" inside of a context.Context
	ContextKeyClusterBinding = contextKeyClusterBinding{}
)

type ClusterEvent struct {
	Cluster *capi.Cluster
	Action  string
	Target  string
}

// ClustersController implements Controller functionality for the Cluster model
type ClustersController struct {
	sync.Mutex
	*render.Renderer

	Events chan ClusterEvent
	Client *client.Client

	clusters []*capi.Cluster
}

// Context is a middleware that parses the Cluster namespace and ID from a request, loads the
// corresponding Cluster model, and makes it available as part of the request's
// context
func (s *ClustersController) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx     context.Context
			cluster *capi.Cluster
			err     error
		)

		namespace := chi.URLParam(r, "namespace")
		id := chi.URLParam(r, "id")

		if namespace != "" && id != "" {
			cluster, err = s.clusterFromIDString(namespace, id)

			if err != nil {
				http.Error(w, fmt.Sprintf("Fatal error: %v", err), http.StatusInternalServerError)

				return
			}

			ctx = context.WithValue(r.Context(), ContextKeyCluster, cluster)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *ClustersController) clusterFromIDString(namespace, id string) (*capi.Cluster, error) {
	log.Println("Getting cluster")

	cluster := &capi.Cluster{}

	if err := s.Client.Kubernetes.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: id}, cluster); err != nil {
		return nil, err
	}

	log.Printf("Found cluster %q", id)

	return cluster, nil
}

// Index shows a list of all Clusters.
func (s *ClustersController) Index(w http.ResponseWriter, r *http.Request) {
	log.Println("Listing clusters")

	content := map[string]interface{}{"Clusters": s.clusters, "Current": "Clusters"}

	s.RenderFull(w, http.StatusOK, "list clusters", content)
}

// Get retrieves a single Cluster by its ID.
func (s *ClustersController) Get(w http.ResponseWriter, r *http.Request) {
	var (
		cluster *capi.Cluster
	)

	switch v := r.Context().Value(ContextKeyCluster).(type) {
	case nil:
	case *capi.Cluster:
		cluster = v
	default:
		log.Println("type unknown")
		return
	}

	content := map[string]interface{}{
		"Cluster": cluster,
		"Current": "Clusters",
	}

	s.RenderFull(w, http.StatusOK, "get cluster", content)
}

// Summary shows a summary of all Clusters.
func (s *ClustersController) Summary(w http.ResponseWriter, r *http.Request) {
	count := 0

	for _, cluster := range s.clusters {
		if cluster.Status.ControlPlaneReady {
			count++
		}
	}

	content := map[string]interface{}{
		"TotalClusters":      len(s.clusters),
		"TotalClustersReady": count,
	}

	s.Renderer.RenderPartial(w, http.StatusOK, "clusters summary", content)
}

func (s *ClustersController) Clusters() []*capi.Cluster {
	s.Lock()
	defer s.Unlock()

	return s.clusters
}

func (s *ClustersController) Watch() error {
	skheme := runtime.NewScheme()
	_ = scheme.AddToScheme(skheme)
	_ = capi.AddToScheme(skheme)

	config, err := s.Client.KubernetesClientConfig.ClientConfig()
	if err != nil {
		return err
	}
	cache, err := cache.New(config, cache.Options{Scheme: skheme})
	if err != nil {
		return err
	}

	informer, err := cache.GetInformer(context.TODO(), &capi.Cluster{})
	if err != nil {
		return err
	}

	informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s.Lock()
			defer s.Unlock()

			cluster := obj.(*capi.Cluster)

			s.clusters = append(s.clusters, cluster)

			log.Printf("Added cluster %q", cluster.Name)

			s.Events <- ClusterEvent{Cluster: cluster, Action: "append", Target: "cluster-list"}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			s.Lock()
			defer s.Unlock()

			cluster := newObj.(*capi.Cluster)

			for i, old := range s.clusters {
				if old.UID == cluster.UID {
					s.clusters[i] = cluster

					log.Printf("Updated cluster %q", cluster.Name)

					s.Events <- ClusterEvent{Cluster: cluster, Action: "replace", Target: fmt.Sprintf("cluster-%s-%s", cluster.Namespace, cluster.Name)}

					break
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			s.Lock()
			defer s.Unlock()

			cluster := obj.(*capi.Cluster)

			for i, old := range s.clusters {
				if old.UID == cluster.UID {
					s.clusters[i] = s.clusters[len(s.clusters)-1]
					s.clusters[len(s.clusters)-1] = nil
					s.clusters = s.clusters[:len(s.clusters)-1]

					log.Printf("Deleted cluster %q", cluster.Name)

					s.Events <- ClusterEvent{Cluster: cluster, Action: "remove", Target: fmt.Sprintf("cluster-%s-%s", cluster.Namespace, cluster.Name)}

					break
				}
			}
		},
	})

	ctx := context.Background()

	go cache.Start(ctx)

	if ok := cache.WaitForCacheSync(ctx); ok {
		log.Println("Cluster cache synced.")
	}

	return nil
}

func (s *ClustersController) Socket(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		log.Println("Streaming cluster updates")

		for {
			var (
				content map[string]interface{}
				err     error
			)

			select {
			case event := <-s.Events:
				content = map[string]interface{}{
					"Cluster": event.Cluster,
					"Action":  event.Action,
					"Target":  event.Target,
				}

				log.Printf("Performing %q action on %q target", event.Action, event.Target)
			}

			msg, err := s.Render("cluster list entry stream", content)
			if err != nil {
				log.Println(err)
			}

			websocket.Message.Send(ws, msg)
		}
	}).ServeHTTP(w, r)
}
