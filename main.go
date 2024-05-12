package main

import (
	"package/db"
	"package/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
// основная исполняющая функция, здесь мы инициализируем БД, настраиваем маршруты запросов
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pool := db.DbStart(databaseURL)

	db := db.NewDB(pool)
	handler := handlers.NewBaseHandler(db)
	r := gin.Default()

	// здесь игнорим CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})
// роут получения инфы о пользователе
	r.GET("/getUserInfo", func(c *gin.Context) {
		handler.GetUserInfo(c)
	})
	v1 := r.Group("/auth")
	{
		// роут отправки подтверждающего письма на почту
		v1.POST("/sendMail", func(c *gin.Context) {
			handler.SendMail(c)
		})
		// роут регистрации
		v1.POST("/register", func(c *gin.Context) {
			handler.RegisterUser(c)
		})
		// роут входа
		v1.POST("/login", func(c *gin.Context) {
			handler.LoginUser(c)
		})
	}
//остальное базовые настройки сервера
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
