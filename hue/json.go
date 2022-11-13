package hue

import (
	"encoding/json"
)

type hueResponse struct {
	Errors []interface{}   `json:"errors"`
	Data   json.RawMessage `json:"data"`
}

type hueMetadata struct {
	Name string `json:"name"`
}

type hueOn struct {
	On bool `json:"on"`
}

type hueDimming struct {
	Brightness float64 `json:"brightness"`
}

type hueDynamics struct {
	Duration int64 `json:"duration"`
}

type hueLight struct {
	ID       string       `json:"id,omitempty"`
	Metadata *hueMetadata `json:"metadata,omitempty"`
	On       *hueOn       `json:"on,omitempty"`
	Dimming  *hueDimming  `json:"dimming,omitempty"`
	Dynamics *hueDynamics `json:"dynamics,omitempty"`
}
