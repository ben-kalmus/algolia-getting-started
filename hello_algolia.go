package main

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

func main() {
	appID := "ALGOLIA_APPLICATION_ID"
	// API key with `addObject` and `search` ACL
	apiKey := "ALGOLIA_API_KEY"
	indexName := "test-index"

	record := map[string]any{
		"objectID": "object-1",
		"name":     "test record",
	}

	client, err := search.NewClient(appID, apiKey)
	if err != nil {
		panic(err)
	}

	// Add record to an index
	saveResp, err := client.SaveObject(
		client.NewApiSaveObjectRequest(indexName, record),
	)
	if err != nil {
		panic(err)
	}

	// Wait until indexing is done
	_, err = client.WaitForTask(indexName, saveResp.TaskID)
	if err != nil {
		panic(err)
	}

	// Search for 'test'
	searchResp, err := client.Search(
		client.NewApiSearchRequest(
			search.NewEmptySearchMethodParams().SetRequests(
				[]search.SearchQuery{
					*search.SearchForHitsAsSearchQuery(
						search.NewEmptySearchForHits().SetIndexName(indexName).SetQuery("time"),
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
