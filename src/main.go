package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

const slashCommand = "/donowall"

func slashCommandHandler(client *Client) http.HandlerFunc {
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

			params := &slack.Msg{Text: s.Text}

			code := params.Text
			jsonResponse := client.Get(code)

			// Create a JSON response
			w.Header().Set("Content-Type", "application/json")

			jsonResponseMarshal, err := json.Marshal(jsonResponse)

			if err != nil {
				log.Fatalf("error while JSON marshalling response: %s", err)
			}

			w.Write(jsonResponseMarshal)

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	godotenv.Load(".env")

	client, err := NewClient()

	if err != nil {
		log.Fatal("Failed to initialise client")
	}

	http.HandleFunc("/receive", slashCommandHandler(client))

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
