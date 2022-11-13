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
	ID       string       `json:"id"`
	Metadata *hueMetadata `json:"metadata"`
	On       *hueOn       `json:"on"`
	Dimming  *hueDimming  `json:"dimming"`
	Dynamics *hueDynamics `json:"dynamics"`
}
