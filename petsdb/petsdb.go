package petsdb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/google/uuid"
	"cloud.google.com/go/datastore"
)

var projectID string

// Pet model stored in Datastore
type Pet struct {
	Added   time.Time `datastore:"added"`
	Caption string    `datastore:"caption"`
	Email   string    `datastore:"email"`
	Image   string    `datastore:"image"`
	Likes   int       `datastore:"likes"`
	Owner   string    `datastore:"owner"`
	Petname string    `datastore:"petname"`
	Name    string    // The ID used in the datastore.
}

// GetPets Returns all pets from datastore ordered by likes in Desc Order
func GetPets() ([]Pet, error) {

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	var pets []Pet
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	// Create a query to fetch all Pet entities".
	fmt.Println("Starting reading")
	query := datastore.NewQuery("Pet").Order("-likes")
	keys, err := client.GetAll(ctx, query, &pets)
	fmt.Println(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		pets[i].Name = key.Name
	}

	client.Close()
	return pets, nil
}

// AddPets adds two pets
func AddPets() error {

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	//var pets []Pet
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	e := Pet{Email: "dragon@gmail.com", Petname: "Dragon", Image: "https://i.pinimg.com/originals/43/78/8c/43788cc3cf3ba73b03c4380f2ba6e4e1.jpg", Likes: 7}
	k := datastore.NameKey("Pet", "Dragon", nil)
	if _, err := client.Put(ctx, k, &e); err != nil {
		// Handle error.
	}

	e = Pet{Email: "unicorn@gmail.com", Petname: "Unicorn", Image: "https://i.pinimg.com/736x/e6/ca/93/e6ca93645d17b70eb612cd6312e0de36.jpg", Likes: 5}
	k = datastore.NameKey("Pet", "Unicorn", nil)
	if _, err := client.Put(ctx, k, &e); err != nil {
		// Handle error.
	}
	client.Close()
	return err
}

func PutPet(pet Pet){
	id := uuid.New()
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	k := datastore.NameKey("Pet", "Pet" + id.String(), nil)
	_, err = client.Put(ctx, k, &pet)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("new Pet:", pet)
	client.Close()
}

func DeletePet(petId string){
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	key := datastore.NameKey("Pet", petId, nil)
	fmt.Println("Key to delete: ")
	fmt.Println(key)
	if err := client.Delete(ctx, key); err != nil {
		log.Fatalf("Could not delete item: %v", err)
	}
	log.Println("deleted pet with id:", petId)
	client.Close()
}
