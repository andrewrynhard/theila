package controllers

import (
	"context"
	"net/http"

	metal "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	capi "sigs.k8s.io/cluster-api/api/v1alpha3"

	"github.com/talos-systems/theila/pkg/client"
	"github.com/talos-systems/theila/pkg/render"
)

// DashboardController implements Controller functionality for the Server model
type DashboardController struct {
	*render.Renderer

	Client *client.Client
}

// Index shows the dashboard.
func (d *DashboardController) Index(w http.ResponseWriter, r *http.Request) {
	servers := &metal.ServerList{}

	if err := d.Client.Kubernetes.List(context.TODO(), servers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	serverClasses := &metal.ServerClassList{}

	if err := d.Client.Kubernetes.List(context.TODO(), serverClasses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	clusters := &capi.ClusterList{}

	if err := d.Client.Kubernetes.List(context.TODO(), clusters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	count := 0

	for _, server := range servers.Items {
		if server.Status.InUse {
			count++
		}
	}

	content := map[string]interface{}{
		"TotalServers":          len(servers.Items),
		"TotalServersAllocated": count,
		"TotalServerClasses":    len(serverClasses.Items),
		"TotalClusters":         len(clusters.Items),
		"Current":               "Dashboard",
	}

	d.RenderFull(w, http.StatusOK, "dashboard", content)
}
