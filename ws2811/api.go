package ws2811

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type api_ws2811_leds_POST_params struct {
	NumLEDs int `json:"num_leds"`
}

func (w *Ws2811) api_ws2811_leds_POST(c *gin.Context) {
	p := &api_ws2811_leds_POST_params{}
	if err := c.ShouldBindJSON(p); err != nil {
		panic(err)
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if err := w.db.SetIntSetting(KeyNumLeds, p.NumLEDs); err != nil {
		panic(err)
	}
	w.numLEDs = p.NumLEDs
	w.free()
	if err := w.init(); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, p)
}
