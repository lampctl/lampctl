package hue

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	hue_db "github.com/lampctl/lampctl/hue/db"
)

var errInvalidBridge = errors.New("invalid bridge specified")

func (h *Hue) api_hue_bridges_POST(c *gin.Context) {
	v := &hue_db.Bridge{}
	if err := c.ShouldBindJSON(v); err != nil {
		panic(err)
	}
	b := NewBridge(v)
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
	if err := b.Init(); err != nil {
		panic(err)
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
	if err := h.db.Delete(&hue_db.Bridge{}, c.Param("id")); err != nil {
		panic(err)
	}
	func() {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		if _, ok := h.bridges[c.Param("id")]; !ok {
			panic(errInvalidBridge)
		}
		delete(h.bridges, c.Param("id"))
	}()
	c.JSON(http.StatusOK, gin.H{})
}
