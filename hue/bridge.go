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
)

const appName = "lampctl"

var errInvalidResponse = errors.New("invalid response received")

// Bridge represents a connection to a Hue bridge.
type Bridge struct {
	*hue_db.Bridge
	client *http.Client
	lights map[string]*hueLight
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

func (b *Bridge) setState(light_id string, on bool) error {
	l := &hueLight{
		On: &hueOn{
			On: on,
		},
		Dynamics: &hueDynamics{},
	}
	if on {
		l.Dimming = &hueDimming{
			Brightness: 100,
		}
	}
	_, err := b.doPut(
		fmt.Sprintf("/clip/v2/resource/light/%s", light_id),
		l,
	)
	if err != nil {
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
		lights: make(map[string]*hueLight),
	}
}

// Initialize loads the list of lights in a bridge.
func (b *Bridge) Init() error {
	r, err := b.doGet("/clip/v2/resource/light")
	if err != nil {
		return err
	}
	lights := []*hueLight{}
	if err := json.Unmarshal(r.Data, &lights); err != nil {
		return err
	}
	for _, l := range lights {
		b.lights[l.ID] = l
	}
	return nil
}
