package db

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// структура базы данный где хранится пулы подключения к базе данных
type DB struct {
	pool *pgxpool.Pool
}

// структура юзера
type User struct {
	UserId      uuid.UUID `json:"id"`
	Username    string    `json:"name,omitempty"         binding:"required"`
	Email       string    `json:"email,omitempty"        binding:"required"`
	Password    string    `json:"password,omitempty"     binding:"required,min=8"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	ConfirmCode int       `json:"confirmCode,omitempty"`
}

// структура передаваемых значений при логине
type UserLoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// структура токена
type Token struct {
	TokenString string `json:"accessToken"`
}

// структура Email
type UserEmailData struct {
	Email string `json:"email" binding:"required"`
}

// функция создания нового подключения к бд
func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

// функция хеширования пароля
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// функция старта подключения к бд (используется в main)
func DbStart(baseUrl string) *pgxpool.Pool {
	urlExample := baseUrl
	dbpool, err := pgxpool.New(context.Background(), string(urlExample))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	return dbpool
}

// функция которая является методом структуры DB позволяет сделать запрос в бд для добавления пользователя после регистрации
func (db DB) RegisterUser(userData User) (string, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return "Проблема с установкой соединения", fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	userData.UserId = uuid.New()
	password, hashErr := hashPassword(userData.Password)
	if hashErr != nil {
		return "Проблема с хешированием", fmt.Errorf("unable to hashPass: %v", hashErr)
	}

	err = conn.QueryRow(context.Background(),
		`INSERT INTO users(userid, username, email, password) VALUES ($1, $2, $3, $4) RETURNING userid`,
		userData.UserId, userData.Username, userData.Email, password).Scan(&userData.UserId)
	if err != nil {
		return "Проблема с запросом в базу данных", fmt.Errorf("unable to INSERT: %v", err)
	}

	return "Вы успешно зарегистрировались", nil
}

// функция которая является методом структуры DB позволяет сделать запрос в бд для получения инфы о пользователе
func (db DB) GetUserByEmail(email string) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var user User
	err = conn.QueryRow(context.Background(), "SELECT userid, username, email, password FROM users WHERE email = $1", email).
		Scan(&user.UserId, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user: %v", err)
	}

	return &user, err
}

// опциональная функция, которая нужна для проверки наличия юзера в БД
func (db DB) userExists(userID uuid.UUID) (bool, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE userid = $1)", userID).
		Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return exists, err
}

// функция которая является методом структуры DB позволяет сделать запрос в бд для получения инфы о юзере
func (db DB) GetUserInfo(userID uuid.UUID) (User, error) {
	exists, err := db.userExists(userID)
	if err != nil {
		return User{}, err
	}

	if !exists {
		return User{}, fmt.Errorf("пользователь с ID %s не существует", userID.String())
	}

	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return User{}, fmt.Errorf("невозможно получить соединение с базой данных: %v", err)
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		"SELECT userid, username, email, description, avatar FROM users WHERE userid = $1", userID)

	var user User
	err = row.Scan(&user.UserId, &user.Username, &user.Email, &user.Description, &user.Avatar)
	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, fmt.Errorf("пользователь с ID %s не найден", userID.String())
		}
		return User{}, fmt.Errorf("невозможно прочитать данные из базы данных: %v", err)
	}

	return user, err
}

// функция которая проверяет при отправке на почту письма есть ли юзер в БД
func (db DB) UserExistsByEmail(email string) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// функция которая изменяет данные пользователя
func (db DB) EditUserInfo(user *User) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("невозможно получить соединение с базой данных: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		"UPDATE users set (username, description, avatar) = ($1, $2, $3) WHERE userid = $4",
		user.Username, user.Description, user.Avatar, user.UserId)

	return err
}
