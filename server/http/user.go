package http

import (
	"github.com/gin-gonic/gin"
	"github.com/huizm/go-chatroom/logic"
	"github.com/huizm/go-chatroom/model"
	"net/http"
)

func searchUserByUsername(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "empty username",
			"error":   "",
		})
		return

	}

	inq := &model.User{Username: username}
	data, err := model.SearchUser(inq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "search user error",
			"error":   err.Error(),
		})
		return
	}

	resp := data[0].ID

	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func createUser(ctx *gin.Context) {
	data := &model.User{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "bad create user request",
			"error":   err.Error(),
		})
		return
	}

	if err := model.CreateUser(data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "create user error",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"username": data.Username,
	})
}

func auth(ctx *gin.Context) {
	data := &model.User{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "bad auth request",
			"error":   err.Error(),
		})
		return
	}

	inq := &model.User{Username: data.Username}
	users, err := model.SearchUser(inq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
			"error":   err.Error(),
		})
		return
	}

	if len(users) == 0 || users[0].Password != data.Password {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "incorrect username or password",
			"error":   "",
		})
		return
	}

	token, err := logic.GenerateToken(users[0])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
