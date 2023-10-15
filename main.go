package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	request_id()
}

func googleMeetFunction() {
	client, err := getGoogleClient()
	if err != nil {
		log.Fatalf("Unable to authenticate: %v", err)
	}

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}
	// request_ID := uuid.New().String()

	event := &calendar.Event{
		Summary:     "Meeting with Google Meet Link",
		Description: "This is a Google Meet link",
		Start: &calendar.EventDateTime{
			DateTime: "2023-10-15T09:00:00-07:00", // Replace with your desired start time.
			TimeZone: "America/Los_Angeles",       // Replace with your desired time zone.
		},
		End: &calendar.EventDateTime{
			DateTime: "2023-10-15T10:00:00-07:00", // Replace with your desired end time.
			TimeZone: "America/Los_Angeles",       // Replace with your desired time zone.
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				// Replace with a unique request ID.
			},
		},
	}

	calendarID := "primary" // Use "primary" for the user's primary calendar.
	event, err = srv.Events.Insert(calendarID, event).ConferenceDataVersion(1).Do()
	if err != nil {
		log.Fatalf("Unable to create event: %v", err)
	}

	log.Printf("Event created: %s", event.HtmlLink)
}

func getGoogleClient() (*http.Client, error) {
	// Replace 'credentials.json' with the path to your JSON credentials file.
	credFile := "path/to/your/credentials.json"
	b, err := os.ReadFile(credFile)
	if err != nil {
		return nil, err
	}

	// Use the 'calendar.CalendarReadonlyScope' scope for read-only access.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		return nil, err
	}

	client := getClient(config)
	return client, nil
}
func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser, then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func generateUniqueRequestID() string {
	// Use a timestamp to ensure uniqueness.
	timestamp := time.Now().Unix()

	// Generate a random number to further ensure uniqueness.
	randomNumber := rand.Intn(10000) // Adjust the range as needed.

	// Combine the timestamp and random number to create a unique request ID.
	requestID := fmt.Sprintf("%d-%d", timestamp, randomNumber)

	return requestID
}

func request_id() {
	/* types := []string{"eventHangout", "eventNamedHangout", "hangoutsMeet", "addOn"}
	typeIndex := rand.Intn(len(types))
	uniqueKey := GenerateUniqueKey()

	return calendar.ConferenceSolutionKey{
		Type: types[typeIndex],
		Key:  uniqueKey,
	} */
}

func GenerateUniqueKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength := 10 // Adjust the length as needed.
	rand.Seed(time.Now().UnixNano())

	key := make([]byte, keyLength)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}
	return string(key)
}
