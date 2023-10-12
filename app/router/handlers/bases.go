package handlers

import (
	"net/http"
    "encoding/json"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"

    "github.com/julienschmidt/httprouter"
    "github.com/lucas-kern/tower-of-babel_server/app/model"
    requests "github.com/lucas-kern/tower-of-babel_server/app/model/requests"
    "github.com/lucas-kern/tower-of-babel_server/app/auth"
)

// Function to place a building in a base
func (env *HandlerEnv) PlaceBuilding(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    currUser := new(model.User)
    currBase := new(model.Base)
    var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // How to access the "claims" object so the user properties from the auth token
    claims := r.Context().Value("claims").(*auth.SignedDetails)

    // Parse the request body to get building placement data
    var placementData requests.BuildingPlacement
    err := json.NewDecoder(r.Body).Decode(&placementData)
    if err != nil {
        // Handle JSON decoding error
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    // Convert the hexadecimal string to an ObjectID
    objectID, err := primitive.ObjectIDFromHex(claims.Uid)
    if err != nil {
        // Handle the error if the hex string is not a valid ObjectID
        panic(err)
    }

    //TODO have the user to check things like resources
    var userCollection model.Collection = env.database.GetUsers()
    err = userCollection.FindOne(currUser, ctx, bson.M{"ID": objectID})

    var baseCollection model.Collection = env.database.GetBases()
    err = baseCollection.FindOne(currBase, ctx, bson.M{"owner": objectID})
    building := model.NewBuilding(&placementData)

    // Perform validation checks (e.g., user permissions, available resources)
    // TODO ensure that user has the correct amount of resources. Will need to pull user from DB with user ID from claims and then ensure the info is correct.
    // validate that building is unlocked for user

    // Need to persist the base after the building is placed
    if err := currBase.AddBuildingToBase(building); err != nil {
        // Validation failed, return an error response
        WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
        return
    }

    // Define the filter to identify the base document by its ID (assuming it's a unique identifier).
    filter := bson.M{"_id": currBase.ID}

    // Define the update operation to set the "Grid" field in the base document.
    update := bson.M{
        "$set": bson.M{
            "Grid": currBase.Grid,
            "Buildings": currBase.Buildings,
            // Add any other fields you need to update here
        },
    }

    // Specify update options (optional). For example, you can enable upsert or specify additional options.
    options := options.Update().SetUpsert(false)

    // Perform the update operation on the "bases" collection.
    result, err := baseCollection.UpdateOne(ctx, filter, update, options)
    if err != nil {
        // Handle the error if the update operation fails
        WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update base: "+err.Error())
        return
    }

    if result.MatchedCount == 0 {
        // Handle the case where no documents were matched by the filter
        WriteErrorResponse(w, http.StatusNotFound, "No matching base found for the update")
        return
    }
    
    if result.ModifiedCount == 0 {
        // Handle the case where the update didn't modify any documents
        WriteErrorResponse(w, http.StatusNoContent, "No changes were made during the update")
        return
    }

    // Return a success response to the client
    WriteSuccessResponse(w, currBase)
}
