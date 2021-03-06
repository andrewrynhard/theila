// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package kubernetes implements the connector that can pull data from the Kubernetes control plane.
package kubernetes

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gertd/go-pluralize"
	"github.com/talos-systems/talos/pkg/machinery/api/resource"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	toolscache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/talos-systems/theila/api/common"
	"github.com/talos-systems/theila/api/socket/message"
	"github.com/talos-systems/theila/internal/backend/runtime"
)

// Name kubernetes runtime string id.
var Name = common.Source_Kubernetes.String()

// New creates new Runtime.
func New() (*Runtime, error) {
	scheme, err := getScheme()
	if err != nil {
		return nil, err
	}

	return &Runtime{
		configs:   map[string]*rest.Config{},
		clients:   map[string]*client{},
		scheme:    scheme,
		pluralize: pluralize.NewClient(),
	}, nil
}

// Runtime implements runtime.Runtime.
type Runtime struct {
	scheme    *k8sruntime.Scheme
	pluralize *pluralize.Client
	configs   map[string]*rest.Config
	clients   map[string]*client
	configsMu sync.RWMutex
	clientsMu sync.RWMutex
}

// Watch creates new kubernetes watch.
func (r *Runtime) Watch(ctx context.Context, request *message.WatchSpec, events chan runtime.Event) error {
	var (
		cfg         *rest.Config
		err         error
		contextName string
	)

	if request.Context != nil {
		contextName = request.Context.Name
	}

	cfg, err = config.GetConfigWithContext(contextName)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)

	if request.Context != nil && request.Context.Cluster != nil {
		cfg, err = r.getKubeconfig(ctx, request.Context)
		if err != nil {
			cancel()

			return err
		}
	}

	w := &Watch{
		resource: request.Resource,
		selector: request.Selector,
		events:   events,
		config:   cfg,
		ctx:      ctx,
		cancel:   cancel,
	}

	if err := w.run(ctx); err != nil {
		cancel()

		return err
	}

	return nil
}

// Get implements runtime.Runtime.
func (r *Runtime) Get(ctx context.Context, setters ...runtime.QueryOption) (interface{}, error) {
	opts := runtime.NewQueryOptions(setters...)

	client, err := r.getOrCreateClient(ctx, opts)
	if err != nil {
		return nil, err
	}

	object, err := r.createObject(client, opts)
	if err != nil {
		return nil, err
	}

	if opts.Name != "" {
		err = client.Get(ctx, types.NamespacedName{
			Namespace: opts.Namespace,
			Name:      opts.Name,
		}, object)
	} else {
		return nil, fmt.Errorf("both resource type and resource name must be defined for Kubernetes Get")
	}

	if err != nil {
		return nil, err
	}

	return object, nil
}

// List implements runtime.Runtime.
func (r *Runtime) List(ctx context.Context, setters ...runtime.QueryOption) (interface{}, error) {
	opts := runtime.NewQueryOptions(setters...)

	client, err := r.getOrCreateClient(ctx, opts)
	if err != nil {
		return nil, err
	}

	// unsafe guess list type
	parts := strings.Split(opts.Resource, ".")
	if !strings.HasSuffix(strings.ToLower(parts[0]), "list") {
		parts[0] = r.pluralize.Singular(parts[0])
		parts[0] += "list"
		opts.Resource = strings.Join(parts, ".")
	}

	object, err := r.createObject(client, opts)
	if err != nil {
		return nil, err
	}

	options := []runtimeclient.ListOption{}

	if opts.LabelSelector != "" {
		var selector labels.Selector

		selector, err = labels.Parse(opts.LabelSelector)
		if err != nil {
			return nil, err
		}

		options = append(options, runtimeclient.MatchingLabelsSelector{
			Selector: selector,
		})
	}

	err = client.List(ctx, object, options...)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// AddContext implements runtime.Runtime.
func (r *Runtime) AddContext(id string, data []byte) error {
	c, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return err
	}

	r.configsMu.Lock()
	r.configs[id] = c
	r.configsMu.Unlock()

	return nil
}

