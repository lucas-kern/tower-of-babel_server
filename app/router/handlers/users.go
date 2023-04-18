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
			msg = fmt.Sprintf("The username or password is incorrect")
			check = false
	}

	return check, msg
}

// Sign up allows a user with a unique email address to create an account and persists the account
// TODO see what code is repeated with login and make external function for it
func (env *HandlerEnv) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	var user model.User
	var clientUser *model.ClientUser
	// TODO make this method faster if possible
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 422, "There was an error with the client request")
		return
	}

	err = validate.Struct(user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 400, "There was an error with user validation")
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
			log.Println(err)
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}

	if count > 0 {
		log.Println("error: this email already exists")
		// TODO update with email responses like described here
		// https://security.stackexchange.com/questions/51856/e-mail-already-in-use-exploit
		WriteErrorResponse(w, 401, "There was an error registering this account")
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password

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
			log.Println("User item was not created")
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}
	defer cancel()

	WriteSuccessResponse(w, "Account created successfully")
}

//Login will allow a user to login to an account
func (env *HandlerEnv) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userCollection model.Collection = env.database.GetUsers()
	var user model.User
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	foundUser := new(model.User)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 422, "There was an error with the client request")
		return
	}

	err = validate.Struct(user)
	if err != nil {
		log.Println(err)
		WriteErrorResponse(w, 400, "There was an error with user validation")
		return
	}

	err = userCollection.FindOne(foundUser, ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
			log.Println(err)
			WriteErrorResponse(w, 502, "There was an error connecting with the server")
			return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		log.Println(msg)
		WriteErrorResponse(w, 401, "The username or password is incorrect")
		return
	}

	token, refreshToken, _ := auth.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

	auth.UpdateAllTokens(userCollection, token, refreshToken, foundUser.User_id)

	clientUser := model.NewUser(&user)

	WriteSuccessResponse(w, clientUser)
}
