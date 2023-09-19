package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

type ServerCache struct {
	Keys          []string
	AnnotatedKeys []string
	datasetMap    map[string]Dataset
	JSON          []byte
}

func (cache *ServerCache) Dataset(key string) (Dataset, bool) {
	ds, ok := cache.datasetMap[key]
	return ds, ok
}

func NewServerCache(datasets []Dataset) *ServerCache {
	datasetMap := make(map[string]Dataset)
	keys := make([]string, len(datasets))
	annotatedKeys := make([]string, len(datasets))

	for i, dataset := range datasets {
		keys[i] = strcase.ToCamel(dataset.Name)
		types := strings.ReplaceAll(dataset.Types, " ", "")
		keywords := make([]string, len(dataset.DatasetKeywords))
		for i, keyword := range dataset.DatasetKeywords {
			keywords[i] = keyword.Keywords.Keyword
		}
		keyworkds_str := strings.Join(keywords, " ")
		annotatedKeys[i] = fmt.Sprintf(
			"%s-%s: %s attrs=%d instances=%d keywords: %s", dataset.TaskName(),
			types, keys[i], dataset.NumAttributes, dataset.NumInstances,
			keyworkds_str,
		)
		datasetMap[keys[i]] = dataset
	}

	json, err := json.Marshal(datasetMap)
	if err != nil {
		panic(err)
	}

	return &ServerCache{
		Keys:          keys,
		AnnotatedKeys: annotatedKeys,
		datasetMap:    datasetMap,
		JSON:          json,
	}
}

func serve(datasets []Dataset) {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	cache := NewServerCache(datasets)

	attachEndpoints(r, cache)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Println(err)
			log.Fatal("Server closed unexpectedly")
		}
	}

	log.Println("Server exiting")
}

func attachEndpoints(r *gin.Engine, cache *ServerCache) {
	r.GET("/datasets", func(c *gin.Context) {
		if c.GetHeader("Accept") == "text/plain" {
			c.String(200, strings.Join(cache.AnnotatedKeys, "\n"))
			return
		}
		c.Data(200, gin.MIMEJSON, cache.JSON)
	})

	r.GET("/datasets/:name", func(c *gin.Context) {
		name := c.Param("name")
		dataset, ok := cache.Dataset(name)
		if !ok {
			c.AbortWithStatus(404)
			return
		}
		if c.GetHeader("Accept") == "text/markdown" {
			str, err := dataset.String("markdown")
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			c.String(200, str)
			return
		}
		if c.GetHeader("Accept") == "text/plain" {
			str, err := dataset.String("julia")
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			c.String(200, str)
			return
		}
		c.JSON(200, dataset)
	})
}
