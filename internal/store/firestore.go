package store

import (
	model "assignment-2/internal/models"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Store struct {
	client *firestore.Client
}

func NewFirestoreStore(client *firestore.Client) *Store {
	return &Store{client: client}
}

func (f *Store) CreateRegistration(ctx context.Context, reg model.Registration) (string, error) {

	doc := f.client.Collection("registrations").NewDoc()

	reg.ID = doc.ID

	_, err := doc.Set(ctx, reg)
	if err != nil {
		return "", err
	}

	return doc.ID, nil
}
func (f *Store) ApiKeyExists(ctx context.Context, apiKey string) bool {

	_, err := f.client.
		Collection("all_api_keys").
		Doc(apiKey).
		Get(ctx)

	if err != nil {
		//api key cant be found e.g it is unique
		fmt.Println("Api key cant be found")
		return false
	}

	fmt.Println("found api key")
	return true
}

func (f *Store) CreateApiStorage(ctx context.Context, reg model.Authentication) error {
	//setts api
	AllApi := f.client.Collection("all_api_keys").Doc(reg.ApiKeyHash)
	_, err := AllApi.Set(ctx, map[string]interface{}{
		"createdAt": reg.CreatedAt,
	})

	emailDoc := f.client.Collection("authentication_info").Doc(reg.Email)
	//creating nested api key structure
	apiDoc := emailDoc.Collection("api_keys").Doc(reg.ApiKeyHash)

	_, err = apiDoc.Set(ctx, reg)
	if err != nil {
		return err
	}

	return nil
}

func (h *Store) CountApiPerUser(ctx context.Context, email string) (int, error) {
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

// TODO write doxygen commenting
func (f *Store) DeleteAPIkey(ctx context.Context, id string) error {
	docRef := f.client.Collection("all_api_keys").Doc(id)

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

	nestedDocRef := userDoc.Collection("api_keys").Doc(id)

	_, err = nestedDocRef.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (f *Store) GetRegistration(ctx context.Context, id string) (*model.Registration, error) {
	doc, err := f.client.Collection("registrations").Doc(id).Get(ctx)
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

func (f *Store) GetAllRegistrations(ctx context.Context) ([]model.Registration, error) {
	iter := f.client.Collection("registrations").Documents(ctx)
	defer iter.Stop()

	var registrations []model.Registration

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		var reg model.Registration
		if err := doc.DataTo(&reg); err != nil {
			return nil, err
		}

		reg.ID = doc.Ref.ID

		registrations = append(registrations, reg)
	}

	return registrations, nil
}

func (f *Store) UpdateRegistration(ctx context.Context, id string, reg model.Registration) error {

	docRef := f.client.Collection("registrations").Doc(id)

	// Check if exists
	_, err := docRef.Get(ctx)
	if err != nil {
		return err
	}

	reg.ID = id

	_, err = docRef.Set(ctx, reg) // replaces entire document
	if err != nil {
		return err
	}

	return nil
}

func (f *Store) DeleteRegistration(ctx context.Context, id string) error {

	docRef := f.client.Collection("registrations").Doc(id)

	// check if exists
	_, err := docRef.Get(ctx)
	if err != nil {
		return err
	}

	_, err = docRef.Delete(ctx)
	return err
}

// function to change a specific part/parts of a registration with use of the patch method
//func (f *Store) TweakRegistration(ctx context.Context) error {}
