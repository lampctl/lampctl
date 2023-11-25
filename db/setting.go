package db

import (
	"errors"
	"strconv"

	"gorm.io/gorm"
)

// Setting provides a simple key / value store.
type Setting struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value string `gorm:"not null" json:"value"`
}

func (c *Conn) getSetting(key string) (string, error) {
	s := &Setting{}
	if err := c.DB.First(s).Error; err != nil {
		return "", err
	}
	return s.Value, nil
}

// GetStringSetting retrieves a setting's string value by key.
func (c *Conn) GetStringSetting(key, def string) (string, error) {
	v, err := c.getSetting(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return def, nil
		}
		return "", err
	}
	return v, nil
}

// GetIntSetting retrieves a setting's integer value by key.
func (c *Conn) GetIntSetting(key string, def int) (int, error) {
	v, err := c.getSetting(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return def, nil
		}
		return 0, err
	}
	return strconv.Atoi(v)
}

// GetIntSetting retrieves a setting's boolean value by key.
func (c *Conn) GetBoolSetting(key string, def bool) (bool, error) {
	var intDef int
	if def {
		intDef = 1
	}
	v, err := c.GetIntSetting(key, intDef)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return def, nil
		}
		return false, err
	}
	return v != 0, nil
}

// SetStringSetting stores a string value for the specified key.
func (c *Conn) SetStringSetting(key, value string) error {
	return c.DB.Save(&Setting{
		Key:   key,
		Value: value,
	}).Error
}

// SetIntSetting stores an integer value for the specified key.
func (c *Conn) SetIntSetting(key string, value int) error {
	return c.SetStringSetting(key, strconv.Itoa(value))
}

// SetBoolSetting stores a boolean value for the specified key.
func (c *Conn) SetBoolSetting(key string, value bool) error {
	var intVal int
	if value {
		intVal = 1
	}
	return c.SetIntSetting(key, intVal)
}
