// Package routes wires up request paths to their Controllers
package routes

import (
	"github.com/go-chi/chi"

	"github.com/talos-systems/theila/pkg/client"
	"github.com/talos-systems/theila/pkg/controllers"
	"github.com/talos-systems/theila/pkg/render"
)

// Register wires up request paths to controllers for the given router
func Register(r chi.Router, renderer *render.Renderer) {
	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	renderer.Client = c

	dashboardController := &controllers.DashboardController{Renderer: renderer, Client: c}
	serversController := &controllers.ServersController{Renderer: renderer, Client: c, Events: make(chan controllers.ServerEvent, 5)}
	clustersController := &controllers.ClustersController{Renderer: renderer, Client: c, Events: make(chan controllers.ClusterEvent, 5)}

	serversController.Watch()
	clustersController.Watch()

	r.Get("/", dashboardController.Index)

	r.Route("/servers", func(r chi.Router) {
		r.Get("/", serversController.Index)
		r.Get("/summary", serversController.Summary)
		r.Get("/socket", serversController.Socket)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(serversController.Context)

			r.Get("/", serversController.Get)
		})
	})

	r.Route("/clusters", func(r chi.Router) {
		r.Get("/", clustersController.Index)
		r.Get("/summary", clustersController.Summary)
		r.Get("/socket", clustersController.Socket)

		r.Route("/{namespace}/{id}", func(r chi.Router) {
			r.Use(clustersController.Context)

			r.Get("/", clustersController.Get)
		})
	})
}
