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

func searchRequest(client *search.APIClient, index, query string) (*search.SearchResponse, error) {
	searchParams := search.SearchParams{
		SearchParamsObject: search.
			NewEmptySearchParamsObject().
			SetQuery(query),
	}

	response, err := client.SearchSingleIndex(
		client.
			NewApiSearchSingleIndexRequest(index).
			WithSearchParams(&searchParams),
		// Add a custom HTTP header to this request
		search.WithHeaderParam("extra-header", "greetings"),
		// Add query parameters to this request
		search.WithQueryParam("queryParam", "value"),
	)
	if err != nil {
		return nil, err
	}

	return response, nil
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

	client, err := setupClient(appID, apiKey)
	if err != nil {
		panic(err)
	}

	record := map[string]any{
		"objectID": "object-1",
		"name":     "test record",
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
	fmt.Println("client.Search response:")
	fmt.Println(searchResp.Results)

	results, err := searchRequest(client, indexName, "test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n")
	fmt.Println("client.SearchSingleIndex response:")
	// hits := results.GetHits()
	for _, hit := range results.GetHits() {
		fmt.Println(hit)
	}
}
