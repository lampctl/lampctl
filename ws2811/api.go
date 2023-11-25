package ws2811

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/db"
	ws2811_db "github.com/lampctl/lampctl/ws2811/db"
)

func (w *Ws2811) api_ws2811_channels_POST(c *gin.Context) {
	if err := w.db.Transaction(func(conn *db.Conn) error {
		ch := &ws2811_db.Channel{}
		if err := c.ShouldBindJSON(ch); err != nil {
			return err
		}
		if err := conn.Save(ch).Error; err != nil {
			return err
		}
		func() {
			w.mutex.RLock()
			defer w.mutex.RUnlock()
			w.channels = append(w.channels, ch)
		}()
		c.JSON(http.StatusOK, ch)
		return nil
	}); err != nil {
		panic(err)
	}
}
