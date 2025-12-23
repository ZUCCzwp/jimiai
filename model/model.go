package model

import (
	"database/sql/driver"
	"encoding/json"
)

type IntList []int

func (t *IntList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t IntList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type StringList []string

func (t *StringList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), t)
}

func (t StringList) Value() (driver.Value, error) {
	return json.Marshal(t)
}
