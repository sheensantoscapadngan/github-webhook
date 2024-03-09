package publisher

import (
	"github-webhook/app"
	branchtagcreationevt "github-webhook/events/branch_tag_creation"
	airopsconnect "github-webhook/publishers/connectors/airops"
	eventspublisher "github-webhook/publishers/events"
	"log"
	"net/http"
)

func HandlePublishEvents(a *app.App, w http.ResponseWriter, r *http.Request) {
	var forPublishSlices []eventspublisher.UnpublishedEventSlice

	// use goroutine
	unpublishedBranchTagCreationSlice, err := branchtagcreationevt.GetUnpublishedBranchTagCreation(
		a.Pool,
		r.Context(),
	)

	if err != nil {
		log.Fatal(err.Error())
	}
	if len(unpublishedBranchTagCreationSlice) > 0 {
		forPublishSlices = append(forPublishSlices, unpublishedBranchTagCreationSlice)
	}
	if len(forPublishSlices) > 0 {
		if err := airopsconnect.Publish(forPublishSlices, a.Pool, r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
