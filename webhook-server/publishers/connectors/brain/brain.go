package brainconnect

import (
	"bytes"
	"context"
	"encoding/json"
	branchtagcreationevt "github-webhook/events/branch_tag_creation"
	pushrepositoryevt "github-webhook/events/push_repository"
	eventspublisher "github-webhook/publishers/events"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PublishBatchToBrainData struct {
	BatchContent eventspublisher.UnpublishedEventSlice `json:"batchContent"`
}

func Publish(s []eventspublisher.UnpublishedEventSlice, p *pgxpool.Pool, ctx context.Context) error {
	var wg sync.WaitGroup
	errChannel := make(chan error, len(s))

	for _, eventSlice := range s {
		wg.Add(1)
		go func (slice eventspublisher.UnpublishedEventSlice)  {
			defer wg.Done()

			byteData, err := json.Marshal(slice)
			if err != nil {
				log.Println(err.Error())
				errChannel <- err
				return
			}
			
			req := map[string]json.RawMessage {
				"batchContent": byteData,
			}

			reqJSON, err := json.Marshal(req)
			if err != nil {
				log.Println(err.Error())
				errChannel <- err
				return
			}

			postURL := os.Getenv("BRAIN_UPLOAD_URL")
			switch slice.GetEventType() {
				case pushrepositoryevt.EVENT_TYPE:
					postURL += "/repository-push/batch"
				case branchtagcreationevt.EVENT_TYPE:
					postURL += "/branch-tag-creation/batch"
			}

			request, err := http.NewRequest("POST", postURL, bytes.NewReader(reqJSON))
			if err != nil {
				log.Println(err.Error())
				errChannel <- err
				return
			}

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer " + os.Getenv("AIROPS_API_KEY"))
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Println(err.Error())
				errChannel <- err
				return
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				log.Println("Batch write to brain FAILED")
				return
			}
			
			slice.MarkEventsAsPublished(p, ctx)
		}(eventSlice)
	}

	go func ()  {
		wg.Wait()
		close(errChannel)
	}()

	for err := range errChannel {
		if err != nil {
			return err
		}
	}

	return nil
}