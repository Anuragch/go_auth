package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Anuragch/go_auth/errors"
	models "github.com/Anuragch/go_auth/models/user"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDB struct {
	Instance *mongo.Database
	Name     string
}

func (userDb *UserDB) GetAllUsers() []*models.User {
	var userList []*models.User
	collection := userDb.Instance.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})

	// Find() method raised an error
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(ctx)

		// If the API call was a success
	} else {
		// iterate over docs using Next()
		for cursor.Next(ctx) {

			var result models.User
			err := cursor.Decode(&result)

			// If there is a cursor.Decode error
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
				os.Exit(1)

				// If there are no cursor.Decode errors
			}
			userList = append(userList, &result)
		}
	}
	return userList
}

func (userDb *UserDB) GetUser(id string) (*models.User, error) {
	collection := userDb.Instance.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	user := models.User{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, &errors.DBError{"User not found"}
	}
	return &user, nil
}

func (userDb *UserDB) GetUserCredential(id string) (*models.Credentials, error) {
	collection := userDb.Instance.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	user := models.User{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, &errors.DBError{"User not found"}
	}
	userCredential := models.Credentials{ID: user.ID, PasswordHash: user.PasswordHash}
	return &userCredential, nil
}

func (userDb *UserDB) AddUser(user *models.User) bool {
	fmt.Println("Attempting AddUser")
	collection := userDb.Instance.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 8)
	user.PasswordHash = string(hashedPassword)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("Could not insert user")
		panic(err)
	}
	fmt.Println("Inserted:", result.InsertedID)
	fmt.Println("Adding user")
	return false
}

// handle error
func NewUserRepo(db *mongo.Database, dbName string) models.UserRepository {
	return &UserDB{Instance: db, Name: dbName}
}
