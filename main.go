package main

import (
	"bytes"
	"drehnstrom.com/go-pets/petsdb"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type deletePet struct {
	PetId string
}

var projectID string

func main() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}
	log.Printf("GOOGLE_CLOUD_PROJECT is set to %s", projectID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	log.Printf("Port set to: %s", port)

	fs := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()

	// This serves the static files in the assets folder
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// The rest of the routes
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/pets", petsHandler)
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/api/add", apiAddPetHandler)
	mux.HandleFunc("/api/delete", apiDeletePetHandler)

	log.Printf("Webserver listening on Port: %s", port)
	http.ListenAndServe(":"+port, mux)
}

func apiAddPetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST"{
		var pet petsdb.Pet
		err := json.NewDecoder(r.Body).Decode(&pet)
		if err != nil {
			fmt.Print(err)
		}
		log.Println("[INFO] Add Pet ", pet)
		petsdb.PutPet(pet)
	}else{
		http.Error(w, "Bad request to api/add endpoint, it supports only POST", http.StatusInternalServerError)
		log.Println("Bad request to api/add endpoint, it supports only POST")
		return
	}
}

//curl -d '{"petId":"Pet750f0c2e-1c34-41d5-ace7-894108931230"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/delete
//curl -X POST http://localhost:8080/api/delete
func apiDeletePetHandler(w http.ResponseWriter, r *http.Request) {
	var pet deletePet
	err := json.NewDecoder(r.Body).Decode(&pet)
	log.Println(pet)
	if err != nil {
		fmt.Print(err)
	}
	if r.Method == "POST"{
		log.Println("[INFO] Delete Pet with id: ", pet.PetId)
		petsdb.DeletePet(pet.PetId)
	}else{
		http.Error(w, "Bad request to api/delete endpoint, it supports only POST", http.StatusInternalServerError)
		log.Println("Bad request to api/delete endpoint, it supports only POST")
		return
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	data := AboutPageData{
		PageTitle: "Add Pet",
	}

	var tpl = template.Must(template.ParseFiles("templates/add.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("Add Page Served")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var pets []petsdb.Pet
	pets, error := petsdb.GetPets()
	if error != nil {
		fmt.Print(error)
	}

	data := HomePageData{
		PageTitle: "Pets Home Page",
		Pets:      pets,
	}

	var tpl = template.Must(template.ParseFiles("templates/index.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("Home Page Served")
}

func petsHandler(w http.ResponseWriter, r *http.Request) {

	if error := petsdb.AddPets(); error != nil {
		fmt.Print(error)
	}

	data := PetsPageData{
		PageTitle: "Adding Pets Page",
	}

	var tpl = template.Must(template.ParseFiles("templates/pets.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("Pets Page Served")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := AboutPageData{
		PageTitle: "About Go Pets",
	}

	var tpl = template.Must(template.ParseFiles("templates/about.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("About Page Served")
}

// HomePageData for Index template
type HomePageData struct {
	PageTitle string
	Pets      []petsdb.Pet
}

// AboutPageData for About template
type AboutPageData struct {
	PageTitle string
}

// PetsPageData for About template
type PetsPageData struct {
	PageTitle string
}
