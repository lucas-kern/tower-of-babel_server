package handlers

import (
    "context"
    "fmt"
    "log"

    "net/http"
    "time"

		"github.com/go-playground/validator/v10"
		"github.com/julienschmidt/httprouter"

    "github.com/lucas-kern/tower-of-babel_server/app/server/database"

    "github.com/lucas-kern/tower-of-babel_server/app/auth"
    "github.com/lucas-kern/tower-of-babel_server/app/model"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetUsers()
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

func SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user model.User

	validationErr := validate.Struct(user)
	if validationErr != nil {
			log.Panic(validationErr)
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

	count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	defer cancel()
	if err != nil {
			log.Panic(err)
			return
	}

	if count > 0 {
			log.Println("error: this email or phone number already exists")
			return
	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := auth.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			log.Println(msg)
			return
	}
	defer cancel()

	json.NewEncoder(w).Encode(http.StatusOK, resultInsertionNumber)

}

//Login is the api used to get a single user
//TODO update to use the current router
// func Login() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 			var user models.User
// 			var foundUser models.User

// 			if err := c.BindJSON(&user); err != nil {
// 					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 					return
// 			}

// 			err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
// 			defer cancel()
// 			if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
// 					return
// 			}

// 			passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
// 			defer cancel()
// 			if passwordIsValid != true {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 					return
// 			}

// 			token, refreshToken, _ := auth.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

// 			auth.UpdateAllTokens(token, refreshToken, foundUser.User_id)

// 			c.JSON(http.StatusOK, foundUser)

// 	}
// }
