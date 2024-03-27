package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleNewUser(db *sql.DB, ctx *gin.Context) {
	type NewUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var newUser NewUser
	var err error

	if err = ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := db.Prepare("INSERT INTO users (UserID, Username, Password) VALUES ($1, $2, $3);")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}
	defer user.Close()

	userID := uuid.New().String()

	_, err = user.Exec(userID, newUser.Username, newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "User Created", "data": gin.H{"user_id": userID}})
}
