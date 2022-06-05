package main

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dbName       string
	bindAddr     string
	dbURL        string
	mpDatabaseRS *mongo.Database
	ctx          context.Context
)

func init() {
	flag.StringVar(&dbName, "db-name", "mai", "MongoDB database name")
	flag.StringVar(&bindAddr, "bind-addr", ":8080", "Server IP bind")
	flag.StringVar(&dbURL, "db-url", "mongodb://127.0.0.1:27017", "MongoDB database URL")
}

func main() {

	//Parsing flags
	flag.Parse()
	//gin.SetMode(gin.ReleaseMode)

	//connect to DB, returns client and context
	client, ctx := ConnectToDatabase(dbURL)
	mpDatabaseRS = client.Database(dbName)

	defer client.Disconnect(ctx)

	//Init router
	r := gin.Default()

	//InitRoutes
	InitRouters(r)

	r.Run(bindAddr)
}
