package handlers

import (
	"encoding/json"
	"log"
	"minna-style-hub/database"
	models "minna-style-hub/model"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllItems handles GET request to fetch all items// GetAllItems handles GET request to retrieve all items
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := database.GetAllItems()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/item/")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	item, err := database.GetItem(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// AddItem handles POST request to add an item
func AddItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Generate a new ObjectId for the item
	newItemID := primitive.NewObjectID()
	newItem.ID = newItemID.Hex() // Convert ObjectID to string

	err := database.AddItem(newItem)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateItem handles PUT request to update an item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := database.UpdateItem(updatedItem)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteItem handles DELETE request to delete an item by its MongoDB _id
// DeleteItem handles DELETE request to delete an item by its MongoDB _id
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	err := database.DeleteItem(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
