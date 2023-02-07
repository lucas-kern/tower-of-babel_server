package handlers

import (
		"encoding/json"
    "context"
    "fmt"
		"log"

    "net/http"
    "time"

		"github.com/go-playground/validator/v10"
		"github.com/julienschmidt/httprouter"

    // "github.com/lucas-kern/tower-of-babel_server/app/server/database"

    "github.com/lucas-kern/tower-of-babel_server/app/auth"
    "github.com/lucas-kern/tower-of-babel_server/app/model"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    // "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
			log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
			msg = fmt.Sprintf("login or passowrd is incorrect")
			check = false
	}

	return check, msg
}

// Sign up allows a user with a unique email address to create an account and persists the account
func (env *HandlerEnv) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	var user model.User

	//TODO ensure that we are receiving the correct structure for this endpoint.
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Panic(err)
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
			log.Panic(err)
			return
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	if count > 0 {
			log.Println("error: this email already exists")
			return
	}

	//TODO ensure we are not just taking input, but are sanitizing it to improve security
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := auth.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			log.Println(msg)
			return
	}
	defer cancel()

	// TODO need to create a result object that will include a body to hold items we want returned to client
	json.NewEncoder(w).Encode(http.StatusOK)
}

//Login will allow a user to login to an account
func (env *HandlerEnv) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user model.User
	foundUser := new(model.User)

	//TODO we don't need to panic at every failed login rather inform user it is wrong
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Panic(err)
		return
	}

	err = userCollection.FindOne(foundUser, ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		  log.Panic(err)
		  return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
			log.Panic(msg)
			return
	}

	token, refreshToken, _ := auth.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

	auth.UpdateAllTokens(userCollection, token, refreshToken, foundUser.User_id)

	//TODO need to return the user information and tokens
	json.NewEncoder(w).Encode(http.StatusOK)
}
