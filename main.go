package main

import (
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"github.com/yudha-ab/go-elastic-geolocation/Handlers"
)

func main() {
	cfg := &elasticsearch.Config{
		Addresses: []string{
			"http://es_go:9200",
		},
	}

	// Use default settings from gin and ealsticsearch client
	router := gin.Default()
	es, err_c := elasticsearch.NewClient(*cfg)

	if err_c != nil {
		log.Printf("Error: %s", err_c)
	}

	// the URL will be `http(s)://{{host}}/api`
	routeGroup := router.Group("/api")
	{
		routeGroup.GET("/", func(ctx *gin.Context) {
			Handlers.HomeHandler(es, ctx)
		})

		// the endpoint for API search will addressed at `http(s)://{{host}}/api/search`
		routeGroup.GET("/search", func(ctx *gin.Context) {

			/*
				The code below will read URL query as an input to process.
				You can use order, unit, limit, and latlon within the query.
				Example:
				1. http(s)://{{host}}/api/search?latlon=-7.810448,110.4172433
				2. http(s)://{{host}}/api/search?latlon=-7.810448,110.4172433?limit=5
				3. http(s)://{{host}}/api/search?latlon=-7.810448,110.4172433?order=asc (or desc)
				4. http(s)://{{host}}/api/search?latlon=-7.810448,110.4172433?unit=km (or m or miles)
				You can use all query params together, but, the `latlon` is mandatory parameter.
			*/
			order := ctx.DefaultQuery("order", "asc")
			unit := ctx.DefaultQuery("unit", "km")
			limit := ctx.DefaultQuery("limit", "10")
			latLon := ctx.Query("latlon")

			// This conditional statement will filter if the latlon is empty.
			// If you leave it blank, the API will return error message
			if latLon == "" {
				ctx.JSON(
					http.StatusBadRequest, // HTTP bad request (400)
					gin.H{
						"error": "You must specify latlon in query.",
					},
				)
				return
			}

			Handlers.SearchHandler(order, unit, limit, latLon, es, ctx)
		})
	}

	err := router.Run()
	if err != nil {
		log.Printf("Error creating routes: %s", err)
	}
}