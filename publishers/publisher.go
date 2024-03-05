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
	forPublishSlices = append(forPublishSlices, unpublishedBranchTagCreationSlice)

	airopsconnect.Publish(forPublishSlices, a.Pool, r.Context())
}