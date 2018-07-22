package main

var kind2emoji = map[string]string{
	"err":  "⁉",
	"warn": "⚠️",
	"done": "✅",
	"info": "ℹ️",
}

const projectsInRow = 3

// Pre and post stand for 1s project creation here
const messageWelcomePre = "" +
	"Hi, I'm a bot that recives your messages via simple webhook" +
	"and sends them here."

//
const messageWelcomePost = "" +
	"You can rename this project later by using /rename and choosing it"

// ping.Kind, proj.Name, ping.Time (Unix), ping.Val
const templatePing = "*%s %s*\n%s\n\n%s"

const messageDeleteChoose = "Choose a project to delte"

// Proj.Name
const templateDeleteDone = "Deleted %s"

// proj.Name, proj.UUID
const templateCreated = "" +
	"New project *%s* was created. Try using wh.v-trof.ru/%s?message=test?t=done"

const messageRenameChoose = "Choose a project to rename"
const messageRenameNewName = "Input new name for the project"

// old proj.Name, new proj.Name
const templateRenameDone = "Renamed %s to %s"

// proj.Name
const templateNotFound = "Project called %s not found"

const messageHelp = "" +
	"/new {my-name} -- create new project called my-name\n" +
	"/rename -- rename existing project\n" +
	"/list -- list all your projects \n\n" +
	"Use GET wh.v-trof.ru/{id}?m=message&t=type\n" +
	"to recive a message from your project\n" +
	"You can see your project ids using /list\n" +
	"Known types are done, info, warn, err"

const messageNoIntent = "" +
	"Sorry, I have no idea what this was about." +
	"Use some command to get me back on the track"
