package store

import (
	model "assignment-2/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FireStore wraps the Firestore client and provides methods
// for interacting with data in firestore
type FireStore struct {
	client *firestore.Client
}

// NewFirestoreStore initializes a new firestore instance with the
// provided Firestore client
func NewFirestoreStore(client *firestore.Client) *FireStore {
	return &FireStore{client: client}
}

// CreateRegistration stores a new registration under an API key.
// A document id is generated automatically by firestore and assigned to the registratrion
//
// Returns:
// - string: generated document ID
// - error:  if the operation fails
func (f *FireStore) CreateRegistration(ctx context.Context, apiKey string, reg model.Registration) (string, error) {

	// Navigate to user_registrations subcollection
	col := f.client.Collection("registrations").Doc(apiKey).Collection("user_registrations")

	//generate new doc reference per auto-ID
	doc := col.NewDoc()

	//assign generated id to the document
	reg.ID = doc.ID

	//store registration in firestore
	_, err := doc.Set(ctx, reg)
	if err != nil {
		return "", err
	}

	return doc.ID, nil
}

// GetRegistration retrieves a single registration by ID
//
// Returns:
// - *model.Registration: the requested registration
// - error: if the document does not exist or decoding fails
func (f *FireStore) GetRegistration(ctx context.Context, apiKey string, id string) (*model.Registration, error) {

	doc, err := f.client.Collection("registrations").
		Doc(apiKey).
		Collection("user_registrations").
		Doc(id).
		Get(ctx)
	if err != nil {
		return nil, err
	}

	var reg model.Registration
	if err := doc.DataTo(&reg); err != nil {
		return nil, err
	}

	// Include the document ID
	reg.ID = doc.Ref.ID

	return &reg, nil
}

/*
You may want to use this function to check is API key gives acces.

Checks if api key exists
Takes api key, hashes it, checks database if it exists.
This will be used when authenticating incomming api requests, if it exists returns true

This method is part of the Store struct, which holds the Firestore client.

@see			-hashAPIKey() for hashing implementation

@param ctx 		-keeping track of firestore connection(timeout etc)
@apiKey			-the key you want to check
@return bool	-if api key exists:true, if not in Firestore:false
*/
func (f *FireStore) ApiKeyExists(ctx context.Context, apiKey string) bool {
	hashedApiKey := hashAPIKey(apiKey)
	_, err := f.client.
		Collection("all_api_keys").
		Doc(hashedApiKey).
		Get(ctx)

	if err != nil {
		//api key cant be found e.g it is unique
		return false
	}

	//This api key exists!
	return true
}

/*
Storess API
Apis are currently stored in two different ways:

	1: All apis stored in one collection ass documents
			-Data stored: "time of creation" and what email used
	2: All users (email, addresses) have nested collection storing each api key
		These are the same api keys, stored in different ways
		This is donne for effecient lookup (if we letssay have 1 million users this would still work)
			-Data stored: "time of creation" and "name of api key"

# This function hashes api key so no clairtext api key is stored on server

This method is part of the Store struct, which holds the Firestore client.

@see			-hashAPIKey() for hashing implementation

@param ctx 		- keeping track of firestore connection(timeout etc)
@param reg 		- struct of all data that we want to store (api key gets hashed)
@return error 	- if anny errors cam when storing api key in firestore, if nil, the keys were stored!
*/
func (f *FireStore) CreateApiStorage(ctx context.Context, reg model.Authentication) error {
	//first hashes api key generated:
	hashedApiKey := hashAPIKey(reg.ApiKey)
	//setts api
	AllApi := f.client.Collection("all_api_keys").Doc(hashedApiKey)
	_, err := AllApi.Set(ctx, map[string]interface{}{
		"time of creation": reg.CreatedAt,
		"user":             reg.Email,
	})

	emailDoc := f.client.Collection("authentication_info").Doc(reg.Email)
	//creating nested api key structure
	EmailApiDoc := emailDoc.Collection("api_keys").Doc(hashedApiKey)

	_, err = EmailApiDoc.Set(ctx, map[string]interface{}{
		"time of creation": reg.CreatedAt,
		"name of api key":  reg.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

/*
Counts how manny Api's the speccified user have (checks Firestore)
and return the number of Apis's and anny errors if appropirate
If error, return 0 apis registerd to this user.

This method is part of the Store struct, which holds the Firestore client.

@param ctx	-keeping track of firestore connection(timeout etc)
@param email-The email you want to check, email is the user all apis are registerd under
@return int	-Return if anny, how manny Api's this user have registerd in Firestore, 0 if error
return error-Returns anny error and dont complete the function
*/
func (h *FireStore) CountApiPerUser(ctx context.Context, email string) (int, error) {
	//getting info about spesific email
	EmailDoc := h.client.Collection("authentication_info").Doc(email)
	//seeing how manny api keys that user hve
	ApiKeyDoc := EmailDoc.Collection("api_keys")

	doc, err := ApiKeyDoc.Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	//returns length of
	return len(doc), nil
}

/*
Deletes Api stored in Firestore. Deletes both places where Api is stored (global storage and per user)
First extract what email(user) this is api is registerd to, then delete in global storage (All_api_Keys)
Then delete this exact api from user.

This function don't delete user from database if this is the last api. This is because we
want to keep our user stored. We may want to enhance the functionality, and want to link maybe some other
information about this user

This method is part of the Store struct, which holds the Firestore client.

@see			-hashAPIKey() for hashing implementation

@param ctx		-keeping track of firestore connection(timeout etc)
@param apiKey	-api key from the user
@return error	-returns error if something goes wrong, example: wrong format stored in Firestore
*/
func (f *FireStore) DeleteAPIkey(ctx context.Context, apiKey string) error {
	apiKeyHashed := hashAPIKey(apiKey)
	docRef := f.client.Collection("all_api_keys").Doc(apiKeyHashed)

	// check if exists
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return err
	}

	//Finds mail to this user that this api is registerd under
	data := docSnap.Data()

	userMail, ok := data["user"].(string)
	if !ok {
		//TODO: log this in logg file
		return fmt.Errorf("Cant get email, user field missing or not a string (Firestore)")
	}

	//delete api under "ALL_API_KEYS"
	_, err = docRef.Delete(ctx)
	if err != nil {
		return err
	}

	//now goes to right user, and deletes that API key:
	userDoc := f.client.Collection("authentication_info").Doc(userMail)

	nestedDocRef := userDoc.Collection("api_keys").Doc(apiKeyHashed)

	_, err = nestedDocRef.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

/*
Hashes API key
This is done BEFORE being stored in database
Use this function when checking api key in database, since all are hashed
Uses sha 256 hash. Stores as string.

@param apiKeyUnhashed	-String (api key) you want hashed
@return string			-Returns sha 256 hashed string
*/
func hashAPIKey(apiKeyUnhashed string) string {
	apiKeyHashed := sha256.Sum256([]byte(apiKeyUnhashed))
	apiKeyHashedString := hex.EncodeToString(apiKeyHashed[:])
	return apiKeyHashedString
}

// GetAllRegistrations retrieves all registrations associated with an API key.
//
// Returns:
// - []model.Registration: list of registrations
// - error: if iteration or decoding fails
func (f *FireStore) GetAllRegistrations(ctx context.Context, apiKey string) ([]model.Registration, error) {
	iter := f.client.Collection("registrations").
		Doc(apiKey).
		Collection("user_registrations").
		Documents(ctx)

	defer iter.Stop()

	var registrations []model.Registration

	// Iterate through all documents in collection
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		var reg model.Registration
		doc.DataTo(&reg)

		// Include the document ID
		reg.ID = doc.Ref.ID

		registrations = append(registrations, reg)
	}

	return registrations, nil
}

// UpdateRegistration replaces an existing registration entirely.
//
// NOTE: This performs a full overwrite of the document.
func (f *FireStore) UpdateRegistration(ctx context.Context, apiKey string, id string, reg model.Registration) error {
	docRef := f.client.Collection("registrations").
		Doc(apiKey).
		Collection("user_registrations").
		Doc(id)
	_, err := docRef.Set(ctx, reg)
	if err != nil {
		return err
	}
	// Overwrites the entire document
	_, err = docRef.Set(ctx, reg)
	return err
}

// DeleteRegistration removes a registration by ID.
//
// Returns error if the document does not exist or deletion fails.
func (f *FireStore) DeleteRegistration(ctx context.Context, apiKey string, id string) error {

	docRef := f.client.Collection("registrations").Doc(apiKey).Collection("user_registrations").Doc(id)

	// check if exists
	_, err := docRef.Get(ctx)
	if err != nil {
		return err
	}

	_, err = docRef.Delete(ctx)
	return err
}

// TweakRegistration performs a partial update (PATCH) on a registration.
//
// Only fields provided in the patch object are updated.
// Uses Firestore's Update() to modify specific fields instead of overwriting.
//
// Supports nested updates for the "features" object using reflection.
func (f *FireStore) TweakRegistration(
	ctx context.Context,
	apiKey string,
	id string,
	patch model.RegistrationPatch,
) error {

	docRef := f.client.Collection("registrations").
		Doc(apiKey).
		Collection("user_registrations").
		Doc(id)

	var updates []firestore.Update

	// Handle top-level fields
	if patch.Country != nil {
		updates = append(updates, firestore.Update{
			Path:  "country",
			Value: *patch.Country,
		})
	}

	if patch.IsoCode != nil {
		updates = append(updates, firestore.Update{
			Path:  "isoCode",
			Value: *patch.IsoCode,
		})
	}

	// Handle nested "features" fields dynamically using reflection

	if patch.Features != nil {
		v := reflect.ValueOf(*patch.Features)
		t := reflect.TypeOf(*patch.Features)

		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i)
			fieldType := t.Field(i)

			// Skip nil pointers safely
			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				continue
			}

			// Extract JSON tag to match Firestore field naming
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "" {
				continue
			}
			jsonTag = strings.Split(jsonTag, ",")[0]

			// Extract actual value (handle pointer vs non-pointer)
			var value interface{}
			if fieldValue.Kind() == reflect.Ptr {
				value = fieldValue.Elem().Interface()
			} else {
				value = fieldValue.Interface()
			}
			// Append nested field update (e.g., "features.temperature")
			updates = append(updates, firestore.Update{
				Path:  "features." + jsonTag,
				Value: value,
			})
		}
	}

	// Update timestamp
	updates = append(updates, firestore.Update{
		Path:  "lastChange",
		Value: time.Now().Format("20060102 15:04"),
	})

	// Execute partial update
	_, err := docRef.Update(ctx, updates)
	return err
}

// APIKeyExists checks whether a given API key (hashed) exists in Firestore.
//
// Returns:
// - true if the key exists
// - false if not found or an error occurs
//
// NOTE: Errors are treated as "not found" for simplicity.
func (f *FireStore) APIKeyExists(ctx context.Context, keyHash string) bool {
	_, err := f.client.Collection("all_api_keys").Doc(keyHash).Get(ctx)
	if err != nil {
		return false
	}
	return true

}
