package handlers

import (
	// "fmt"
	"net/http"
	// "encoding/json"
	// "strconv"

	"github.com/julienschmidt/httprouter"
	// "github.com/lucas-kern/tower-of-babel_server/app/model"
)

func (env *HandlerEnv) Bases(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
        // id, err := strconv.Atoi(params.ByName("id"))
        // if err != nil {return};
        // json.NewEncoder(w).Encode(bases[id])
        WriteSuccessResponse(w, true)
}

// Function to place a building in a base
func (env *HandlerEnv) PlaceBuilding(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // How to access the "claims" object so the user properties from the auth token
    // claims := r.Context().Value("claims").(*auth.SignedDetails)
    // // Parse the request body to get building placement data
    // var placementData model.BuildingPlacement
    // err := json.NewDecoder(r.Body).Decode(&placementData)
    // if err != nil {
    //     // Handle JSON decoding error
    //     http.Error(w, "Error decoding request body", http.StatusBadRequest)
    //     return
    // }

    // // Perform validation checks (e.g., user permissions, available resources)

    // // Update the database with the new building placement or upgrade

    // // Return a success or error response to the client
    // if validationPassed {
    //     WriteSuccessResponse(w, "Building placed/updated successfully")
    // } else {
    //     WriteErrorResponse(w, http.StatusForbidden, "Validation failed or error occurred")
    // }
}

// var bases = []model.Base{
// }
