package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Tile struct {
	Id      string   `json:"id"`
	Label   string   `json:"label"`
	Tags    []string `json:"tags"`
	Details string   `json:"details"`
}

var tokens []string

func HandleRoot(writer http.ResponseWriter, request *http.Request) {
	dir := "./client/dist"
	// Create a file server handler for the chosen directory
	fs := http.FileServer(http.Dir(dir))

	// Strip the prefix from the URL path and serve the files
	http.StripPrefix("/", fs).ServeHTTP(writer, request)

}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("code")
	if token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}

	// Set the token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true, // Set to true if using HTTPS
		Path:     "/",
	})

	// Redirect to the home page or another page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetImg(w http.ResponseWriter, r *http.Request) {
	// Get the image from the argument
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the appropriate content type based on the file extension
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file information", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)
}

func PostImg(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a unique filename for the image
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filePath := filepath.Join("uploads", fileName)

	// Create the directory if it does not exist
	err = os.MkdirAll("uploads", 0755)
	if err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}
	// Create a new file in the uploads directory
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		return
	}

	// Generate the URL for the saved image
	imageURL := "http://localhost:8080/portfolio/img?filename=" + fileName

	// Return the URL in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"url": imageURL})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	token, _ := randomHex(20)
	tokens = append(tokens, token)
	//// Set the cookie domain to localhost
	//domain := "localhost" // Update this with your specific domain
	//
	//http.SetCookie(w, &http.Cookie{
	//	Name:     "Authorization",
	//	Value:    token,
	//	HttpOnly: true,
	//	Secure:   false, // Set to true if using HTTPS
	//	Path:     "/",
	//	Domain:   domain, // Set the cookie domain here
	//	SameSite: http.SameSiteNoneMode,
	//})
	w.Header().Set("token", token)
	w.WriteHeader(http.StatusOK)
}

func HandleCreateTile(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var nTile Tile
	err := d.Decode(&nTile)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if the tile with the same ID already exists
	existingTile, err := DBClient.GetPortfolioEntryByID(CTX, nTile.Id)
	if err != nil {
		// Handle error when checking for existing tile
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingTile != nil {
		// Update the existing tile instead of adding a new entry
		_, err := DBClient.UpdatePortfolioEntry(CTX, nTile.Id, nTile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Add a new entry if the tile doesn't exist
		id, err := DBClient.AddPortfolioEntry(CTX, nTile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nTile.Id = id.(primitive.ObjectID).Hex()
	}

	jsonData, err := json.Marshal(nTile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func HandleLoadTiles(w http.ResponseWriter, r *http.Request) {
	tiles, err := DBClient.GetPortfolioEntries(CTX)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Marshal the tiles array into JSON
	jsonData, err := json.Marshal(tiles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func HandleEditTile(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var nTile Tile
	err := d.Decode(&nTile)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := DBClient.AddPortfolioEntry(CTX, nTile)
	nTile.Id = id.(primitive.ObjectID).Hex()

	jsonData, err := json.Marshal(nTile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func HandleDelete(w http.ResponseWriter, r *http.Request) {

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var nTile Tile
	err := d.Decode(&nTile)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract the ID from the incoming Tile data
	deleteID := nTile.Id

	_, err = DBClient.DeletePortfolioEntry(CTX, deleteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(ExpectedUser))
			expectedPasswordHash := sha256.Sum256([]byte(ExpectedKey))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {

				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func CheckToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bToken := r.Header.Get("Authorization")
		if bToken != "" {
			reqToken := strings.Split(bToken, " ")[1]
			for _, token := range tokens {
				if token == reqToken {
					next.ServeHTTP(w, r)
					return
				}
			}
		}
		// If token is not found or does not match, return unauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
