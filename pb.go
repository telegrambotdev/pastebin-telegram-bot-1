package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Pastebin uploader
func PasteBinBot() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("SECRET_TOKEN_PB_BOT"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// Get pastebin private token key
	pasteBinToken := os.Getenv("PASTEBIN_TOKEN")

	if pasteBinToken == "" {
		log.Fatal("PASTEBIN_TOKEN env is missing")
		return
	}

	// Options
	//private := 1 // 1 = unlisted, 0 = public, 2 = private
	pasteBinPostUrl := "https://pastebin.com/api/api_post.php"
	pasteExpiration := "1D"

	// The bot handlings
	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "Pastebin bot")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if m.Text != "" {
			formData := url.Values{
				"api_dev_key":    {pasteBinToken},
				"api_option":     {"paste"},
				"api_paste_code": {m.Text},
				//"api_paste_private": {private},
				"api_paste_expire_date": {pasteExpiration},
			}
			resp, err := http.PostForm(pasteBinPostUrl, formData)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			log.Print(string(body))
			if resp.StatusCode == 200 {
				b.Send(m.Sender, string(body))
			}

		}
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
	})

	b.Start()
}

func main() {
	go PasteBinBot()
	select {}
}
