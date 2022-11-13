package hue

import (
	"net/http"

	"github.com/gin-gonic/gin"
	hue_db "github.com/lampctl/lampctl/hue/db"
)

func (h *Hue) api_hue_bridges_POST(c *gin.Context) {
	v := &hue_db.Bridge{}
	if err := c.ShouldBindJSON(v); err != nil {
		panic(err)
	}
	b, err := NewBridge(v)
	if err != nil {
		panic(err)
	}
	if b.Username == "" {
		if err := b.register(); err != nil {
			panic(err)
		}
	}
	if b.ID == "" {
		if err := b.getID(); err != nil {
			panic(err)
		}
	}
	if err := h.db.Save(b).Error; err != nil {
		panic(err)
	}
	func() {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		h.bridges[b.ID] = b
	}()
	c.JSON(http.StatusOK, b.Bridge)
}

func (h *Hue) api_hue_bridges_id_DELETE(c *gin.Context) {
	//...
}
