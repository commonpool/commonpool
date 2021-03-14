package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"golang.org/x/sync/errgroup"
)

type GetResourceWithSharingsAndValues struct {
	getResource               *GetResource
	getResourceSharings       *GetResourceSharings
	getEvaluations            *GetResourceEvaluations
	getUserResourceEvaluation *GetUserResourceEvaluation
}

func NewGetResourceWithSharingsAndValues(
	getResource *GetResource,
	getResourceSharings *GetResourceSharings,
	getEvaluations *GetResourceEvaluations,
	getUserResourceEvaluation *GetUserResourceEvaluation,
) *GetResourceWithSharingsAndValues {
	return &GetResourceWithSharingsAndValues{
		getResource:               getResource,
		getResourceSharings:       getResourceSharings,
		getEvaluations:            getEvaluations,
		getUserResourceEvaluation: getUserResourceEvaluation,
	}
}

func (q *GetResourceWithSharingsAndValues) Get(
	ctx context.Context,
	resourceKey keys.ResourceKey,
	userKey keys.UserKey,
) (*readmodel.ResourceWithSharingsAndValuesReadModel, error) {

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

	var evaluations []*readmodel.ResourceEvaluationReadModel
	g.Go(func() error {
		var err error
		evaluations, err = q.getEvaluations.Get(ctx, resourceKey)
		return err
	})

	var userEvaluations []*readmodel.ResourceEvaluationReadModel
	g.Go(func() error {
		var err error
		userEvaluations, err = q.getUserResourceEvaluation.Get(ctx, resourceKey, userKey)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	sumsTo := map[string]domain.Value{}
	sumsFrom := map[string]domain.Value{}
	counts := map[string]int{}
	averagesFrom := map[string]domain.Value{}
	averagesTo := map[string]domain.Value{}
	dimensionNames := map[string]bool{}

	for _, evaluation := range evaluations {
		sumsTo[evaluation.DimensionName] += evaluation.ValueRange.From
		sumsFrom[evaluation.DimensionName] += evaluation.ValueRange.To
		counts[evaluation.DimensionName] += 1
		dimensionNames[evaluation.DimensionName] = true
	}

	for key, sum := range sumsFrom {
		count := counts[key]
		avg := float64(sum) / float64(count)
		averagesFrom[key] = domain.Value(avg)
	}

	for key, sum := range sumsTo {
		count := counts[key]
		avg := float64(sum) / float64(count)
		averagesTo[key] = domain.Value(avg)
	}

	var estimations []*domain.AverageDimensionValue
	for dimensionName, _ := range dimensionNames {
		avgEstimation := &domain.AverageDimensionValue{
			domain.DimensionValue{
				DimensionName: dimensionName,
				ValueRange: domain.ValueRange{
					From: averagesFrom[dimensionName],
					To:   averagesTo[dimensionName],
				},
			},
			counts[dimensionName],
		}
		estimations = append(estimations, avgEstimation)
	}

	var userEstimationList domain.ValueEstimations
	for _, model := range userEvaluations {
		userEstimationList = append(userEstimationList, model.DimensionValue)
	}

	return &readmodel.ResourceWithSharingsAndValuesReadModel{
		ResourceReadModel: *resource,
		Sharings:          sharings,
		AverageValues:     estimations,
		Values:            userEstimationList,
	}, nil

}
