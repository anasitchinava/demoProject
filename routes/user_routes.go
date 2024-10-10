package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v4/pgxpool"
    "net/http"
    "go-crud-app/models"
    "go-crud-app/metrics"
    "fmt"
)

func SetupRouter(db *pgxpool.Pool) *gin.Engine {
    router := gin.Default()

    router.Use(func(c *gin.Context) {
		c.Next()
		status := fmt.Sprint(c.Writer.Status())
		metrics.RequestCounter.WithLabelValues(c.Request.Method, status).Inc()
		metrics.ResponseStatus.WithLabelValues(status).Inc()
	})

    router.POST("/users", func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := models.ValidateUser(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var existingUser models.User
        err := db.QueryRow(c, "SELECT id, name, email FROM users WHERE email = $1", user.Email).Scan(&existingUser.ID, &existingUser.Name, &existingUser.Email)
        if err == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
            return
        } else if err.Error() != "no rows in result set" {
            fmt.Println("err:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Check failed"})
            return
        }

        err = db.QueryRow(c, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
            return
        }

        c.JSON(http.StatusCreated, user)
    })

    router.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        var user models.User
        err := db.QueryRow(c, "SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            if err.Error() == "no rows in result" {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user"})
            }
            return
        }

        c.JSON(http.StatusOK, user)
    })

    router.GET("/users", func(c *gin.Context) {
        rows, err := db.Query(c, "SELECT id, name, email FROM users")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
            return
        }
        defer rows.Close()

        users := []models.User{}
        for rows.Next() {
            var user models.User
            if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not scan user"})
                return
            }
            users = append(users, user)
        }

        c.JSON(http.StatusOK, users)
    })

    router.PUT("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := models.ValidateUser(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        _, err := db.Exec(c, "UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
            return
        }

        c.JSON(http.StatusOK, user)
    })

    router.DELETE("/users/:id", func(c *gin.Context) {
        id := c.Param("id")

        _, err := db.Exec(c, "DELETE FROM users WHERE id = $1", id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
            return
        }

        c.JSON(http.StatusNoContent, nil)
    })

    return router
}