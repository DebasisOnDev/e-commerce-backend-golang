package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//some test cases failed as they are written by chatgpt

func TestPostUser(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserFromParams{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "john_doe",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}

}

func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/", userHandler.HandleGetUsers)

	// Add some users to the database for testing
	usersToAdd := []types.CreateUserFromParams{
		{Email: "user1@example.com", FirstName: "Alice", LastName: "Johnson", Password: "alice123"},
		{Email: "user2@example.com", FirstName: "Bob", LastName: "Smith", Password: "bob456"},
	}

	for _, params := range usersToAdd {
		user, err := types.NewUserFromParams(params)
		if err != nil {
			t.Error(err)
		}
		_, err = tdb.UserStore.InsertUser(context.TODO(), user)
		if err != nil {
			t.Error(err)
		}
	}

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Error(err)
	}

	var users []types.User
	json.NewDecoder(resp.Body).Decode(&users)

	// Check if the number of returned users matches the number added to the database
	if len(users) != len(usersToAdd) {
		t.Errorf("expected %d users but got %d", len(usersToAdd), len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/:id", userHandler.HandleGetUser)

	// Add a user to the database for testing
	params := types.CreateUserFromParams{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "john_doe",
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Error(err)
	}
	insertedUser, err := tdb.UserStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Error(err)
	}

	// Retrieve the user by ID
	url := fmt.Sprintf("/%s", insertedUser.ID)
	resp, err := app.Test(httptest.NewRequest("GET", url, nil))
	if err != nil {
		t.Error(err)
	}

	var returnedUser types.User
	json.NewDecoder(resp.Body).Decode(&returnedUser)

	// Check if the returned user matches the inserted user
	if returnedUser.ID != insertedUser.ID {
		t.Errorf("expected user ID %s but got %s", insertedUser.ID, returnedUser.ID)
	}
}

func TestDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Delete("/:id", userHandler.HandleDeleteUser)

	// Add a user to the database for testing
	params := types.CreateUserFromParams{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "john_doe",
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Error(err)
	}
	insertedUser, err := tdb.UserStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Error(err)
	}

	// Delete the user by ID
	url := fmt.Sprintf("/%s", insertedUser.ID)
	resp, err := app.Test(httptest.NewRequest("DELETE", url, nil))
	if err != nil {
		t.Error(err)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	InsId := insertedUser.ID.Hex()

	// Check if the response indicates successful deletion
	if result["Deleted"] != InsId {
		t.Errorf("expected deleted user ID %s but got %s", insertedUser.ID, result["Deleted"])
	}
}

func TestUpdateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Put("/:id", userHandler.HandlePutUser)

	// Add a user to the database for testing
	params := types.CreateUserFromParams{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "john_doe",
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Error(err)
	}
	insertedUser, err := tdb.UserStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Error(err)
	}

	// Update the user by ID
	url := fmt.Sprintf("/%s", insertedUser.ID)
	updateParams := types.UpdateUserParams{
		FirstName: "UpdatedFirstName",
		LastName:  "UpdatedLastName",
	}

	b, _ := json.Marshal(updateParams)
	req := httptest.NewRequest("PUT", url, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	// Check if the response indicates successful update
	if result["updated"] != insertedUser.ID.Hex() {
		t.Errorf("expected updated user ID %s but got %s", insertedUser.ID, result["updated"])
	}

	// Retrieve the updated user and check if the fields are updated
	updatedUser, err := tdb.UserStore.GetUserByID(context.TODO(), insertedUser.ID.Hex())
	if err != nil {
		t.Error(err)
	}

	if updatedUser.FirstName != updateParams.FirstName {
		t.Errorf("expected updated first name %s but got %s", updateParams.FirstName, updatedUser.FirstName)
	}
	if updatedUser.LastName != updateParams.LastName {
		t.Errorf("expected updated last name %s but got %s", updateParams.LastName, updatedUser.LastName)
	}

}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	dburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	userStore := db.NewMongoUserStore(client)

	return &testdb{
		client:    client,
		UserStore: userStore,
	}
}

type testdb struct {
	client *mongo.Client
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	dbname := os.Getenv("MONGO_DB_NAME")
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
