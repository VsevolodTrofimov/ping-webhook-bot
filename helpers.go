package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"strings"
)

// Map liek in js
func Map(vs []Project, f func(Project) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func printProject(proj Project) string {
	return fmt.Sprintf("%s %s", proj.UUID, proj.Name)
}

func printProjects(projects []Project) string {
	title := "Id        Name"
	lines := Map(projects, printProject)
	return "```\n" + title + "\n" + strings.Join(lines, "\n") + "```"
}

func createProjectButtons(db *gorm.DB, user int64) [][]tgbotapi.KeyboardButton {
	var projects []Project
	db.Where("user = ?", user).Find(&projects)

	total := len(projects)
	rowsCount := total / projectsInRow
	if total%projectsInRow > 0 {
		rowsCount++
	}
	buttons := make([][]tgbotapi.KeyboardButton, 0, rowsCount)
	buffer := make([]tgbotapi.KeyboardButton, 0, projectsInRow)

	for _, project := range projects {
		buffer = append(buffer, tgbotapi.NewKeyboardButton(project.Name))

		if len(buffer) == cap(buffer) {
			buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(buffer...))
			buffer = buffer[:0]
		}
	}
	if len(buttons) < cap(buttons) {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(buffer...))
	}

	return buttons
}

type config struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
}

var conf config

func getConf() config {
	if conf.Token == "" {
		file, err := ioutil.ReadFile("./config.json")
		if err != nil {
			fmt.Println("[Config] error reading config file:", err)
			panic(err)
		}

		json.Unmarshal(file, &conf)
	}

	return conf
}
