package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Message _
type Message struct {
	//Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	ChatID string `json:"chat_id"`

	// Text of the message to be sent
	Text string `json:"text"`

	// 	Send Markdown or HTML, if you want Telegram apps to show bold, italic, fixed-width text or inline URLs in your bot's message.
	ParseMode string `json:"parse_mode"`

	// Disables link previews for links in this message
	DisableWebPagePreview bool `json:"disable_web_page_preview"`

	// Sends the message silently. Users will receive a notification with no sound.
	DisableNotification bool `json:"disable_notification"`
}

const parseModeMarkdown = "Markdown"

var messageSkeleton Message

func main() {
	checkForDefinedEnvVars()
	defineDefaultMessage()
	tellJobIsSuccessful()
}

func defineDefaultMessage() {
	messageSkeleton.ChatID = os.Getenv("TELEGRAM_CHAT_ID")
	messageSkeleton.ParseMode = parseModeMarkdown
	messageSkeleton.DisableWebPagePreview = true
	messageSkeleton.DisableNotification = false
}

func sendMessage(text string) error {
	messageSkeleton.Text = text

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(messageSkeleton)

	url := fmt.Sprintf("%s%s/%s", "https://api.telegram.org/bot", os.Getenv("TELEGRAM_BOT_TOKEN"), "sendMessage")

	request, _ := http.NewRequest("POST", url, buf)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

func tellJobIsSuccessful() {
	user := fmt.Sprintf("@%s", os.Getenv("GITLAB_USER_LOGIN"))

	url := fmt.Sprintf("%s/commit/%s", os.Getenv("CI_PROJECT_URL"), os.Getenv("CI_COMMIT_SHA"))
	url = fmt.Sprintf("%s (%s)", url, os.Getenv("CI_COMMIT_SHA"))

	message := fmt.Sprintf("%s just deployed an update on branch %s\n%s", user, os.Getenv("CI_COMMIT_REF_NAME"), url)

	err := sendMessage(message)
	if err != nil {
		log.Fatalln(err)
	}
}

func checkForDefinedEnvVars() {
	for _, v := range []string{"TELEGRAM_BOT_TOKEN", "TELEGRAM_CHAT_ID"} {
		if len(os.Getenv(v)) == 0 {
			fmt.Printf("The required environment variable \"%s\" is not defined.\n", v)
			os.Exit(1)
		}
	}
}
