package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Tile struct {
	Id      string   `json:"_id"`
	Label   string   `json:"label"`
	Tags    []string `json:"tags"`
	Details string   `json:"details"`
}

func HandleRoot(writer http.ResponseWriter, request *http.Request) {
	dir := "./client/dist"
	// Create a file server handler for the chosen directory
	fs := http.FileServer(http.Dir(dir))

	// Strip the prefix from the URL path and serve the files
	http.StripPrefix("/", fs).ServeHTTP(writer, request)

}

func HandlePortfolioImg(w http.ResponseWriter, r *http.Request) {
	var isValid = true
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			isValid = false
		}
		isValid = false
	}

	tokenString := cookie.Value
	_, err = ValidateJWT(tokenString) //this is just for portfolio,
	// otherwise the claims with username would be a good way to get images to the correct user folder
	if err != nil {
		isValid = false
	}

	switch r.Method {
	case http.MethodGet:
		//get the image from the argument
		break
	case http.MethodPost:
		if !isValid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//upload the image and save it in the file system. then return the url
		break
	case http.MethodDelete:
		if !isValid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//delete the image from the filesystem
		break
	}
}
func HandleCreateTile(w http.ResponseWriter, r *http.Request) {
	//cookie, err := r.Cookie("token")
	//if err != nil {
	//	if err == http.ErrNoCookie {
	//		w.WriteHeader(http.StatusUnauthorized)
	//		return
	//	}
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//tokenString := cookie.Value
	//claims, err := ValidateJWT(tokenString)
	//if err != nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	// If we reach here, the token is valid
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
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

	//DBClient.AddPortfolioEntry(CTX,string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	//if it does, update it
	//else create a new tile
	//add the tile to the database
	//return json of the new tile
}

func HandleLoadTiles(w http.ResponseWriter, r *http.Request) {
	//tiles := []Tile{
	//	{
	//		Id:      "lkajf9o924laksgladflkjvafarg334",
	//		Label:   "tile loaded from api",
	//		Tags:    []string{},
	//		Details: "# this is the headline /n /n and these are some details",
	//	},
	//}
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
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If we reach here, the token is valid
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	//check to see if the tile exists
	//if it does, update it
	//else create a new tile
	//add the tile to the database
	//return json of the new tile
}
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If we reach here, the token is valid
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	//delete all images linked in the tile
	//remove the tile from the database
	//return a code to indicate success
}
