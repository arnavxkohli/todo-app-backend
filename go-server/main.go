package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	db, err := sql.Open("postgres", "postgres://postgres:1234@db:5432/todo-db?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return
	}
	defer db.Close()

	log.Println("Successfully connected to the PostgreSQL database")

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Welcome to this todo app!")
	})

	router.POST("/add-user", func(ctx *gin.Context) {
		handleNewUser(db, ctx)
	})

	router.POST("/add-todo", func(ctx *gin.Context) {
		handleNewTODO(db, ctx)
	})

	router.GET("/get-todos", func(ctx *gin.Context) {
		handleFetchTodos(db, ctx)
	})

	router.DELETE("/delete-todo", func(ctx *gin.Context) {
		handleDeleteTodo(db, ctx)
	})

	// gin.SetMode(gin.ReleaseMode)

	err = router.Run(":8000")
	if err != nil {
		log.Fatal("Server startup failed: ", err)
		return
	}
}
