package client

import (
	"log"

	cabpt "github.com/talos-systems/cluster-api-bootstrap-provider-talos/api/v1alpha3"
	cacpt "github.com/talos-systems/cluster-api-control-plane-provider-talos/api/v1alpha3"
	caps "github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metal "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	KubernetesClientConfig clientcmd.ClientConfig
	Kubernetes             client.Client
}

func NewClient() (c *Client, err error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	configOverrides := &clientcmd.ConfigOverrides{}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	log.Printf("Using %v", clientConfig.ConfigAccess().GetDefaultFilename())

	scheme := runtime.NewScheme()

	if err = clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = capi.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = cacpt.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = cabpt.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = caps.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = metal.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = caps.AddToScheme(scheme); err != nil {
		return nil, err
	}

	if err = metal.AddToScheme(scheme); err != nil {
		return nil, err
	}

	kubernetes, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return &Client{
		KubernetesClientConfig: clientConfig,
		Kubernetes:             kubernetes,
	}, nil
}
