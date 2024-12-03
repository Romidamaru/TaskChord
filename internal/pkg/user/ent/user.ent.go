package ent

import "gorm.io/gorm"

// User represents a Discord user in your application.
type User struct {
	gorm.Model
	Nickname string // Discord nickname
}
