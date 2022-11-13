package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	hue_db "github.com/lampctl/lampctl/hue/db"
)

// Bridge represents a connection to a Hue bridge.
type Bridge struct {
	*hue_db.Bridge
	client *http.Client
	lights map[string]*hueLight
}

func (b *Bridge) doRequest(method, path string, body interface{}) (*hueResponse, error) {
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
	resp, err := b.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	hResp := &hueResponse{}
	if err := json.NewDecoder(resp.Body).Decode(hResp); err != nil {
		return nil, err
	}
	return hResp, nil
}

func (b *Bridge) doGet(method string) (*hueResponse, error) {
	return b.doRequest(http.MethodGet, method, nil)
}

func (b *Bridge) doPut(method string, body interface{}) (*hueResponse, error) {
	return b.doRequest(http.MethodPut, method, body)
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

// NewBridge creates and initializes a new Bridge instance.
func NewBridge(bridge *hue_db.Bridge) (*Bridge, error) {
	b := &Bridge{
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
	r, err := b.doGet("/clip/v2/resource/light")
	if err != nil {
		return nil, err
	}
	lights := []*hueLight{}
	if err := json.Unmarshal(r.Data, &lights); err != nil {
		return nil, err
	}
	for _, l := range lights {
		b.lights[l.ID] = l
	}
	return b, nil
}
