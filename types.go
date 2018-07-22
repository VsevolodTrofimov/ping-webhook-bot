package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

// Project represents of a user's project
type Project struct {
	gorm.Model
	User int64 `gorm:"index"`
	Name string
	UUID string `gorm:"primary_key"`
}

// Ping is minimal representation of reported event
type Ping struct {
	Val  string
	Proj Project
	Kind string
	Time time.Time
}

// Intent represent a way to treat non-command user ection
type Intent struct {
	Kind    int
	Payload string
}

// IntentRenameChooseProjest means that next message will
// be the name of the project to use wants to rename
const IntentRenameChooseProjest = 1

// IntentRenameNewName means that next message will
// be new name of the project with payload in id
const IntentRenameNewName = 2

// IntentDelete means that next message will
// be the name of the project to delete
const IntentDelete = 3
