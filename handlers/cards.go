package handlers

import (
	"fmt"
	"net/http"
	"package/db"

	"github.com/gin-gonic/gin"
)

// хендлер создания карточки
func (h BaseHandler) CreateCard(c *gin.Context) {
	var card *db.Card
	err := c.BindJSON(&card)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	userUUID, err := getUserUUIDFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	card.User = &db.User{UserId: userUUID}

	card, err = h.db.CreateCard(card)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, card)
}

// хендлер удаления карточки
func (h BaseHandler) DeleteCard(c *gin.Context) {
	var card *db.Card
	err := c.BindJSON(&card)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	userUUID, err := getUserUUIDFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	card.User = &db.User{UserId: userUUID}

	_, err = h.db.DeleteCard(card)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// хендлер получения всех карточек
func (h BaseHandler) GetAllCards(c *gin.Context) {
	cards, err := h.db.GetAllCards()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, cards)
}
