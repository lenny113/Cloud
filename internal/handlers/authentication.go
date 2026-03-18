package handlers

import (
	"crypto/md5"   //for generarting hash to create api key
	"encoding/hex" //for converting md5 hash to string
	"encoding/json"
	"fmt"
	"net/http"
	"time" //time of creating api key and for generating unique api key based partly on time hash
)

type Login struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type key struct {
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt"`
}

func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //if not a POST request, return method not allowed
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body of the request into a Login struct
	var login Login
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate the input (check if name and email are not empty, and if email contains @)
	if login.Email == "" || login.Name == "" {
		http.Error(w, "Missing name or email", http.StatusBadRequest)
		return
	}

	//check if email contains @, if not, return bad request
	isAtValid := false
	for i := 0; i < len(login.Email); i++ {
		if login.Email[i] == '@' {
			isAtValid = true
			break
		}
	}
	if !isAtValid {
		http.Error(w, "Invalid email format, no @ found", http.StatusBadRequest)
		return
	}

	// Generate API key for the user
	var createAPI, timeCreateApi string
	//maxAttempt should be a const, but I want to ask team member or TA for what is a reasonable number
	//max attempt is used to avoid loop with checking api key duplicate!
	maxTestAttempts := 10
	for i := 0; i < maxTestAttempts; i++ {
		createAPI, timeCreateApi = createAPIKey(login.Email)
		if !isAPIKeyUsed(createAPI) {
			break
		}
		//this would only be in the senario where some number breakes the md5 hash or full storage in database
		//then it would be an infinte loop of generating api keys, this should not happen, so we want to break this after a certain number of attempts and log it
		if i == maxTestAttempts-1 {
			fmt.Printf("Failed to generate a unique API key after %d attempts\n", i+1)
			http.Error(w, "Failed to generate a unique API key after several attempts", http.StatusLoopDetected)
			return
		}
	}

	//formatting the response to the user, which includes the generated API key and the time of API creation
	key := key{
		Key:       createAPI,
		CreatedAt: timeCreateApi,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(key); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println("API key generated and sent to user")

}

func createAPIKey(email string) (string, string) {

	timeCreateApi := time.Now().Format("20060102 15:04") //this is the format specified in the assignement for time
	fmt.Println("Time of API creation:", timeCreateApi)  //maybe logg this

	startOfUserAPI := "sk-envdash-"

	// Generate hash of email + current time using md5
	//this will be used as the api key for the user, and will be unique for each registration
	//unique even, even same email cant make same key, because of the time component
	hash := md5.Sum([]byte(email + time.Now().String()))
	hashString := hex.EncodeToString(hash[:])
	createAPI := startOfUserAPI + hashString

	// Display the character
	fmt.Println("createAPI:", createAPI)

	return createAPI, timeCreateApi
}

func isAPIKeyUsed(key string) bool {
	//check if api key is already in use, if so, generate a new one (this is very unlikely, but we want to be sure)
	usedAPIKey := "sk-envdash-4597dc2d89e56c8e0cde3d3b9f42bdfa" // simulates a used api key

	//here we want to check if the generated API key were made or not, if no api key were made, send true
	if key == "" {
		//Logging not implemented, but needs to be logged because this should not happen
		//log.Println("Generated API key is empty")
		return true
	}
	//if only the start of the API key is generated, send true, because this means that the random string part is not generated yet, and this is not a valid API key
	if key == "sk-envdash-" {
		//logging waiting for implementation
		//log.Println("Generated API key is incomplete")
		return true
	}

	//when firebase is added this will be where api key is looking for duplicates
	if key == usedAPIKey {
		//logging waiting for implementation
		//log.Println("Generated API key is already in use")
		return true
	}
	return false

}
