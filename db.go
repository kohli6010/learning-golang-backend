package main

import (
	"log"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"

)

// registerDB ...
func registerDB() {
	dbAlias := os.Getenv("DB_ALIAS")
	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// print all the environment variables
	log.Printf("DB_ALIAS: %s, DB_DRIVER: %s, DB_HOST: %s, DB_PORT: %s, DB_USER: %s, DB_PASSWORD: %s, DB_NAME: %s, DB_SSLMODE: %s",
		dbAlias, dbDriver, dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	dsn := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " port=" + dbPort + " dbname=" + dbName + " sslmode=" + dbSSLMode

	maxRetries := 5
	retryDelay := 2 * time.Second
	timeout := 10 * time.Second

	var err error
	start := time.Now()
	for i := 0; i < maxRetries; i++ {
		err = orm.RegisterDataBase(dbAlias, dbDriver, dsn)
		if err == nil {
			break
		}
		if time.Since(start) > timeout {
			break
		}
		log.Printf("DB connection failed: %v. Retrying (%d/%d)...", err, i+1, maxRetries)
		time.Sleep(retryDelay)
	}
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	// orm.RunSyncdb(dbAlias, true, true)
	log.Printf("Database %s registered successfully with driver %s", dbAlias, dbDriver)
}
