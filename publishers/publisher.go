package publisher

import (
	"github-webhook/app"
	airopsconnect "github-webhook/publishers/connectors/airops"
	eventspublisher "github-webhook/publishers/events"
	branchpublisher "github-webhook/publishers/events/branch"
	"log"
	"net/http"
)

func HandlePublishEvents(a *app.App, w http.ResponseWriter, r *http.Request) {
	var forPublishSlices []eventspublisher.UnpublishedEventSlice

	// use goroutine
	unpublishedBranchTagCreationSlice, err := branchpublisher.GetUnpublishedBranchTagCreation(
		a,
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
