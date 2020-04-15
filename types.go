package main

import "math/big"

type Door struct {
	Id     int64  `json:"id,omitempty"`
	System int64  `json:"system,omitempty"`
	KeyX *big.Int     `json:"keyX,omitempty"`
	KeyY *big.Int     `json:"keyY,omitempty"`
	Challenge []byte  `json:"challenge,omitempty"`
	Response  []byte  `json:"response,omitempty"`
	Name   string `json:"name,omitempty"`
	Method int    `json:"method,omitemtpy"`
}

type Role struct {
	System int64  `json:"system,omitempty"`
	Name   string `json:"name,omitempty"`
}

type User struct {
	System int64  `json:"system,omitempty"`
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Pin    string `json:"pin,omitempty"`
	CardNo string `json:"cardno,omitempty"`
}

type System struct {
	Id int64    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	TotpKey []byte
}

type Admin struct {
	Id int64          `json:"id,omitempty"`
	KeyX *big.Int     `json:"keyX,omitempty"`
	KeyY *big.Int     `json:"keyY,omitempty"`
	Challenge []byte  `json:"challenge,omitempty"`
	Response  []byte  `json:"response,omitempty"`
	Name  string      `json:"name,omitempty"`
	Email string      `json:"email,omitempty"`
	Phone string      `json:"phone,omitempty"`
}

