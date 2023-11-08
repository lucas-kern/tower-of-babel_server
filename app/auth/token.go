package auth

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/lucas-kern/tower-of-babel_server/app/model"

    jwt "github.com/dgrijalva/jwt-go"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails
type SignedDetails struct {
	Email      string
	FirstName string
	LastName  string
	Uid        string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens generates both the detailed token and refresh token
func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
    claims := &SignedDetails{
        Email:      email,
        FirstName: firstName,
        LastName:  lastName,
        Uid:        uid,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
        },
    }

    refreshClaims := &SignedDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
        },
    }

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

    if err != nil {
        log.Panic(err)
        return
    }

    return token, refreshToken, err
}

//ValidateToken validates the jwt token
func ValidateToken(userCollection model.Collection, signedToken string) (claims *SignedDetails, msg string) {
	foundUser := new(model.User)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	token, err := jwt.ParseWithClaims(
			signedToken,
			&SignedDetails{},
			func(token *jwt.Token) (interface{}, error) {
					return []byte(SECRET_KEY), nil
			},
	)

	if err != nil {
			msg = err.Error()
			return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
			msg = fmt.Sprintf("the token is invalid")
			msg = err.Error()
			return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
			msg = fmt.Sprintf("token is expired")
			msg = err.Error()
			return
	}
	
	//Check that user from token matches the one in DB
	Id, err := primitive.ObjectIDFromHex(claims.Uid)
	err = userCollection.FindOne(foundUser, ctx, bson.M{"_id": Id})
	if err != nil {
		msg = "error fetching user from the database"
		return
	}

	if foundUser == nil {
		msg = "user not found in the database"
		return
	}

	return claims, msg
}

//UpdateAllTokens renews the user tokens when they login
func UpdateAllTokens(userCollection model.Collection, signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshToken", signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", UpdatedAt})

	upsert := true
	Id, err := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": Id}
	opt := options.UpdateOptions{
			Upsert: &upsert,
	}

	_, err = userCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
					{"$set", updateObj},
			},
			&opt,
	)
	defer cancel()

	if err != nil {
			log.Panic(err)
			return
	}

	return
}
