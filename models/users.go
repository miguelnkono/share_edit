package models

import (
	"net"

	"github.com/google/uuid"
)

type User struct {
	Password string
	UserName string
	Id       uuid.UUID

	Status bool

	address net.Addr
}

func (u *User) SetAddress(addr net.Addr) {
	u.address = addr
}

func (u User) GetAddress() net.Addr { return u.address }
