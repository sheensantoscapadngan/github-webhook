package event

import (
	"bytes"
	"net/http"
	"os"
)

type GithubEvent interface {
	MarshalEvent()([]byte, error)
}

func WriteToMemory(ge GithubEvent) (error) {
	data, err := ge.MarshalEvent()
	if err != nil {
		return err
	}	

	resp, err := http.Post(os.Getenv("MEMORY_URL"), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}