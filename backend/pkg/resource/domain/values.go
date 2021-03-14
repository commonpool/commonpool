package domain

import (
	"encoding/json"
	"math"
)

type ValueThreshold struct {
	Description string `json:"description"`
}

type ValueDimension struct {
	Name         string           `json:"name"`
	Summary      string           `json:"summary"`
	Range        ValueRange       `json:"range"`
	DefaultValue Value            `json:"defaultValue"`
	Thresholds   []ValueThreshold `json:"thresholds"`
}

type AverageDimensionValue struct {
	DimensionValue
	EvaluationCount int `json:"evaluationCount"`
}

type Value float64

func (v *Value) UnmarshalJSON(data []byte) error {
	var res float64
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*v = Value(res)
	return nil
}

func (v Value) Equals(o Value) bool {
	return math.Abs(float64(v)-float64(o)) < 0.000001
}

type ValueRange struct {
	From Value `json:"from"`
	To   Value `json:"to"`
}

func (v ValueRange) Equals(o ValueRange) bool {
	return v.From.Equals(o.From) && v.To.Equals(o.To)
}

type DimensionValue struct {
	DimensionName string     `json:"dimensionName"`
	ValueRange    ValueRange `json:"valueRange" gorm:"embedded"`
}

func (d DimensionValue) Equals(o DimensionValue) bool {
	if d.DimensionName != o.DimensionName {
		return false
	}
	return d.ValueRange.Equals(o.ValueRange)
}

type ValueEstimations []DimensionValue

func (v *ValueEstimations) UnmarshalJSON(data []byte) error {
	var val []DimensionValue
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

func (v ValueEstimations) Equals(o ValueEstimations) bool {

	if len(v) != len(o) {
		return false
	}

	vmap := map[string]DimensionValue{}
	omap := map[string]DimensionValue{}

	for _, value := range v {
		vmap[value.DimensionName] = value
	}
	for _, value := range o {
		omap[value.DimensionName] = value
	}

	for vkey, vvalue := range vmap {
		ovalue, ok := omap[vkey]
		if !ok {
			return false
		}
		if !ovalue.Equals(vvalue) {
			return false
		}
	}

	return true

}

var SentimentalValue = ValueDimension{
	Name:         "sentimental",
	Summary:      "Sentimental Value",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "I hate this with a passion",
		}, {
			Description: "My heart couldn't care less",
		}, {
			Description: "My heart bleeds for this",
		},
	},
}

var UrgencyValue = ValueDimension{
	Name:         "urgency",
	Summary:      "Urgency Value",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "I will die without this",
		}, {
			Description: "I need this today",
		}, {
			Description: "I need this shortly",
		}, {
			Description: "This can wait",
		}, {
			Description: "This is not urgent",
		},
	},
}

var FunctionalValue = ValueDimension{
	Name:         "functional",
	Summary:      "Functional value",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "This is absolutely useless",
		}, {
			Description: "This is not useful",
		}, {
			Description: "This is somewhat useful",
		}, {
			Description: "This is useful",
		}, {
			Description: "This is really useful",
		}, {
			Description: "This infinitely useful",
		},
	},
}

var TimeValue = ValueDimension{
	Name:         "time",
	Summary:      "Time Value",
	DefaultValue: 0.5,
	Range: ValueRange{
		From: 0,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "A second",
		}, {
			Description: "A minute",
		}, {
			Description: "A few minutes",
		}, {
			Description: "Half an hour",
		}, {
			Description: "An hour",
		}, {
			Description: "A few hours",
		}, {
			Description: "Half a day",
		}, {
			Description: "A day",
		}, {
			Description: "A few days",
		}, {
			Description: "A week",
		}, {
			Description: "A few weeks",
		}, {
			Description: "A month",
		}, {
			Description: "A few months",
		}, {
			Description: "A year",
		}, {
			Description: "A few years",
		}, {
			Description: "A decade",
		}, {
			Description: "My entire life",
		},
	},
}

var CulturalValue = ValueDimension{
	Name:         "cultural",
	Summary:      "Cultural Value",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "The world will forever forget this",
		}, {
			Description: "This is a pillar of civilization",
		},
	},
}

var EmotionalValue = ValueDimension{
	Name:         "emotional",
	Summary:      "Emotional Value",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "This makes me feel like shit",
		}, {
			Description: "This makes me feel elated",
		},
	},
}

var MentalValue = ValueDimension{
	Name:         "mental",
	Summary:      "Mental Value",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "This makes people stupid",
		}, {
			Description: "This is edifying",
		},
	},
}

var WorthValue = ValueDimension{
	Name:         "worth",
	Summary:      "Worth",
	DefaultValue: 0,
	Range: ValueRange{
		From: -1,
		To:   1,
	},
	Thresholds: []ValueThreshold{
		{
			Description: "This is worth nothing",
		}, {
			Description: "This is worth something",
		}, {
			Description: "This is worth a fortune",
		},
	},
}

type ValueDimensions []ValueDimension

var AllDimensions = ValueDimensions{
	SentimentalValue,
	TimeValue,
	UrgencyValue,
	FunctionalValue,
	CulturalValue,
	EmotionalValue,
	MentalValue,
	WorthValue,
}
