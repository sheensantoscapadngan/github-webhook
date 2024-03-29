package publisher

import (
	"github-webhook/app"
	branchtagcreationevt "github-webhook/events/branch_tag_creation"
	pushrepositoryevt "github-webhook/events/push_repository"
	airopsconnect "github-webhook/publishers/connectors/airops"
	brainconnect "github-webhook/publishers/connectors/brain"
	eventspublisher "github-webhook/publishers/events"
	"log"
	"net/http"
	"os"
	"sync"
)

func HandlePublishEvents(a *app.App, w http.ResponseWriter, r *http.Request) {
	var forPublishSliceChannel = make(chan eventspublisher.UnpublishedEventSlice)
	var forPublishSlices []eventspublisher.UnpublishedEventSlice
	var wg sync.WaitGroup

	wg.Add(2)
	go func ()  {
		defer wg.Done()
		if unpublishedBranchTagCreationSlice, err := branchtagcreationevt.GetUnpublishedBranchTagCreation(
			a.Pool,
			r.Context(),
		); err != nil {
			log.Println(err.Error())
		} else if len(unpublishedBranchTagCreationSlice) > 0 {
			forPublishSliceChannel <- unpublishedBranchTagCreationSlice
		}
	}()
	
	go func ()  {
		defer wg.Done()
		if unpublishedRepositoryPushSlice, err := pushrepositoryevt.GetUnpublishedRepositoryPush(
			a.Pool,
			r.Context(),
		); err != nil {
			log.Println(err.Error())
		} else if len(unpublishedRepositoryPushSlice) > 0 {
			forPublishSliceChannel <- unpublishedRepositoryPushSlice
		}
	}()

	go func() {
		wg.Wait()
		close(forPublishSliceChannel)
	}()

	for forPublishSlice := range forPublishSliceChannel {
		forPublishSlices = append(forPublishSlices, forPublishSlice)
	}

	if len(forPublishSlices) > 0 {
		connector := os.Getenv("CONNECTOR")
		log.Println("Using publishing connector:", connector)

		switch connector {
		case "airops":
			if err := airopsconnect.Publish(forPublishSlices, a.Pool, r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "brain":
			if err := brainconnect.Publish(forPublishSlices, a.Pool, r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "no configured connector", http.StatusInternalServerError)
		}
	}
}
