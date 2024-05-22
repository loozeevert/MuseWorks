package handlers

import (
	"fmt"
	"net/http"
	"package/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// структура юзера для запроса редактирования
type EditUserReq struct {
	Username    string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

// функция получения UUID пользователя из jwt токена
func getUserUUIDFromToken(tokenString string) (uuid.UUID, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id := claims["id"].(string)
		return uuid.Parse(id)
	}

	return uuid.UUID{}, fmt.Errorf("Invalid token")
}

// хендлер изменения данных пользователя
func (h BaseHandler) EditUser(c *gin.Context) {
	var user EditUserReq
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	userUUID, err := getUserUUIDFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	dbUser := &db.User{
		UserId:      userUUID,
		Username:    user.Username,
		Description: &user.Description,
		Avatar:      &user.Avatar,
	}

	err = h.db.EditUserInfo(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
