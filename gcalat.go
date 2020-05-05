package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var (
	calID = ""
)

type Calendar struct {
	Srv *calendar.Service
	ID  string
}

func NewCalendar(credentialPath string, id string) (*Calendar, error) {
	cli, err := getGoogleClient("credentials.json", calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	srv, err := calendar.New(cli)
	if err != nil {
		return nil, err
	}
	return &Calendar{srv, id}, nil
}

func (c *Calendar) ListEvents() ([]*calendar.Event, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := c.Srv.Events.List(c.ID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	return events.Items, nil
}

func (c *Calendar) UpdateEvent(event *calendar.Event) error {
	fmt.Println(event.Summary)
	in := *event
	_, err := c.Srv.Events.Update(c.ID, in.Id, &in).Do()
	return err
}

func run() (err error) {

	cal, err := NewCalendar("credentials.json", calID)
	if err != nil {
		return
	}

	events, err := cal.ListEvents()
	if err != nil {
		return
	}

	e := events[0]
	fmt.Println(e.Summary)
	e.Summary = "hello"
	cal.UpdateEvent(e)

	return
}
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func getGoogleClient(credentialPath string, scope string) (*http.Client, error) {
	cred, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(cred, scope)
	if err != nil {
		return nil, err
	}
	return conf.Client(oauth2.NoContext), nil
}
