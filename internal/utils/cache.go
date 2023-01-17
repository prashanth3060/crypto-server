package utils

import (
	"crypto-server/internal/model"

	"golang.org/x/sync/syncmap"
)

type ICache interface {
	ReadVal(key string) interface{}
	ReadAllKeys() []string
	Write(key string, v model.Currency)
}

type Cache struct {
	m *syncmap.Map
}

func NewCache() ICache {
	m := syncmap.Map{}
	return &Cache{
		m: &m,
	}
}

func (c Cache) ReadVal(key string) interface{} {
	// if no keys return all
	if len(key) == 0 {
		var vals []model.Currency

		c.m.Range(func(key, value interface{}) bool {
			// cast value to correct format
			val, ok := value.(model.Currency)
			if !ok {
				// this will break iteration
				return false
			}
			// do something with key/value
			vals = append(vals, val)

			// this will continue iterating
			return true
		})

		return vals
	}

	v, ok := c.m.Load(key)
	if !ok {
		return nil
	}
	return v
}

func (c Cache) ReadAllKeys() []string {
	// if no keys return all

	var vals []string

	c.m.Range(func(key, value interface{}) bool {
		// cast value to correct format
		val, ok := key.(string)
		if !ok {
			// this will break iteration
			return false
		}
		// do something with key/value
		vals = append(vals, val)

		// this will continue iterating
		return true
	})

	return vals

}

func (c Cache) Write(k string, v model.Currency) {
	c.m.Store(k, v)
}
