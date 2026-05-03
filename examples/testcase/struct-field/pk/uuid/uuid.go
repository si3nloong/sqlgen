package main

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID `sql:"id,pk"`
	Name string
}
