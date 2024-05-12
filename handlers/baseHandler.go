package handlers

import (
	"package/db"
	"os"
	"sync"

	"github.com/dgrijalva/jwt-go"
)
// основной хендлер в котором хранятся пулы подключения к бд и коды подтверждения юзеров
type UserGet struct {
	Parce []db.User `json:"parce"`
}

type BaseHandler struct {
	db   *db.DB
	Code map[string]*db.User
	mu   sync.Mutex
}
//функция которая будет передавать в другие хендлеры пулы подключения и хештаблицу с кодами
func NewBaseHandler(pool *db.DB) *BaseHandler {
	return &BaseHandler{
		db:   pool,
		Code: make(map[string]*db.User),
	}
}
//функция парса токена
func parseToken(tokenString string) (*jwt.Token, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