// GetContext implements runtime.Runtime.
func (r *Runtime) GetContext(ctx context.Context, context *common.Context, cluster *common.Cluster) ([]byte, error) {
	var err error

	opts := []runtime.QueryOption{
		runtime.WithName(
			fmt.Sprintf("%s-kubeconfig", cluster.Name),
		),
		runtime.WithType(v1.Secret{}),
	}

	if context != nil {
		opts = append(opts, runtime.WithContext(context.Name))
	}

	if cluster.Namespace != "" {
		opts = append(opts, runtime.WithNamespace(cluster.Namespace))
	}

	s, err := r.Get(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	secret, ok := s.(*v1.Secret)
	if !ok {
		return nil, fmt.Errorf("failed to convert the response to v1.Secret")
	}

	raw, ok := secret.Data["value"]
	if !ok {
		return nil, fmt.Errorf("expected kubeconfig to be placed under 'value' in secret, but nothing was found")
	}

	return raw, nil
}

func (r *Runtime) getKubeconfig(ctx context.Context, context *common.Context) (*rest.Config, error) {
	if context == nil || context.Cluster == nil {
		return nil, fmt.Errorf("kubeconfig cluster must specified")
	}

	cluster := context.Cluster

	id := cluster.Uid

	res := func() *rest.Config {
		r.configsMu.RLock()
		defer r.configsMu.RUnlock()

		return r.configs[id]
	}()

	if res != nil {
		return res, nil
	}

	raw, err := r.GetContext(ctx, context, cluster)
	if err != nil {
		return nil, err
	}

	if err = r.AddContext(id, raw); err != nil {
		return nil, err
	}

	r.configsMu.RLock()
	defer r.configsMu.RUnlock()

	return r.configs[id], nil
}

func (r *Runtime) getOrCreateClient(ctx context.Context, opts *runtime.QueryOptions) (*client, error) {
	client, err := r.getClient(opts.Context)
	if err != nil {
		// initialize the client if it's not initialized and we have Cluster information
		if opts.Cluster != nil {
			_, err = r.getKubeconfig(ctx, &common.Context{
				Name:    opts.Context,
				Cluster: opts.Cluster,
			})
			if err != nil {
				return nil, err
			}

			client, err = r.getClient(opts.Context)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return client, nil
}

func (r *Runtime) getClient(id string) (*client, error) {
	client := func() *client {
		r.clientsMu.Lock()
		defer r.clientsMu.Unlock()

		return r.clients[id]
	}()

	if client != nil {
		return client, nil
	}

	r.configsMu.RLock()
	defer r.configsMu.RUnlock()

	var err error

	defer func() {
		if err == nil {
			r.clientsMu.Lock()
			r.clients[id] = client
			r.clientsMu.Unlock()
		}
	}()

	// first try cached kubeconfigs for discovered clusters.
	if r.configs[id] != nil {
		client, err = newClient(r.configs[id], runtimeclient.Options{Scheme: r.scheme})
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	// then fall back to the local kubeconfig
	cfg, err := config.GetConfigWithContext(id)
	if err != nil {
		return nil, err
	}

	client, err = newClient(cfg, runtimeclient.Options{Scheme: r.scheme})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *Runtime) createObject(c *client, opts *runtime.QueryOptions) (k8sruntime.Object, error) {
	var object k8sruntime.Object

	switch {
	case opts.Type != nil:
		if gvk, ok := opts.Type.(schema.GroupVersionKind); ok {
			var err error

			object, err = r.scheme.New(gvk)
			if err != nil {
				return nil, err
			}
		} else {
			val := reflect.New(reflect.TypeOf(opts.Type))
			object, ok = val.Interface().(k8sruntime.Object)

			if !ok {
				return nil, fmt.Errorf("defined type is not runtime.Object")
			}
		}
	case opts.Resource != "":
		gvr, err := parseResource(opts.Resource)
		if err != nil {
			return nil, err
		}

		gvk, err := c.kindFor(*gvr)
		if err != nil {
			return nil, err
		}

		object, err = r.scheme.New(gvk)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("failed to determine resource type")
	}

	return object, nil
}

// Watch watches kubernetes resources.
type Watch struct {
	cancel   context.CancelFunc
	ctx      context.Context
	config   *rest.Config
	events   chan runtime.Event
	resource *resource.WatchRequest
	selector string
}

func (w *Watch) run(ctx context.Context) error {
	dc, err := dynamic.NewForConfig(w.config)
	if err != nil {
		return err
	}

	namespace := v1.NamespaceAll
	if w.resource.Namespace != "" {
		namespace = w.resource.Namespace
	}

	var filter dynamicinformer.TweakListOptionsFunc
	if w.selector != "" {
		filter = func(options *metav1.ListOptions) {
			options.LabelSelector = w.selector
		}
	}

	dynamicInformer := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, namespace, filter)

	gvr, err := parseResource(w.resource.Type)
	if err != nil {
		return err
	}

	informer := dynamicInformer.ForResource(*gvr)

	if err := informer.Informer().SetWatchErrorHandler(func(reflector *toolscache.Reflector, e error) {
		w.events <- runtime.Event{
			Kind: message.Kind_EventError,
			Spec: e.Error(),
		}
	}); err != nil {
		return err
	}

	informer.Informer().AddEventHandler(toolscache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			w.events <- runtime.Event{
				Kind: message.Kind_EventItemAdd,
				Spec: obj,
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			w.events <- runtime.Event{
				Kind: message.Kind_EventItemUpdate,
				Spec: &runtime.EventUpdate{
					Old: oldObj,
					New: newObj,
				},
			}
		},
		DeleteFunc: func(obj interface{}) {
			w.events <- runtime.Event{
				Kind: message.Kind_EventItemDelete,
				Spec: obj,
			}
		},
	})

	go func() {
		dynamicInformer.Start(ctx.Done())

		<-ctx.Done()
	}()

	dynamicInformer.WaitForCacheSync(ctx.Done())

	return nil
}

func parseResource(resource string) (*schema.GroupVersionResource, error) {
	var gvr *schema.GroupVersionResource

	parts := strings.Split(resource, ".")

	if len(parts) == 2 {
		gvr = &schema.GroupVersionResource{
			Resource: parts[0],
			Version:  parts[1],
		}
	} else {
		gvr, _ = schema.ParseResourceArg(resource)
	}

	if gvr == nil {
		return nil, fmt.Errorf("couldn't parse resource name")
	}

	return gvr, nil
}
