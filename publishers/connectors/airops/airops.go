package airopsconnect

import (
	"bytes"
	"context"
	"encoding/json"
	eventspublisher "github-webhook/publishers/events"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PublishToMemoryData struct {
	Text string `json:"text"`
}

func Publish(s []eventspublisher.UnpublishedEventSlice, p *pgxpool.Pool, ctx context.Context) error {
	var collatedString string
	for _, eventSlice := range s {
		collatedString += eventSlice.ParseString() + "\n"
	}

	requestData := PublishToMemoryData{
		Text: collatedString,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// CALL AIROPS API
	request, err := http.NewRequest("POST", os.Getenv("AIROPS_MEMORY_UPLOAD_URL"), bytes.NewReader(jsonData))
	if err != nil {
		log.Println(err.Error())
		return err;
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer " + os.Getenv("AIROPS_API_KEY"))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer response.Body.Close()

	var wg sync.WaitGroup
	// MARK EVENTS AS PUBLISHED
	for _, eventSlice := range s {
		eventSliceCopy := eventSlice
		wg.Add(1)
		go func() {
			defer wg.Done()
			eventSliceCopy.MarkEventsAsPublished(p, ctx)
		}()
	}

	wg.Wait()

	return nil;
}