// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package cli

import (
	"context"

	a1eistore "github.com/onosproject/onos-a1t/pkg/store/a1ei"
	a1pstore "github.com/onosproject/onos-a1t/pkg/store/a1p"
	substore "github.com/onosproject/onos-a1t/pkg/store/subscription"

	a1tapi "github.com/onosproject/onos-a1t/pkg/southbound/a1t"

	"github.com/onosproject/onos-lib-go/pkg/logging/service"

	"google.golang.org/grpc"
)

// NewService returns a new A1T interface service.
func NewService(subscriptionStore substore.Store, policiesStore a1pstore.Store, eijobsStore a1eistore.Store) service.Service {
	return &Service{
		subscriptionStore: subscriptionStore,
		policiesStore:     policiesStore,
		eijobsStore:       eijobsStore,
	}
}

// Service is a service implementation for administration.
type Service struct {
	service.Service
	subscriptionStore substore.Store
	policiesStore     a1pstore.Store
	eijobsStore       a1eistore.Store
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{
		subscriptionStore: s.subscriptionStore,
		policiesStore:     s.policiesStore,
		eijobsStore:       s.eijobsStore,
	}
	a1tapi.RegisterA1TServer(r, server)
}

// Server implements the A1T gRPC service for administrative facilities.
type Server struct {
	subscriptionStore substore.Store
	policiesStore     a1pstore.Store
	eijobsStore       a1eistore.Store
}

func (s *Server) Get(ctx context.Context, request *a1tapi.GetRequest) (*a1tapi.GetResponse, error) {

	response := &a1tapi.GetResponse{}

	return response, nil
}

func (s *Server) List(ctx context.Context, request *a1tapi.GetRequest) (*a1tapi.ListResponse, error) {

	response := &a1tapi.ListResponse{}

	return response, nil
}

func (s *Server) Watch(request *a1tapi.GetRequest, server a1tapi.A1T_WatchServer) error {

	// response := &a1tapi.GetResponse{}

	return nil
}

func (s *Server) Create(ctx context.Context, request *a1tapi.CreateRequest) (*a1tapi.CreateResponse, error) {

	response := &a1tapi.CreateResponse{}

	return response, nil
}

func (s *Server) Update(ctx context.Context, request *a1tapi.UpdateRequest) (*a1tapi.UpdateResponse, error) {

	response := &a1tapi.UpdateResponse{}

	return response, nil
}

func (s *Server) Delete(ctx context.Context, request *a1tapi.DeleteRequest) (*a1tapi.DeleteResponse, error) {

	response := &a1tapi.DeleteResponse{}

	return response, nil
}
