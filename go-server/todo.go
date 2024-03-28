package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PriorityLevel string

const (
	Low    PriorityLevel = "low"
	Medium PriorityLevel = "medium"
	High   PriorityLevel = "high"
)

func handleNewTODO(db *sql.DB, ctx *gin.Context) {
	type NewTODO struct {
		UserID   string        `json:"user_id"`
		DueDate  time.Time     `json:"due_date"`
		Info     string        `json:"info"`
		Priority PriorityLevel `json:"priority"`
	}

	var newTODO NewTODO
	var err error

	if err = ctx.BindJSON(&newTODO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	insert, err := db.Prepare("INSERT INTO todos (TodoID, UserID, CreatedDate, DueDate, Info, Priority) VALUES ($1, $2, $3, $4, $5, $6);")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error preparing statement"})
		return
	}
	defer insert.Close()

	todoID := uuid.New().String()
	createdDate := time.Now()

	_, err = insert.Exec(todoID, newTODO.UserID, createdDate, newTODO.DueDate, newTODO.Info, newTODO.Priority)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error executing statement"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "TODO created"})
}

func handleFetchTodos(db *sql.DB, ctx *gin.Context) {
	var err error
	userID := ctx.Query("uid")

	todos, err := db.Prepare("SELECT TodoID, CreatedDate, DueDate, Info, Priority FROM todos WHERE UserID = $1;")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error preparing statement"})
		return
	}
	defer todos.Close() // Close the prepared statement after it's used

	rows, err := todos.Query(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error executing statement"})
		return
	}
	defer rows.Close() // Close the rows after they are scanned

	var fetchedTodos []map[string]interface{}
	for rows.Next() {
		var todoID string
		var createdDate time.Time
		var dueDate sql.NullTime
		var info string
		var priority PriorityLevel

		err = rows.Scan(&todoID, &createdDate, &dueDate, &info, &priority)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error formatting received data"})
			return
		}

		todo := map[string]interface{}{
			"TodoID":      todoID,
			"CreatedDate": createdDate,
			"DueDate":     dueDate,
			"Info":        info,
			"Priority":    priority,
		}

		fetchedTodos = append(fetchedTodos, todo)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": fetchedTodos})
}

func handleUpdateTodo(db *sql.DB, ctx *gin.Context) {
	type UpdateTODO struct {
		TodoID   string         `json:"todo_id"`
		DueDate  *time.Time     `json:"due_date"`       // Pointer as nullable value
		Info     *string        `json:"info"`           // Pointer as nullable value
		Priority *PriorityLevel `json:"priority_level"` // Pointer as nullable value
	}

	var updateTODO UpdateTODO
	var err error

	err = ctx.BindJSON(&updateTODO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if updateTODO.TodoID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "No todo_id provided"})
		return
	}

	if updateTODO.Priority == nil && updateTODO.DueDate == nil && updateTODO.Info == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "No fields provided for update"})
		return
	}

	updateQuery := "UPDATE todos SET"
	var args []interface{}
	argcount := 1 // Using ? as placeholders does not work, use $ and increment the argument count

	if updateTODO.DueDate != nil {
		updateQuery += " DueDate = $" + strconv.Itoa(argcount) + ","
		args = append(args, *updateTODO.DueDate)
		argcount++
	}

	if updateTODO.Info != nil {
		updateQuery += " Info = $" + strconv.Itoa(argcount) + ","
		args = append(args, *updateTODO.Info)
		argcount++
	}

	if updateTODO.Priority != nil {
		updateQuery += " Priority = $" + strconv.Itoa(argcount) + ","
		args = append(args, *updateTODO.Priority)
		argcount++
	}

	updateQuery = strings.TrimSuffix(updateQuery, ",")

	updateQuery += " WHERE TodoID = $" + strconv.Itoa(argcount) + ";"
	args = append(args, updateTODO.TodoID)

	updateStmt, err := db.Prepare(updateQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error preparing statement"})
		return
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error executing statement"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": updateTODO.TodoID + " updated"})
}

func handleDeleteTodo(db *sql.DB, ctx *gin.Context) {
	var err error
	todoID := ctx.Query("tid")

	del, err := db.Prepare("DELETE FROM todos WHERE TodoID = $1")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error preparing statement"})
		return
	}
	defer del.Close()

	_, err = del.Exec(todoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error executing statement"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": todoID + " Deleted"})
}
