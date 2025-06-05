package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var cnfg Configuration

// Configuration modelo que contendra un mapa de string/inteface para leer y almacenar la data del archivo de configuracion tipo json
type Configuration struct {
	data map[string]interface{}
}

// New agregar descripcion
func New(fullpath string) (*Configuration, error) {
	b, err := LoadFile(fullpath)
	if err != nil {
		return nil, err
	}
	err = LoadBytes(b, &cnfg)
	if err != nil {
		return nil, err
	}
	return &cnfg, nil

}

// LoadFile lee un archivo
func LoadFile(fullpath string) ([]byte, error) {
	f, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	return f, nil

}

// LoadBytes ...
func LoadBytes(d []byte, c *Configuration) error {
	err := json.Unmarshal(d, &c.data)
	if err != nil {
		return err
	}
	return nil

}

// Validate ..
func (c *Configuration) Validate(names ...string) error {

	for _, v := range names {
		_, ok := c.data[v]
		if !ok {
			return fmt.Errorf("no existe el campo %s", v)
		}
	}

	return nil

}

// Get devuelve el valor del campo si exist, tipo string
func (c *Configuration) Get(name string) (string, error) {
	v, ok := c.data[name].(string)

	if !ok {
		return "", fmt.Errorf("no existe el campo %s", name)
	}

	return v, nil

}

// GetInt devuelve el valor del campo si exist, tipo int
func (c *Configuration) GetInt(name string) (int, error) {
	v, ok := c.data[name].(float64)

	if !ok {
		return 0, fmt.Errorf("no existe el campo %s", name)
	}

	return int(v), nil
}
