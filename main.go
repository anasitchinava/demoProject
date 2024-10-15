package main

import (
    "go-crud-app/database"
    "go-crud-app/routes"
    "go-crud-app/metrics"
    "log"
)

func main() {
    db := database.InitDB()
    defer db.Close()

    metrics.Init()

    // Set up the router
    router := routes.SetupRouter(db)

    // Apply metrics middleware
    router.Use(metrics.MetricsMiddleware())

    router.GET("/metrics", metrics.Handler())

    // Start the server
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }

    log.Println("Server running on :8080")
}