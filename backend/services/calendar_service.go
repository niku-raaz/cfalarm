package services

import (
	"context"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func CreateCalendarEvent(tokenJSON []byte, summary string, startTime, endTime string) (string, error) {
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsJSON(tokenJSON))
	if err != nil {
		return "", err
	}

	event := &calendar.Event{
		Summary: summary,
		Start: &calendar.EventDateTime{
			DateTime: startTime,
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
		},
	}

	cEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		return "", err
	}
	return cEvent.Id, nil
}
