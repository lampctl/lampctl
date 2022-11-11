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
	Groups []*registry.Group
	Lamps  []*registry.Lamp
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
	p, err := s.registry.GetProvider(c.Param("id"))
	if err != nil {
		panic(err)
	}
	v := []*registry.Change{}
	if err := c.ShouldBindJSON(&v); err != nil {
		panic(err)
	}
	if err := p.Apply(v); err != nil {
		panic(err)
	}
}
