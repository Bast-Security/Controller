package main

import (
	"math/big"
	"time"
)

type Door struct {
	Id        int64    `json:"id,omitempty"`
	System    int64    `json:"system,omitempty"`
	KeyX      *big.Int `json:"keyX,omitempty"`
	KeyY      *big.Int `json:"keyY,omitempty"`
	Challenge []byte   `json:"challenge,omitempty"`
	Response  []byte   `json:"response,omitempty"`
	Name      string   `json:"name,omitempty"`
	Method    int      `json:"method,omitempty"`
}

type Totp struct {
	Expires int64  `json:"expires"`
	Code    string `json:"code"`
}

type Transaction struct {
	DoorName string `json:"door,omitempty"`
	Door int64
	Time time.Time	`json:"time"`
	Pin string      `json:"pin"`
	Card string     `json:"card"`
}

type Role struct {
	Id          int64        `json:"id"`
	System      int64        `json:"system"`
	Name        string       `json:"name"`
	Doors       []Door       `json:"doors"`
}

type Permission struct {
	System int64 `json:"system,omitempty"`
	Door   int64 `json:"door,omitempty"`
	Role   int64 `json:"role,omitempty"`
}

type User struct {
	System int64  `json:"system,omitempty"`
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Pin    string `json:"pin,omitempty"`
	CardNo string `json:"cardno,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Roles  []Role `json:"roles,omitempty"`
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

