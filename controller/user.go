package controller

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"strconv"

	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
	"github.com/guptaaashutosh/gqlgen-prac/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateNewUser
func CreateNewUser(ctx context.Context, DB *pgxpool.Pool, newUserData model.NewUser) (*model.User, error) {
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
	user := &model.User{
		ID:       int(randNumber.Int64()),
		Username: newUserData.Username,
		Email:    newUserData.Email,
		Password: newUserData.Password,
	}
	_, err := DB.Exec(ctx, "INSERT INTO bookuser (id, username, email, password) VALUES ($1,$2,$3,$4)", user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, nil
}

// getUsers
func GetUsers(ctx context.Context, DB *pgxpool.Pool) ([]*model.User, error) {
	var users []*model.User
	rows, err := DB.Query(ctx, "select id, username, email, password from bookuser")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		var user model.User
		rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		users = append(users, &user)
	}
	return users, nil
}

func GetUserWithUserId(ctx context.Context, DB *pgxpool.Pool, userId int) (*model.User, error) {
	var user model.User
	userRows := DB.QueryRow(ctx, "select id, username, email, password from bookuser where id=$1", userId)
	userRows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return &user, nil
}

func CheckUserExistence(ctx context.Context, DB *pgxpool.Pool, email string, password string) (int, bool, error) {
	var id int
	rows := DB.QueryRow(ctx, "select id from bookuser where email=$1 and password=$2", email, password)
	err := rows.Scan(&id)
	if err != nil {
		return id, false, err
	}
	return id, true, nil
}

func LoginUserController(ctx context.Context, DB *pgxpool.Pool, email string, password string) (*model.LoginDetails, error) {
	var LoginDetails model.LoginDetails
	var loggedInStatus = true
	var userId, isExists, err = CheckUserExistence(ctx, DB, email, password)

	if isExists {
		//generate token here
		loggedInToken, tokenError := utils.GenerateToken(strconv.Itoa(userId))
		if tokenError != nil {
			return nil, tokenError
		}
		LoginDetails.IsLoggedIn = &loggedInStatus
		LoginDetails.Token = loggedInToken
		return &LoginDetails, err
	}
	// token should store in cookie

	return &LoginDetails, nil
}
