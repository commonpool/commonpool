package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceEvaluatedPayload struct {
	EvaluatedBy         keys.UserKey       `json:"evaluatedBy"`
	NewEvaluation       ValueEstimations   `json:"newEvaluation"`
	PreviousEvaluations []ValueEstimations `json:"previousEvaluations"`
	IsNewEvaluation     bool               `json:"isNewEvaluation"`
}

type ResourceEvaluated struct {
	eventsource.EventEnvelope
	ResourceEvaluatedPayload `json:"payload"`
}

func NewResourceEvaluated(
	evaluatedBy keys.UserKey,
	newEvaluation ValueEstimations,
	previousEvaluations []ValueEstimations,
	isNewEvaluation bool) ResourceEvaluated {
	return ResourceEvaluated{
		eventsource.NewEventEnvelope(ResourceEvaluatedEvent, 1),
		ResourceEvaluatedPayload{
			EvaluatedBy:         evaluatedBy,
			NewEvaluation:       newEvaluation,
			PreviousEvaluations: previousEvaluations,
			IsNewEvaluation:     isNewEvaluation,
		},
	}
}
