package gpio

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/db"
	register_db "github.com/lampctl/lampctl/gpio/db"
)

func (g *GPIO) api_registers_POST(c *gin.Context) {
	if err := g.db.Transaction(func(conn *db.Conn) error {
		r := &register_db.Register{}
		if err := c.ShouldBindJSON(r); err != nil {
			return err
		}
		if err := conn.Save(r).Error; err != nil {
			return err
		}
		func() {
			g.mutex.Lock()
			defer g.mutex.Unlock()
			g.registers[r.ID] = NewRegister(r)
		}()
		c.JSON(http.StatusOK, r)
		return nil
	}); err != nil {
		panic(err)
	}
}

func (g *GPIO) api_registers_id_GET(c *gin.Context) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	r, err := g.findRegister(c.Param("id"))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, r)
}

func (g *GPIO) api_registers_id_POST(c *gin.Context) {
	if err := g.db.Transaction(func(conn *db.Conn) error {
		r := &register_db.Register{}
		if err := c.ShouldBindJSON(r); err != nil {
			return err
		}
		if err := conn.Where("id = ?", c.Param("id")).Updates(r).Error; err != nil {
			return err
		}
		if err := func() error {
			g.mutex.Lock()
			defer g.mutex.Unlock()
			o, err := g.findRegister(c.Param("id"))
			if err != nil {
				return err
			}
			delete(g.registers, o.ID)
			g.registers[r.ID] = NewRegister(r)
			return nil
		}(); err != nil {
			return err
		}
		c.JSON(http.StatusOK, r)
		return nil
	}); err != nil {
		panic(err)
	}
}

func (g *GPIO) api_registers_id_DELETE(c *gin.Context) {
	if err := g.db.Transaction(func(conn *db.Conn) error {
		if err := g.db.Delete(&register_db.Register{}, c.Param("id")).Error; err != nil {
			return err
		}
		if err := func() error {
			g.mutex.Lock()
			defer g.mutex.Unlock()
			o, err := g.findRegister(c.Param("id"))
			if err != nil {
				return err
			}
			delete(g.registers, o.ID)
			return nil
		}(); err != nil {
			return err
		}
		c.JSON(http.StatusOK, gin.H{})
		return nil
	}); err != nil {
		panic(err)
	}
}
