package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"golang.org/x/sync/errgroup"
)

type GetResourceWithSharings struct {
	getResource         *GetResource
	getResourceSharings *GetResourceSharings
}

func NewGetResourceWithSharings(getResource *GetResource, getResourceSharings *GetResourceSharings) *GetResourceWithSharings {
	return &GetResourceWithSharings{getResource: getResource, getResourceSharings: getResourceSharings}
}

func (q *GetResourceWithSharings) Get(ctx context.Context, resourceKey keys.ResourceKey) (*readmodel.ResourceWithSharingsReadModel, error) {

	g, ctx := errgroup.WithContext(ctx)

	var resource *readmodel.ResourceReadModel
	g.Go(func() error {
		var err error
		resource, err = q.getResource.Get(ctx, resourceKey)
		return err
	})

	var sharings []*readmodel.ResourceSharingReadModel
	g.Go(func() error {
		var err error
		sharings, err = q.getResourceSharings.Get(ctx, resourceKey)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &readmodel.ResourceWithSharingsReadModel{
		ResourceReadModel: *resource,
		Sharings:          sharings,
	}, nil

}
