package main

import (
	"context"
	"dna-test/config"
	"dna-test/models"
	s "dna-test/server"
	"dna-test/store"
	"gopkg.in/tylerb/graceful.v1"
	"strconv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configuration := config.NewConfiguration()

	var db models.DbClient
	// store.NewNativeDB() can be used instead
	err := store.NewGormDB(configuration, &db)
	if err != nil {
		panic(err)
	}
	defer db.CloseClient()

	var cache models.CacheClient
	err = store.NewRedis(ctx, configuration, &cache)
	if err != nil {
		panic(err)
	}
	defer cache.CloseClient()

	server := s.NewApiServer(configuration, db, cache)
	server.Init()

	graceful.Run(":"+strconv.Itoa(server.Configuration.Port), 0, server.Framework)
}
