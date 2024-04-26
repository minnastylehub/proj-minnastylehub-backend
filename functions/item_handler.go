package functions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"minna-style-hub/database"
	models "minna-style-hub/model"
	"minna-style-hub/sendemail"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllItems handles GET request to fetch all items// GetAllItems handles GET request to retrieve all items
// GetAllItems handles GET request to fetch all items with pagination
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page := 1
	pageSize := 10 // Default page size
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
		page = p
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		pageSize = l
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve items from the database with pagination
	items, err := database.GetItemsWithPagination(offset, pageSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve total count of items
	totalCount, err := database.GetTotalItemCount()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Construct paginated response
	response := struct {
		Meta struct {
			Count  int `json:"count"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		} `json:"meta"`
		Result []models.Item `json:"result"`
	}{
		Meta: struct {
			Count  int `json:"count"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Count:  totalCount,
			Limit:  pageSize,
			Offset: offset,
		},
		Result: items,
	}

	// Set response headers and encode response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

func GetFeedback(w http.ResponseWriter, r *http.Request) {
	emailFrom := os.Getenv("EMAIL_FROM")
	if emailFrom == "" {
		log.Fatal("EMAIL_FROM not found in .env file")
	}
	emailPass := os.Getenv("EMAIL_PASSWORD")
	if emailPass == "" {
		log.Fatal("EMAIL_PASSWORD not found in .env file")
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var feedbackItem models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedbackItem); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	name := feedbackItem.Name
	email := feedbackItem.Email
	message := feedbackItem.Message

	// Send email
	to := "tjminna@gmail.com" // Replace with your Gmail address
	subject := "Feedback from Contact Form"
	body := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s", name, email, message)
	err = sendemail.SendEmail(to, subject, body, emailFrom, emailPass)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	fmt.Println(email)
	// Send thankyou email
	thank_to := email // Replace with your Gmail address
	thank_subject := "Thank You for Your Feedback!"
	thank_body := fmt.Sprintf("Thank you for your feedback, %s! We appreciate your time.<br><br><img src=\"cid:thankyou-image\">", name)

	err = sendemail.ThankyouEmail(thank_to, thank_subject, thank_body, emailFrom, emailPass)
	if err != nil {
		http.Error(w, "Failed to send thank you email", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	fmt.Fprintf(w, "Thank you email sent successfully")
}

func getImageBase64() string {
	// Path to the image file
	imagePath := "C:\\Users\\binwi\\Desktop\\Folders\\study\\proj-minnastylehub-backend\\assets\\thankyou-image.png"

	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return ""
	}

	// Encode image data to base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return imageBase64
}
