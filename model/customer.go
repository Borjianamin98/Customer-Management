package model

import (
	"encoding/json"
	"time"
)

type Customer struct {
	ID           uint      `gorm:"primary_key" json:"cID"`
	Name         string    `json:"cName"`
	Telephone    uint64    `sql:"type:bigint" json:"cTel"`
	Address      string    `json:"cAddress"`
	RegisterDate time.Time `json:"cRegisterDate"`
}

type Alias Customer
type AuxCustomer struct {
	RegisterDate string `json:"cRegisterDate"`
	*Alias
}

func (c Customer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&AuxCustomer{
		RegisterDate: c.RegisterDate.Format("2006-01-02"),
		Alias:        (*Alias)(&c),
	})
}

func (c *Customer) UnmarshalJSON(data []byte) error {
	aux := AuxCustomer{Alias: (*Alias)(c)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	registerDate, error := time.Parse("2006-01-02", aux.RegisterDate)
	if error != nil {
		return error
	}
	c.RegisterDate = registerDate
	return nil
}
