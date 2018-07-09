package dmap

import (
	"encoding/json"
	"fmt"
	"io"
)

var (
	errorEmptyData       = "empty data"
	errorExpectedKey     = "expected key, got %+v of type %T at path %v"
	errorKeyNotFound     = "key %v not found at path %v"
	errorExpectedIndex   = "expected index, got %+v of type %T at path %v"
	errorIndexOutOfRange = "index %v out of range at path %v"
	errorUnexpectedType  = "data at %v is not a map or slice"
	errorNotMapSI        = "data at %v is not a map[string]interface{}"
	errorNotMapII        = "data at %v is not a map[interface{}]interface{}"
	errorNotSliceI       = "data at %v is not a []interface{}"
)

// DMap stores the data and provides a bunch of methods to access and manipulate it.
type DMap struct {
	data interface{}
}

// Init returns a new dmap with the data passed as argument.
func Init(v interface{}) *DMap {
	return &DMap{
		data: v,
	}
}

// ParseJSONBytes returns a new dmap with the JSON bytes unmarshalled.
func ParseJSONBytes(jsonBytes []byte) (*DMap, error) {
	var v interface{}
	err := json.Unmarshal(jsonBytes, &v)
	if err != nil {
		return nil, err
	}

	return &DMap{data: v}, nil
}

// ParseJSONBuffer retuns a new dmap with the JSON buffer unmarshalled.
func ParseJSONBuffer(jsonBuffer io.Reader) (*DMap, error) {
	var v interface{}
	decoder := json.NewDecoder(jsonBuffer)
	err := decoder.Decode(&v)
	if err != nil {
		return nil, err
	}

	return &DMap{data: v}, nil
}

// Data returns the data stored by the dmap.
func (d *DMap) Data() interface{} {
	return d.data
}

// HasData checks if the dmap has any data.
func (d *DMap) HasData() bool {
	return d.data != nil
}

// Get returns the data at a given path. May return a key missing or index out of range error.
func (d *DMap) Get(path ...interface{}) (*DMap, error) {
	if !d.HasData() && len(path) != 0 {
		return nil, fmt.Errorf(errorEmptyData)
	}

	currentData := d.Data()

	for i, p := range path {
		if data, ok := currentData.(map[string]interface{}); ok {
			key, ok := p.(string)
			if !ok {
				return nil, fmt.Errorf(errorExpectedKey, p, p, path[:i+1])
			}

			v, ok := data[key]
			if !ok {
				return nil, fmt.Errorf(errorKeyNotFound, key, path[:i+1])
			}

			currentData = v

		} else if data, ok := currentData.(map[interface{}]interface{}); ok {
			v, ok := data[p]
			if !ok {
				return nil, fmt.Errorf(errorKeyNotFound, p, path[:i+1])
			}

			currentData = v

		} else if data, ok := currentData.([]interface{}); ok {
			index, ok := p.(int)
			if !ok {
				return nil, fmt.Errorf(errorExpectedIndex, p, p, path[:i+1])
			}

			if index < 0 || index >= len(data) {
				return nil, fmt.Errorf(errorIndexOutOfRange, index, path[:i+1])
			}

			currentData = data[index]

		} else {
			return nil, fmt.Errorf(errorUnexpectedType, path[:i+1])
		}
	}

	return &DMap{currentData}, nil
}

// Exists checks whether there is some data at a given path and returns a boolean value.
func (d *DMap) Exists(path ...interface{}) bool {
	_, err := d.Get(path...)
	return err == nil
}

// GetMapSI returns the data at a given path as map[string]interface{}.
func (d *DMap) GetMapSI(path ...interface{}) (map[string]interface{}, error) {
	data, err := d.Get(path...)
	if err != nil {
		return nil, err
	}

	dataMapSI, ok := data.Data().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(errorNotMapSI, path)
	}

	return dataMapSI, nil
}

// GetMapII returns the data at a given path as map[interface{}]interface{}.
func (d *DMap) GetMapII(path ...interface{}) (map[interface{}]interface{}, error) {
	data, err := d.Get(path...)
	if err != nil {
		return nil, err
	}

	dataMapII, ok := data.Data().(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf(errorNotMapII, path)
	}

	return dataMapII, nil
}

// GetSliceI returns the data at a given path as []interface{}.
func (d *DMap) GetSliceI(path ...interface{}) ([]interface{}, error) {
	data, err := d.Get(path...)
	if err != nil {
		return nil, err
	}

	dataSliceI, ok := data.Data().([]interface{})
	if !ok {
		return nil, fmt.Errorf(errorNotSliceI, path)
	}

	return dataSliceI, nil
}

// SetMapSI sets data to a map[string]interface{} at a given path.
func (d *DMap) SetMapSI(data interface{}, key string, path ...interface{}) error {
	parent, err := d.GetMapSI(path...)
	if err != nil {
		return err
	}

	parent[key] = data

	return nil
}

// SetMapII sets data to a map[interface{}]interface{} at a given path.
func (d *DMap) SetMapII(data interface{}, key interface{}, path ...interface{}) error {
	parent, err := d.GetMapII(path...)
	if err != nil {
		return err
	}

	parent[key] = data

	return nil
}

// SetSliceI sets data to a []interface{} at a given path.
func (d *DMap) SetSliceI(data interface{}, index int, path ...interface{}) error {
	parent, err := d.GetSliceI(path...)
	if err != nil {
		return err
	}

	if index < 0 || index >= len(parent) {
		return fmt.Errorf(errorIndexOutOfRange, index, path)
	}

	parent[index] = data

	return nil
}
