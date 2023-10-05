package handlers

import (
	// "fmt"
	"net/http"
    "encoding/json"
    "context"
    "time"
    
	// "strconv"

    "go.mongodb.org/mongo-driver/bson"
    "github.com/julienschmidt/httprouter"
    "github.com/lucas-kern/tower-of-babel_server/app/model"
    requests "github.com/lucas-kern/tower-of-babel_server/app/model/requests"
    "github.com/lucas-kern/tower-of-babel_server/app/auth"
)

func (env *HandlerEnv) Bases(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        // id, err := strconv.Atoi(params.ByName("id"))
        // if err != nil {return};
        // json.NewEncoder(w).Encode(bases[id])
        WriteSuccessResponse(w, true)
}

// TODO finish this method
// Function to place a building in a base
func (env *HandlerEnv) PlaceBuilding(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    currUser := new(model.User)
    currBase := new(model.Base)
    var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    // How to access the "claims" object so the user properties from the auth token
    // TODO use this to pull user from DB for validations
    claims := r.Context().Value("claims").(*auth.SignedDetails)
    // Parse the request body to get building placement data
    var placementData requests.BuildingPlacement
    err := json.NewDecoder(r.Body).Decode(&placementData)
    if err != nil {
        // Handle JSON decoding error
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    //TODO have the user to check things like resources
    var userCollection model.Collection = env.database.GetUsers()
    err = userCollection.FindOne(currUser, ctx, bson.M{"ID": claims.Uid})

    var baseCollection model.Collection = env.database.GetBases()
    err = baseCollection.FindOne(currBase, ctx, bson.M{"ID": claims.Uid})
    building := model.NewBuilding(&placementData)



    // Perform validation checks (e.g., user permissions, available resources)
    // TODO ensure that user has the correct amount of resources. Will need to pull user from DB with user ID from claims and then ensure the info is correct.
    // validate that building is unlocked for user
    // validate building is not already placed
    if err := currBase.AddBuildingToBase(building); err != nil {
        // Validation failed, return an error response
        WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
        return
    }

    // Update the database with the new building placement or upgrade

    // Return a success response to the client
    WriteSuccessResponse(w, "Building placed/updated successfully")
}

// var bases = []model.Base{
// }
