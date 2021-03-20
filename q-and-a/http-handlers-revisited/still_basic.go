package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// User represents a person in our system.
type User struct {
	Name string
}

// UserService provides ways of working with users.
type UserService interface {
	Register(user User) (insertedID string, err error)
}

// UserServer provides a HTTP API for working with users.
type UserServer struct {
	service UserService
}

// NewUserServer creates a UserServer.
func NewUserServer(service UserService) *UserServer {
	return &UserServer{service: service}
}

// RegisterUser is a http handler for storing users.
func (u *UserServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode user payload: %v", err), http.StatusBadRequest)
		return
	}

	insertedID, err := u.service.Register(newUser)

	if err != nil {
		//todo: handle different kinds of errors differently
		http.Error(w, fmt.Sprintf("problem registering new user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, insertedID)
}

// MongoUserService provides storage functionality for Users.
type MongoUserService struct {
}

// NewMongoUserService creates a new MongoUserService managing connection pools etc probably!.
func NewMongoUserService() *MongoUserService {
	//todo: pass in DB URL as argument to this function
	//todo: connect to db, create a connection pool
	return &MongoUserService{}
}

// Register will store a user in mongo.
func (m MongoUserService) Register(user User) (insertedID string, err error) {
	// use m.mongoConnection to perform queries
	panic("implement me")
}

func main() {
	mongoService := NewMongoUserService()
	server := NewUserServer(mongoService)
	log.Fatal(http.ListenAndServe(":8000", http.HandlerFunc(server.RegisterUser)))
}
