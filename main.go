package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	since := time.Time{}

	for {
		notifs, err := GitHubNotifications(since)
		if err != nil {
			log.Print(err)
			continue
		}
		since = time.Now()

		for _, notif := range notifs {
			err = sendMessage(&Message{
				Text: notif.Subject.GetTitle(),
				Blocks: []*Section{
					section(markdown(fmt.Sprintf(
						"*%s*\n_%s_\n<%s|View notifications>",
						notif.Subject.GetTitle(),
						notif.Repository.GetFullName(),
						"https://github.com/notifications",
					))),
				},
			})
			if err != nil {
				log.Print(err)
				continue
			}
		}
		time.Sleep(time.Minute)
	}
}

type Message struct {
	Text   string     `json:"text"`
	Blocks []*Section `json:"blocks"`
}

type Section struct {
	Type string    `json:"type"`
	Text *Markdown `json:"text"`
}
type Markdown struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func sendMessage(message *Message) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = http.Post(os.Getenv("SLACK_HOOK_URL"), "application/json", bytes.NewBuffer(b))
	return err
}

func section(body *Markdown) *Section {
	return &Section{
		Type: "section",
		Text: body,
	}
}

func markdown(md string) *Markdown {
	return &Markdown{
		Type: "mrkdwn",
		Text: md,
	}
}
