package main

import (
	//"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	// ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
