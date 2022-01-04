// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package subscription

import (
	"context"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	toposdk "github.com/onosproject/onos-ric-sdk-go/pkg/topo"
)

type TopoClient interface {
	WatchTopoXapps(ctx context.Context, ch chan topoapi.Event) error
	GetXappAspects(ctx context.Context, nodeID topoapi.ID) (*topoapi.XAppInfo, error)
	UpdateXappAspects(ctx context.Context, xappID topoapi.ID) error
}

// NewClient creates a new topo SDK client
func NewClient() (Client, error) {
	sdkClient, err := toposdk.NewClient()
	if err != nil {
		return Client{}, err
	}
	cl := Client{
		client: sdkClient,
	}
	return cl, nil
}

// Client topo SDK client
type Client struct {
	client toposdk.Client
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

var _ TopoClient = &Client{}
