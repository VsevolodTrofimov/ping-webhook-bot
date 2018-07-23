package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/teris-io/shortid"
	"gopkg.in/telegram-bot-api.v4"
	"time"
)

// Maybe I should back these up
var id2intent = map[int64]Intent{}

// Bot is telegram bot that reports all those pings to the user
func Bot(c chan Ping) {
	// init
	fmt.Println("[Bot] Starting")

	token := getConf().Token
	fmt.Println("[Bot] Using token", token)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println("[Bot] ERROR in init:", err)
		panic(err)
	}

	db := DBConnect()
	createProject := makeCreateProject(bot, db)

	u := tgbotapi.NewUpdate(0)            // WTF
	u.Timeout = 60                        // WTF
	updates, err := bot.GetUpdatesChan(u) // WTF

	fmt.Println("[Bot] Running")

	// The bot itself
	for {
		select {
		case ping := <-c:
			kindStr, ok := kind2emoji[ping.Kind]
			if !ok {
				kindStr = "[" + ping.Kind + "]"
			}

			messageText := fmt.Sprintf(
				templatePing,
				kindStr,
				ping.Proj.Name,
				ping.Time.Format(time.UnixDate),
				ping.Val,
			)

			tgMessage := tgbotapi.NewMessage(ping.Proj.User, messageText)
			tgMessage.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(tgMessage)

		case update := <-updates:
			if update.Message == nil {
				continue
			}

			fmt.Println("[update]", update.Message.Chat.UserName, update.Message.Command())

			command := update.Message.Command()
			user := update.Message.Chat.ID

			switch command {
			case "new":
				name := update.Message.CommandArguments()
				createProject(user, name)

			case "start":
				tgMessage := tgbotapi.NewMessage(user, messageWelcomePre)
				bot.Send(tgMessage)

				createProject(user, "Project 1")

				tgMessage = tgbotapi.NewMessage(user, messageWelcomePost)
				bot.Send(tgMessage)

			case "list":
				var projects []Project
				db.Where("user = ?", user).Find(&projects)

				message := printProjects(projects)
				tgMessage := tgbotapi.NewMessage(user, message)
				tgMessage.ParseMode = tgbotapi.ModeMarkdown

				bot.Send(tgMessage)

			case "rename":
				buttons := createProjectButtons(db, update.Message.Chat.ID)

				id2intent[user] = Intent{
					Kind:    IntentRenameChooseProjest,
					Payload: "",
				}

				tgMessage := tgbotapi.NewMessage(user, messageRenameChoose)
				tgMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)
				bot.Send(tgMessage)

			case "delete":
				buttons := createProjectButtons(db, user)
				tgMessage := tgbotapi.NewMessage(user, messageDeleteChoose)
				tgMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)

				id2intent[user] = Intent{
					Kind:    IntentDelete,
					Payload: "",
				}

				bot.Send(tgMessage)

			case "help":
				user := update.Message.Chat.ID
				bot.Send(tgbotapi.NewMessage(user, messageHelp))

			// Derive from intent
			default:
				user := update.Message.Chat.ID
				text := update.Message.Text
				intent := id2intent[user]

				switch intent.Kind {
				case IntentRenameChooseProjest:
					fmt.Println("[Bot] Chosen for rename", text)
					var project Project

					err := db.Where(
						"user = ?", user,
					).First(
						&project, "name = ?", text,
					).RecordNotFound()

					if err {
						fmt.Println("[Bot] ERROR finding project", err, text)
						message := fmt.Sprintf(templateNotFound, text)
						bot.Send(tgbotapi.NewMessage(user, message))
						continue
					}

					id2intent[user] = Intent{
						Kind:    IntentRenameNewName,
						Payload: project.UUID,
					}

					tgMessage := tgbotapi.NewMessage(user, messageRenameNewName)
					tgMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(tgMessage)
					fmt.Println("[Bot] Ready to rename", project)

				case IntentRenameNewName:
					uuid := intent.Payload
					fmt.Println("Applying rename for uuid", uuid)
					var project Project
					db.First(&project, "uuid = ?", uuid)
					oldName := project.Name

					project.Name = text
					db.Save(&project)

					delete(id2intent, user)
					message := fmt.Sprintf(
						templateRenameDone,
						oldName,
						project.Name,
					)
					bot.Send(tgbotapi.NewMessage(user, message))
					fmt.Println("[Bot] Applied rename", uuid)

				case IntentDelete:
					fmt.Println("[Bot] Chosen for deletion", text)
					var project Project

					err := db.Delete(
						&project,
						"user = ? AND name = ?", user, text,
					).RecordNotFound()

					if err {
						fmt.Println("[Bot] ERROR finding project", err, text)
						message := fmt.Sprintf(templateNotFound, text)
						bot.Send(tgbotapi.NewMessage(user, message))
						continue
					}

					delete(id2intent, user)

					message := fmt.Sprintf(templateDeleteDone, project.Name)
					tgMessage := tgbotapi.NewMessage(user, message)
					tgMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(tgMessage)
					fmt.Println("Deleted")

				default:
					fmt.Println("[Bot] No idea what this was about:", text)
					bot.Send(tgbotapi.NewMessage(user, messageNoIntent))
					continue
				}
			}
		}
	}
}

func makeCreateProject(bot *tgbotapi.BotAPI, db *gorm.DB) func(user int64, name string) {
	return func(user int64, name string) {
		uuid, err := shortid.Generate()
		if err != nil {
			panic(err)
		}
		if name == "" {
			name = uuid
		}

		fmt.Println("[Bot] Creating new project", uuid)

		proj := Project{
			User: user,
			Name: name,
			UUID: uuid,
		}
		db.Create(&proj)

		message := fmt.Sprintf(
			templateCreated,
			name,
			uuid,
		)
		tgMessage := tgbotapi.NewMessage(user, message)
		tgMessage.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(tgMessage)

		fmt.Println("[Bot] Created new project", uuid)
	}
}
