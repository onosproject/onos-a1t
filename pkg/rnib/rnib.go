// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package rnib

import (
	"context"
	"crypto/md5"

	gogotypes "github.com/gogo/protobuf/types"
	uuid2 "github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/env"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/uri"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	toposdk "github.com/onosproject/onos-ric-sdk-go/pkg/topo"
)

var log = logging.GetLogger("rnib")

type TopoClient interface {
	WatchTopoXapps(ctx context.Context, ch chan topoapi.Event) error
	GetXappAspects(ctx context.Context, nodeID topoapi.ID) (*topoapi.XAppInfo, error)
	UpdateXappAspects(ctx context.Context, xappID topoapi.ID) error
	AddA1TEntity(ctx context.Context, nbPort uint32) error
	AddA1TXappRelation(ctx context.Context, xappID topoapi.ID) error
	GetA1TTopoID() topoapi.ID
	GetXappRelationTopoID(xappID topoapi.ID) topoapi.ID
	GetPolicyTypes(ctx context.Context) (map[topoapi.PolicyTypeID]*topoapi.A1PolicyType, error)
	GetXAppIDsForPolicyTypeID(ctx context.Context, policyTypeID string) ([]string, error)
}

// NewClient creates a new topo SDK client
func NewClient() (TopoClient, error) {
	sdkClient, err := toposdk.NewClient()
	if err != nil {
		return &Client{}, err
	}
	cl := &Client{
		client: sdkClient,
	}
	return cl, nil
}

// Client topo SDK client
type Client struct {
	client toposdk.Client
}

func (c *Client) GetXAppIDsForPolicyTypeID(ctx context.Context, policyTypeID string) ([]string, error) {
	targetXAppIDs := make([]string, 0)
	objects, err := c.client.List(ctx, toposdk.WithListFilters(getXappFilter()))
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		xAppObject := &topoapi.XAppInfo{}
		err = object.GetAspect(xAppObject)
		if err != nil {
			return nil, err
		}
		for _, t := range xAppObject.A1PolicyTypes {
			if string(t.ID) == policyTypeID {
				targetXAppIDs = append(targetXAppIDs, string(object.ID))
				break
			}
		}
	}
	return targetXAppIDs, nil
}

func (c *Client) GetPolicyTypes(ctx context.Context) (map[topoapi.PolicyTypeID]*topoapi.A1PolicyType, error) {
	policies := make(map[topoapi.PolicyTypeID]*topoapi.A1PolicyType)
	objects, err := c.client.List(ctx, toposdk.WithListFilters(getXappFilter()))
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		xappObject := &topoapi.XAppInfo{}
		err = object.GetAspect(xappObject)
		if err != nil {
			return nil, err
		}
		for _, t := range xappObject.A1PolicyTypes {
			policies[t.ID] = t
		}
	}

	return policies, nil
}

func (c *Client) GetXappAspects(ctx context.Context, xappID topoapi.ID) (*topoapi.XAppInfo, error) {
	object, err := c.client.Get(ctx, xappID)
	if err != nil {
		return nil, err
	}
	xAppInfo := &topoapi.XAppInfo{}
	err = object.GetAspect(xAppInfo)
	return xAppInfo, err
}

func getXappFilter() *topoapi.Filters {
	controlRelationFilter := &topoapi.Filters{
		KindFilter: &topoapi.Filter{
			Filter: &topoapi.Filter_Equal_{
				Equal_: &topoapi.EqualFilter{
					Value: topoapi.XAPP,
				},
			},
		},
	}
	return controlRelationFilter
}

// WatchTopoXapps watch xapp node connection changes
func (c *Client) WatchTopoXapps(ctx context.Context, ch chan topoapi.Event) error {
	err := c.client.Watch(ctx, ch, toposdk.WithWatchFilters(getXappFilter()))
	if err != nil {
		return err
	}
	return nil
}

// TODO UpdateXappAspects updates xapp aspects
func (c *Client) UpdateXappAspects(ctx context.Context, xappID topoapi.ID) error {
	object, err := c.client.Get(ctx, xappID)
	if err != nil {
		return err
	}

	if object != nil && object.GetEntity().GetKindID() == topoapi.XAPP {
		xappObject := &topoapi.XAppInfo{}
		err := object.GetAspect(xappObject)
		if err != nil {
			return err
		}

		err = object.SetAspect(xappObject)
		if err != nil {
			return err
		}
		err = c.client.Update(ctx, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) AddA1TEntity(ctx context.Context, nbPort uint32) error {
	object := &topoapi.Object{
		ID:   c.GetA1TTopoID(),
		Type: topoapi.Object_ENTITY,
		Obj: &topoapi.Object_Entity{
			Entity: &topoapi.Entity{
				KindID: topoapi.A1T,
			},
		},
		Aspects: make(map[string]*gogotypes.Any),
		Labels:  map[string]string{},
	}

	interfaces := make([]*topoapi.Interface, 1)
	interfaces[0] = &topoapi.Interface{
		IP:   env.GetPodIP(),
		Port: nbPort,
		Type: topoapi.Interface_INTERFACE_A1AP,
	}

	aspect := &topoapi.A1TInfo{
		Interfaces: interfaces,
	}

	err := object.SetAspect(aspect)
	if err != nil {
		return err
	}
	err = c.client.Create(ctx, object)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddA1TXappRelation(ctx context.Context, xappID topoapi.ID) error {
	relationID := c.GetXappRelationTopoID(xappID)
	object := &topoapi.Object{
		ID:   relationID,
		Type: topoapi.Object_RELATION,
		Obj: &topoapi.Object_Relation{
			Relation: &topoapi.Relation{
				KindID:      topoapi.CONTROLS,
				SrcEntityID: c.GetA1TTopoID(),
				TgtEntityID: xappID,
			},
		},
	}
	err := c.client.Create(ctx, object)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			log.Warn("Creating xApp %s control relation %s failed: %v", xappID, relationID, err)
		}
		log.Warnf("xApp Control relation %s already exists (xapp: %s, a1t: %s): %v", relationID, xappID, c.GetA1TTopoID(), err)
		return nil
	}
	return nil
}

func (c *Client) GetA1TTopoID() topoapi.ID {
	return topoapi.ID(uri.NewURI(
		uri.WithScheme("a1"),
		uri.WithOpaque(env.GetPodID())).String())
}

func (c *Client) GetXappRelationTopoID(xappID topoapi.ID) topoapi.ID {
	bytes := md5.Sum([]byte(xappID))
	uuid, err := uuid2.FromBytes(bytes[:])
	if err != nil {
		panic(err)
	}
	return topoapi.ID(uri.NewURI(
		uri.WithScheme("uuid"),
		uri.WithOpaque(uuid.String())).String())
}

var _ TopoClient = &Client{}
