package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/registry"
)

type providerJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) api_providers_GET(c *gin.Context) {
	response := []*providerJSON{}
	for _, p := range s.registry.Providers() {
		response = append(response, &providerJSON{
			ID:   p.ID(),
			Name: p.Name(),
		})
	}
	c.JSON(http.StatusOK, response)
}

type providerMetaJSON struct {
	Groups []*registry.Group `json:"groups"`
	Lamps  []*registry.Lamp  `json:"lamps"`
}

func (s *Server) api_providers_id_GET(c *gin.Context) {
	p, err := s.registry.GetProvider(c.Param("id"))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, &providerMetaJSON{
		Groups: p.Groups(),
		Lamps:  p.Lamps(),
	})
}

func (s *Server) api_providers_id_apply_POST(c *gin.Context) {
	v := []*registry.Change{}
	if err := c.ShouldBindJSON(&v); err != nil {
		panic(err)
	}
	if err := s.Apply(c.Param("id"), v); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) api_providers_id_apply_all_POST(c *gin.Context) {
	v := &registry.Change{}
	if err := c.ShouldBindJSON(&v); err != nil {
		panic(err)
	}
	p, err := s.registry.GetProvider(c.Param("id"))
	if err != nil {
		panic(err)
	}
	if err := p.ApplyToAll(v); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) api_sequencer_GET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

type sequencerLoadJSON struct {
	MidiFilename    string `json:"midi_filename"`
	MappingFilename string `json:"mapping_filename"`
}

func (s *Server) api_sequencer_load_POST(c *gin.Context) {
	v := &sequencerLoadJSON{}
	if err := c.ShouldBindJSON(v); err != nil {
		panic(err)
	}
	if err := s.sequencer.Load(v.MidiFilename, v.MappingFilename); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) api_sequencer_play_POST(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (s *Server) api_sequencer_stop_POST(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
