package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"your_project_name/pkg/models"
	"your_project_name/pkg/utils"
)

var ctx = context.TODO()

// UserController handles user-related operations.
type UserController struct {
	collection *mongo.Collection
}
// NewUserController is a constructor for user controller
// NewUserController creates a new UserController instance.
func NewUserController(client *mongo.Client, dbName, collectionName string) *UserController {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserController{collection: collection}
}

// CreateUser creates a new user.
// Signup creates a new user.
func (uc *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// check if the user is admin to be able to create user
	adminId := os.Getenv("ADMIN_ID")
	if user.ID != adminId {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Hash the password before storing it
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword
	user.IsAdmin = false

	//check if the user with the same email exists
	existingUser := models.User{}
	err = uc.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User with this email already exists", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments{
		log.Println("Error getting user:", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
	}

	_, err := uc.collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
// Login handles user login.
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//parse the request body
	var loginUser models.LoginUser
	_ = json.NewDecoder(r.Body).Decode(&loginUser)

	//fetch the user from database
	var user models.User
	err := uc.collection.FindOne(ctx, bson.M{"email": loginUser.Email}).Decode(&user)
	if err != nil {
		log.Println("Error getting user:", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	//compare the hash password
	err = utils.ComparePassword(user.Password, loginUser.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	//generate jwt token
	token, err := utils.CreateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	json.NewEncoder(w).Encode(response)
}