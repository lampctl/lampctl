package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	hue_db "github.com/lampctl/lampctl/hue/db"
	"github.com/lucasb-eyer/go-colorful"
)

const appName = "lampctl"

var (
	errInvalidResource = errors.New("invalid resource specified")
	errInvalidResponse = errors.New("invalid response received")
)

type bridgeResource struct {
	Name     string
	Path     string
	Resource *hueResource
}

// Bridge represents a connection to a Hue bridge.
type Bridge struct {
	*hue_db.Bridge
	client        *http.Client
	resources     map[string]*bridgeResource
	allResourceID string
}

func (b *Bridge) getResource(id string) (*bridgeResource, error) {
	r, ok := b.resources[id]
	if !ok {
		return nil, errInvalidResource
	}
	return r, nil
}

func (b *Bridge) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		v, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(v)
	}
	u := &url.URL{
		Scheme: "https",
		Host:   b.Host,
		Path:   path,
	}
	r, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, err
	}
	r.Header.Add("hue-application-key", b.Username)
	return b.client.Do(r)
}

func (b *Bridge) doRequestAndResponse(method, path string, body interface{}) (*hueResponse, error) {
	r, err := b.doRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	response := &hueResponse{}
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

func (b *Bridge) doGet(method string) (*hueResponse, error) {
	return b.doRequestAndResponse(http.MethodGet, method, nil)
}

func (b *Bridge) doPut(method string, body interface{}) (*hueResponse, error) {
	return b.doRequestAndResponse(http.MethodPut, method, body)
}

func (b *Bridge) register() error {
	r, err := b.doRequest(
		http.MethodPost,
		"/api",
		&hueRegisterRequest{
			DeviceType:        appName,
			GenerateClientKey: true,
		},
	)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	response := hueRegisterResponse{}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return err
	}
	if len(response) < 1 {
		return errInvalidResponse
	}
	if response[0].Error != nil {
		return errors.New(response[0].Error.Description)
	}
	if response[0].Success == nil {
		return errInvalidResponse
	}
	b.Username = response[0].Success.Username
	return nil
}

func (b *Bridge) getID() error {
	r, err := b.doGet("/clip/v2/resource/bridge")
	if err != nil {
		return err
	}
	bridges := []*hueBridge{}
	if err := json.Unmarshal(r.Data, &bridges); err != nil {
		return err
	}
	if len(bridges) < 1 {
		return errInvalidResponse
	}
	b.ID = bridges[0].BridgeID
	return nil
}

func (b *Bridge) setState(
	light_id string,
	on bool,
	brightness float64,
	color string,
	duration int64,
) error {
	r, err := b.getResource(light_id)
	if err != nil {
		return err
	}
	if brightness == 0 {
		brightness = 1.0
	}
	l := &hueResource{
		On: &hueOn{
			On: on,
		},
		Dimming: &hueDimming{
			Brightness: brightness,
		},
		Dynamics: &hueDynamics{
			Duration: duration,
		},
	}
	if on {
		l.Dimming = &hueDimming{
			Brightness: 100,
		}
	}
	if color != "" {
		c, err := colorful.Hex(color)
		if err != nil {
			return err
		}
		x, y, _ := c.Xyy()
		l.Color = &hueColor{
			XY: &hueColorXY{
				X: x,
				Y: y,
			},
		}
	}
	if _, err := b.doPut(r.Path, l); err != nil {
		return err
	}
	l.On.On = on
	return nil
}

// NewBridge creates a new Bridge instance.
func NewBridge(bridge *hue_db.Bridge) *Bridge {
	return &Bridge{
		Bridge: bridge,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{

					// This is not ideal - but the certificate presented by the
					// bridge does not contain a SAN for its IP address and I
					// can't find a way around that
					InsecureSkipVerify: true,
				},
			},
		},
		resources: make(map[string]*bridgeResource),
	}
}

// Init enumerates the contents of the bridge, looking for lights and grouped
// lights, extracting information from what is retrieved.
func (b *Bridge) Init() error {
	r, err := b.doGet("/clip/v2/resource")
	if err != nil {
		return err
	}
	resources := []*hueResource{}
	if err := json.Unmarshal(r.Data, &resources); err != nil {
		return err
	}

	// Due to the way types work in the API, we first need to build a map of
	// resource IDs => resource names
	nameMap := make(map[string]string)
	for _, r := range resources {
		if r.Metadata != nil {
			nameMap[r.ID] = r.Metadata.Name
		}
	}

	for _, r := range resources {
		switch r.Type {

		// For a light, simply add it to the map by its ID
		case hueTypeLight:
			b.resources[r.ID] = &bridgeResource{
				Name:     r.Metadata.Name,
				Path:     fmt.Sprintf("/clip/v2/resource/light/%s", r.ID),
				Resource: r,
			}

		// For grouped lights, do the same, but lookup the name
		case hueTypeGroupedLight:
			name := nameMap[r.Owner.RID]
			if r.Owner.RType == hueTypeBridgeHome {
				name = "All"
				b.allResourceID = r.ID
			}
			b.resources[r.ID] = &bridgeResource{
				Name:     name,
				Path:     fmt.Sprintf("/clip/v2/resource/grouped_light/%s", r.ID),
				Resource: r,
			}
		}
	}
	return nil
}
