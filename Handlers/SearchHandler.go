package Handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
)

// SearchHandler function will handling request from inputs and return the result as json response.
func SearchHandler(order, unit, limit, latlon string, es *elasticsearch.Client, context *gin.Context) {
	var response map[string]interface{}
	var buf bytes.Buffer

	// Split value from latlon separated by a comma. So we can define latitude and longitude easily.
	splitLatLon := strings.Split(latlon, ",")

	// Convert limit from string to int
	limitInt, _ := strconv.Atoi(limit)

	/*
		Query sort in Elasticsearch.
		It will produce query like this:
		{
			"sort": {
				"_geo_distance": {
					"location": {
						"lat": splitLatLon[0],
						"lon": splitLatLon[1]
					},
					"order": __order__,
					"unit": __unit__
				}
			},
		}
	*/
	sort := map[string]interface{}{
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": splitLatLon[0],
					"lon": splitLatLon[1],
				},
				"order": order,
				"unit":  unit,
			},
		},
	}

	// We encode from map string-interface into json format.
	if err := json.NewEncoder(&buf).Encode(sort); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Process the query
	search, err := es.Search(
		es.Search.WithSize(limitInt),
		es.Search.WithIndex("trial_geo"), // the index you defined in Elasticsearch
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	)
	defer search.Body.Close()

	if err != nil {
		log.Fatalf("Error when get data. %s", err)
	}

	if err := json.NewDecoder(search.Body).Decode(&response); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	headerStatus := http.StatusOK
	// Because there are some responses from elasticsearch that we don't use, we only use the value from `hits` field.
	body := gin.H{
		"result": response["hits"].(map[string]interface{})["hits"],
	}
	context.JSON(
		headerStatus,
		body,
	)
}