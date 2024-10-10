package main

import (
    "go-crud-app/database"
    "go-crud-app/routes"
    "log"
)

func main() {
    db := database.InitDB()
    defer db.Close()

    // Set up the router
    router := routes.SetupRouter(db)

    log.Fatal(router.Run(":8080"))
}
