package main

import (
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

func setupClient(appID, apiKey string) (*search.APIClient, error) {
	client, err := search.NewClient(appID, apiKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func saveObject(client *search.APIClient, indexName string, record map[string]any) error {
	// Add record to an index
	saveResp, err := client.SaveObject(
		client.NewApiSaveObjectRequest(indexName, record),
	)
	if err != nil {
		return err
	}

	// Wait until indexing is done
	_, err = client.WaitForTask(indexName, saveResp.TaskID)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	appID, ok := os.LookupEnv("APP_ID")
	if !ok {
		panic("env APP_ID not set. Create and run source .env")
	}
	// API key with `addObject` and `search` ACL
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		panic("env API_KEY not set. Create and run source .env")
	}

	indexName := "algolia-tutorial"

	record := map[string]any{
		"objectID": "object-1",
		"name":     "test record",
	}

	client, err := setupClient(appID, apiKey)
	if err != nil {
		panic(err)
	}

	err = saveObject(client, indexName, record)
	if err != nil {
		panic(err)
	}

	// Search for 'test'
	searchResp, err := client.Search(
		client.NewApiSearchRequest(
			search.NewEmptySearchMethodParams().SetRequests(
				[]search.SearchQuery{
					*search.SearchForHitsAsSearchQuery(
						search.NewEmptySearchForHits().SetIndexName(indexName).SetQuery("test"),
					),
				},
			),
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(searchResp.Results)
}
