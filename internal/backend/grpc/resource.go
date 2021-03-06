// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/talos-systems/theila/api/common"
	"github.com/talos-systems/theila/api/rpc"
	"github.com/talos-systems/theila/internal/backend/grpc/router"
	"github.com/talos-systems/theila/internal/backend/runtime"
)

type resourceServer struct {
	rpc.UnimplementedClusterResourceServiceServer
}

func (s *resourceServer) register(server grpc.ServiceRegistrar) {
	rpc.RegisterClusterResourceServiceServer(server, s)
}

func (s *resourceServer) gateway(ctx context.Context, mux *gateway.ServeMux, address string, opts []grpc.DialOption) error {
	return rpc.RegisterClusterResourceServiceHandlerFromEndpoint(ctx, mux, address, opts)
}

// Get returns resource from cluster using Talos or Kubernetes.
func (s *resourceServer) Get(ctx context.Context, in *rpc.GetFromClusterRequest) (*rpc.GetFromClusterResponse, error) {
	r, err := runtime.Get(getSource(ctx).String())
	if err != nil {
		return nil, err
	}

	opts := withContext(router.ExtractContext(ctx))

	opts = append(opts, withResource(in.Resource)...)

	if in.Resource.Id != "" {
		opts = append(opts, runtime.WithName(in.Resource.Id))
	}

	res := &rpc.GetFromClusterResponse{}

	md, _ := metadata.FromIncomingContext(ctx)

	if nodes := md.Get("nodes"); nodes != nil {
		opts = append(opts, runtime.WithNodes(nodes...))
	}

	result, err := r.Get(ctx, opts...)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	res.Body = string(data)

	return res, nil
}

// List returns resources from cluster using Talos or Kubernetes.
func (s *resourceServer) List(ctx context.Context, in *rpc.ListFromClusterRequest) (*rpc.ListFromClusterResponse, error) {
	r, err := runtime.Get(getSource(ctx).String())
	if err != nil {
		return nil, err
	}

	opts := withContext(router.ExtractContext(ctx))

	opts = append(opts, withResource(in.Resource)...)

	for _, s := range in.Selectors {
		opts = append(opts, runtime.WithLabelSelector(s))
	}

	res := &rpc.ListFromClusterResponse{}

	md, _ := metadata.FromIncomingContext(ctx)

	if nodes := md.Get("nodes"); nodes != nil {
		opts = append(opts, runtime.WithNodes(nodes...))
	}

	result, err := r.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	res.Messages = []string{
		string(data),
	}

	return res, nil
}

// GetConfig returns kubeconfig or talos config.
// It's a bit more than just getting a resource that's why it has this custom getter.
func (s *resourceServer) GetConfig(ctx context.Context, cluster *common.Cluster) (*rpc.ConfigResponse, error) {
	r, err := runtime.Get(getSource(ctx).String())
	if err != nil {
		return nil, err
	}

	context := router.ExtractContext(ctx)
	if context == nil {
		return nil, fmt.Errorf("context parameters are required for the config request")
	}

	res, err := r.GetContext(ctx, context, cluster)
	if err != nil {
		return nil, err
	}

	return &rpc.ConfigResponse{
		Data: string(res),
	}, nil
}

func getSource(ctx context.Context) common.Source {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		source := md.Get("source")
		if source != nil {
			if res, ok := common.Source_value[source[0]]; ok {
				return common.Source(res)
			}
		}
	}

	return common.Source_Kubernetes
}

type resource interface {
	GetType() string
	GetNamespace() string
}

func withContext(ctx *common.Context) []runtime.QueryOption {
	opts := []runtime.QueryOption{}
	if ctx == nil {
		return opts
	}

	opts = append(opts, runtime.WithContext(ctx.Name))

	if ctx.Cluster != nil {
		opts = append(opts, runtime.WithCluster(ctx.Cluster))
	}

	return opts
}

func withResource(r resource) []runtime.QueryOption {
	opts := []runtime.QueryOption{}
	if r == nil {
		return opts
	}

	if r.GetNamespace() != "" {
		opts = append(opts, runtime.WithNamespace(r.GetNamespace()))
	}

	if r.GetType() != "" {
		opts = append(opts, runtime.WithResource(r.GetType()))
	}

	return opts
}
