package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pajlada/gobttv"
	"github.com/slack-go/slack"
)

const slashCommand = "/donowall"

func slashCommandHandler(bttv *gobttv.BTTVAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case slashCommand:
			log.Printf("%s slack command is issued\n", slashCommand)
			// Add code here!
			response := "Hello World!"
			params := &slack.Msg{Text: s.Text}
			text := params.Text
			emotes, err := bttv.GetEmotes()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			emoteJson, err := json.Marshal(emotes)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Create a JSON response
			w.Header().Set("Content-Type", "application/json")
			jsonResponse := make(map[string]string)
			jsonResponse["response_type"] = "ephemeral"
			jsonResponse["text"] = fmt.Sprintf("%s : %s\n%s", response, text, emoteJson)
			jsonResponseMarshal, err := json.Marshal(jsonResponse)
			if err != nil {
				log.Fatalf("error while JSON marshalling response: %s", err)
			}

			w.Write([]byte(jsonResponseMarshal))

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	godotenv.Load(".env")

	bttv := gobttv.New()

	http.HandleFunc("/receive", slashCommandHandler(bttv))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("PORT environment variable not found defaulting to port %s", port)
	}

	log.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
