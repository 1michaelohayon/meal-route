package types

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email             string             `bson:"email" json:"email"`
	FullName          string             `bson:"fullName" json:"fullName"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	Admin             bool               `bson:"admin" json:"admin"`
}

const (
	bcryptCost  = 10
	minFullName = 3
	maxFullName = 75
	minPass     = 7
	maxPass     = 40
)

type NewUser struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
}

func (params *NewUser) CreateUser() (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FullName:          params.FullName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		Admin:             false,
	}, nil
}

/* TODO not allow duplicate emails  */
func (params *NewUser) Validate() map[string]string {
	validEmail := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

	errMap := map[string]string{}
	if len(params.FullName) < minFullName || len(params.FullName) > maxFullName {
		errMap["fullName error"] = fmt.Sprintf("fullName must be between %d and %d", minFullName, maxFullName)
	}
	if !validEmail.MatchString(params.Email) {
		errMap["email error"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	if len(params.Password) < minPass || len(params.Password) > maxPass {
		errMap["password error"] = fmt.Sprintf("password must be between %d and %d", minPass, maxPass)
	}
	if len(errMap) != 0 {
		return errMap
	}

	return nil
}

func (user *User) GenerateToken() string {
	now := time.Now()
	expires := now.Add(time.Hour * 1).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if len(secret) == 0 {
		log.Fatal("Secret is empty")
	}
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
