package authcontroller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hisyam99/go-jwt-mux/config"
	"github.com/hisyam99/go-jwt-mux/helper"
	"github.com/hisyam99/go-jwt-mux/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	defer r.Body.Close()

	var user models.User
	collection := models.DB.Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"username": userInput.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
		} else {
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
		return
	}

	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "login berhasil"})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	defer r.Body.Close()

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	collection := models.DB.Collection("users")
	_, err := collection.InsertOne(context.TODO(), userInput)
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "logout berhasil"})
}
