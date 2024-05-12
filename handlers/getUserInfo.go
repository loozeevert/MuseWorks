package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
//хендлер получения юзера по токену
func (h BaseHandler) GetUserInfo(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует токен авторизации"})
		return
	}

	token, err := parseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}
	userID, err := uuid.Parse(token.Claims.(jwt.MapClaims)["id"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при извлечении ID пользователя"})
		return
	}

	user, err := h.db.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ошибка при запросе данных из базы данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
