package hue

import (
	"encoding/json"
)

const (
	hueTypeLight        = "light"
	hueTypeGroupedLight = "grouped_light"
	hueTypeZone         = "zone"
	hueTypeBridgeHome   = "bridge_home"
)

type hueRegisterRequest struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}

type hueRegisterResponse []*struct {
	Error *struct {
		Description string `json:"description"`
	} `json:"error"`
	Success *struct {
		Username string `json:"username"`
	} `json:"success"`
}

type hueResponse struct {
	Errors []interface{}   `json:"errors"`
	Data   json.RawMessage `json:"data"`
}

type hueOwner struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
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

type hueColorXY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type hueColor struct {
	XY *hueColorXY `json:"xy"`
}

type hueDynamics struct {
	Duration int64 `json:"duration"`
}

type hueResource struct {
	ID       string       `json:"id,omitempty"`
	Owner    *hueOwner    `json:"owner,omitempty"`
	Metadata *hueMetadata `json:"metadata,omitempty"`
	On       *hueOn       `json:"on,omitempty"`
	Dimming  *hueDimming  `json:"dimming,omitempty"`
	Color    *hueColor    `json:"color,omitempty"`
	Dynamics *hueDynamics `json:"dynamics,omitempty"`
	Type     string       `json:"type,omitempty"`
}

type hueBridge struct {
	BridgeID string `json:"bridge_id"`
}
