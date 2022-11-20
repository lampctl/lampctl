package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/registry"
	"github.com/nathan-osman/go-herald"
)

func (s *Server) api_ws_GET(c *gin.Context) {
	s.herald.AddClient(c.Writer, c.Request, nil)
}

type wsMessage struct {
	ProviderID string             `json:"provider_id"`
	Changes    []*registry.Change `json:"changes"`
}

func (s *Server) messageHandler(m *herald.Message, client *herald.Client) {
	v := &wsMessage{}
	if err := json.Unmarshal(m.Data, v); err != nil {
		s.logger.Error().Msg(err.Error())
		return
	}
	if err := s.Apply(v.ProviderID, v.Changes); err != nil {
		s.logger.Error().Msg(err.Error())
	}
}
