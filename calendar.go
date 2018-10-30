package main

import (
	"cbf-calendar/cbf"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	ASKBlanskoId = "10098"
	CalendarId   = "nlnuj41b0lkvqgs4r2fssb3ajk@group.calendar.google.com"
	GamesUrl     = "http://www.cbf.cz/xml/sched.php?p="
	TableUrl     = "http://www.cbf.cz/xml/table.php?p="
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Missing argument")
	}

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	games, err := cbf.ImportGames(fmt.Sprintf(GamesUrl+"%s", os.Args[1]))
	if err != nil {
		log.Fatalf("Failed to load games from cbf: %v", err)
	}

	table, err := cbf.ImportTable(fmt.Sprintf(TableUrl+"%s", os.Args[1]))
	if err != nil {
		log.Fatalf("Failed to load table from cbf: %v", err)
	}

	for i := 0; i < len(games.Games); i++ {
		if games.Games[i].Teams[0].Id == ASKBlanskoId || games.Games[i].Teams[1].Id == ASKBlanskoId {

			summary := "(" + cbf.FindPositionInTable(games.Games[i].Teams[0].Id, table) + ")" + games.Games[i].Teams[0].Name + " - (" + cbf.FindPositionInTable(games.Games[i].Teams[1].Id, table) + ")" + games.Games[i].Teams[1].Name
			event := &calendar.Event{
				Summary:  summary,
				Id:       games.Games[i].Id,
				Location: games.Games[i].Place,
				Start: &calendar.EventDateTime{
					Date:     games.Games[i].Date,
					TimeZone: "Europe/Prague",
				},
				End: &calendar.EventDateTime{
					Date:     games.Games[i].Date,
					TimeZone: "Europe/Prague",
				},
			}

			if games.Games[i].Time != "00:00:00" {
				event.Start = &calendar.EventDateTime{
					DateTime: games.Games[i].Date + "T" + games.Games[i].Time,
					TimeZone: "Europe/Prague",
				}
				event.End = &calendar.EventDateTime{
					DateTime: games.Games[i].Date + "T" + games.Games[i].Time,
					TimeZone: "Europe/Prague",
				}
			}

			if games.Games[i].Result.Score.A != "0" || games.Games[i].Result.Score.B != "0" {
				event.Summary = event.Summary + " [" + games.Games[i].Result.Score.A + ":" + games.Games[i].Result.Score.B + "]"
			}

			exist, _ := srv.Events.Get(CalendarId, games.Games[i].Id).Do()
			if exist == nil {
				event, err = srv.Events.Insert(CalendarId, event).Do()

				if err != nil {
					log.Fatalf("Unable to create event. %v\n", err)
				}
				fmt.Printf("Event created: %s\n", event.HtmlLink)
			} else {
				event, err = srv.Events.Update(CalendarId, games.Games[i].Id, event).Do()

				if err != nil {
					log.Fatalf("Unable to update event. %v\n", err)
				}
				fmt.Printf("Event updated: %s\n", event.HtmlLink)
			}
		}
	}
}
