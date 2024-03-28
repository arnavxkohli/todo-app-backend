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

func handleFetchTodos(db *sql.DB, ctx *gin.Context) {
	var err error
	userID := ctx.Query("uid")

	todos, err := db.Prepare("SELECT TodoID, CreatedDate, DueDate, Info FROM todos WHERE UserID = $1;")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}
	defer todos.Close() // Close the prepared statement after it's used

	rows, err := todos.Query(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}
	defer rows.Close() // Close the rows after they are scanned

	var fetchedTodos []map[string]interface{}
	for rows.Next() {
		var todoID string
		var createdDate time.Time
		var dueDate sql.NullTime
		var info string

		err = rows.Scan(&todoID, &createdDate, &dueDate, &info)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
			return
		}

		todo := map[string]interface{}{
			"TodoID":      todoID,
			"CreatedDate": createdDate,
			"DueDate":     dueDate,
			"Info":        info,
		}

		fetchedTodos = append(fetchedTodos, todo)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": fetchedTodos})
}

func handleDeleteTodo(db *sql.DB, ctx *gin.Context) {
	var err error
	todoID := ctx.Query("tid")

	del, err := db.Prepare("DELETE FROM todos WHERE TodoID = $1")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}
	defer del.Close()

	_, err = del.Exec(todoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": todoID + " Deleted"})
}
