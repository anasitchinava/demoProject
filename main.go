package main

import (
    "go-crud-app/database"
    "go-crud-app/routes"
    "go-crud-app/metrics"
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    db := database.InitDB()
    defer db.Close()

    metrics.Init()

	// Set up the router
    router := routes.SetupRouter(db)

    router.GET("/metrics", gin.WrapH(metrics.Handler()))

    log.Fatal(router.Run("localhost:8080"))
}