package userHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	models "github.com/Anuragch/go_auth/models/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type UserHandler struct {
	log  *logrus.Logger
	repo models.UserRepository
}

type WelcomeHandler struct {
	name string
}

func NewWelcomeHandler() *WelcomeHandler {
	return &WelcomeHandler{"welcome_handler"}
}

func (w *WelcomeHandler) Welcome(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Handeling welcome ")
}

func NewUserHandler(log *logrus.Logger, repo models.UserRepository) *UserHandler {
	return &UserHandler{log, repo}
}

// New user signup handler
func (u *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	u.log.Println("Handling new user signup")

	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.log.Panic("Could not read request")
	}
	defer r.Body.Close()
	var user models.User
	err = json.Unmarshal(rBody, &user)

	ret := u.repo.AddUser(&user)
	if ret != false {
		u.log.Panic("Could not add user")
	}
	fmt.Println(user)
}

func (u *UserHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.log.Println("Getting all users")
	rw.Header().Set("Content-Type", "application/json")
	users := u.repo.GetAllUsers()
	res, err := json.Marshal(&users)
	if err != nil {
		panic(err)
	}
	rw.Write(res)
}

func (u *UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	u.log.Println("Getting one user")
	params := mux.Vars(r)
	userId := params["id"]
	rw.Header().Set("Content-Type", "application/json")
	fmt.Println("Getting user info for ", userId)
	user, err := u.repo.GetUser(userId)
	res, err := json.Marshal(&user)
	if err != nil {
		panic(err)
	}
	rw.Write(res)
}

func (u *UserHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) AuthenticateUser(rw http.ResponseWriter, r *http.Request) {

	u.log.Println("Authenticating user")

	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.log.Panic("Could not read request")
	}
	defer r.Body.Close()
	var userCredential models.Credentials
	err = json.Unmarshal(rBody, &userCredential)
	storedCredential, err := u.repo.GetUserCredential(userCredential.ID)

	if err != nil {
		u.log.Panic("Could not find user")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCredential.PasswordHash), []byte(userCredential.PasswordHash)); err != nil {
		// If the two passwords don't match, return a 401 status
		fmt.Println("Password mismatch ")
		rw.WriteHeader(http.StatusUnauthorized)
	}

	// Isue JWT token
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		ID: userCredential.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}
