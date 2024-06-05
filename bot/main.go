package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/valyala/fasthttp"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("ENV ERROR !!!!")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			youtubeURL := update.Message.Text

			if !isValidYouTubeURL(youtubeURL) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter YouTube link please")
				bot.Send(msg)
				continue
			}

			waitMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Wait please...")
			waitMsg.ReplyToMessageID = update.Message.MessageID
			bot.Send(waitMsg)

			audioBytes, err := getYouTubeAudio(youtubeURL)
			if err != nil {
				log.Printf("Error getting audio: %s", err)
				continue
			}

			audioFile := tgbotapi.FileBytes{
				Name:  "audio.mp3",
				Bytes: audioBytes,
			}
			audioMsg := tgbotapi.NewAudio(update.Message.Chat.ID, audioFile)
			if _, err := bot.Send(audioMsg); err != nil {
				log.Printf("Error sending audio: %s", err)
			}
		}
	}
}

func isValidYouTubeURL(url string) bool {
	re := regexp.MustCompile(`^(https?://)?(www\.)?(youtube|youtu|youtube-nocookie)\.(com|be)/.+`)
	return re.MatchString(url)
}

func getYouTubeAudio(youtubeURL string) ([]byte, error) {
	apiURL := "http://worker:8080/get_audio"

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(apiURL)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")

	reqBody := []byte(`{"youtube_url":"` + youtubeURL + `"}`)
	req.SetBody(reqBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		if resp.StatusCode() == http.StatusBadRequest {
			return nil, fmt.Errorf("Invalid YouTube URL")
		}
		return nil, fmt.Errorf("Error: %d - %s", resp.StatusCode(), resp.Body())
	}

	return resp.Body(), nil
}
