package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleNewTODO(db *sql.DB, ctx *gin.Context) {
	type NewTODO struct {
		UserID  string    `json:"user_id"`
		DueDate time.Time `json:"due_date"`
		Info    string    `json:"info"`
	}

	var newTODO NewTODO
	var err error

	if err = ctx.BindJSON(&newTODO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	insert, err := db.Prepare("INSERT INTO todos (TodoID, UserID, CreatedDate, DueDate, Info) VALUES ($1, $2, $3, $4, $5);")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}
	defer insert.Close()

	todoID := uuid.New().String()
	createdDate := time.Now()

	_, err = insert.Exec(todoID, newTODO.UserID, createdDate, newTODO.DueDate, newTODO.Info)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "TODO created"})
}

// func handleFetchTODO(db *sql.DB, ctx *gin.Context) {
// 	var fetchTODO TODO
// 	var err error

// 	if err = ctx.BindJSON(&fetchTODO); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "error", "message": "TODOs fetched", "data": gin.H{}})
// }
