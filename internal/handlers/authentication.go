package handlers

import (
	"crypto/md5"   //for generarting hash to create api key
	"encoding/hex" //for converting md5 hash to string
	"fmt"
	"time" //for generating unique api key based partly on time hash
)

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
