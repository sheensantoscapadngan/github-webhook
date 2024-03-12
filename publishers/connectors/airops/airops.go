package airopsconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	eventspublisher "github-webhook/publishers/events"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/nu7hatch/gouuid"
)

type PublishToMemoryData struct {
	Text string `json:"text"`
	Name string `json:"name"`
}

func Publish(s []eventspublisher.UnpublishedEventSlice, p *pgxpool.Pool, ctx context.Context) error {
	var collatedString string
	for _, eventSlice := range s {
		collatedString += eventSlice.ParseString() + "\n"
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("AIROPS COLLATED STRING", collatedString)
	requestData := PublishToMemoryData{
		Text: collatedString,
		Name: "Github events: " + uuid.String(),
	}
	
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Publishing collation uuid:", uuid.String())

	// CALL AIROPS API
	request, err := http.NewRequest("POST", os.Getenv("AIROPS_MEMORY_UPLOAD_URL"), bytes.NewReader(jsonData))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer " + os.Getenv("AIROPS_API_KEY"))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("write to memorystore failed")
	}

	var wg sync.WaitGroup
	// MARK EVENTS AS PUBLISHED
	for _, eventSlice := range s {
		wg.Add(1)
		go func(eventSlice eventspublisher.UnpublishedEventSlice) {
			defer wg.Done()
			eventSlice.MarkEventsAsPublished(p, ctx)
		}(eventSlice)
	}

	wg.Wait()

	return nil;
}