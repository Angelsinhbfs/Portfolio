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

func LoadTube(w http.ResponseWriter, request *http.Request) {
	// get the basic document info from mongo and build the response object from the user data there
	// for now just return something that looks right

	response := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		"id":                "https://dontremember.me/tubes/angelsin",
		"type":              "Person",
		"preferredUsername": "angelsin",
		"inbox":             "https://dontremember.me/tubes/angelsin/inbox",
		"outbox":            "https://dontremember.me/tubes/angelsin/outbox",
		"endpoints": []map[string]string{
			{
				"sharedInbox": "https://bugle.lol/inbox",
			},
		},
		"name":                      "AngelSin",
		"summary":                   "Site owner",
		"url":                       "https://dontremember.me/tubes/angelsin",
		"manuallyApprovesFollowers": false,
		"discoverable":              true,
		"published":                 "2024-9-30T00:00:00Z",
		"publicKey": []map[string]string{
			{
				"id":           "https://dontremember.me/tubes/angelsin#main-key",
				"owner":        "https://dontremember.me/tubes/angelsin",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	//resource := r.URL.Query().Get("resource")
	// Generate JSON response based on the resource. todo: actually check the database to get this info
	// search db to see if acct (r.URL.Query().Get("resource")) exists
	// if it does then return the value if not return error
	response := map[string]interface{}{

		"subject": "acct:angelsin@dontremember.me",
		"links": []map[string]string{
			{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://dontremember.me/tubes/angelsin",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleInbox(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		// serve inbox?
		break
	case http.MethodPost:
		// add to inbox in db?
		break
	case http.MethodDelete:
		// remove from inbox
		break
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleUserInbox(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		// serve inbox?
		break
	case http.MethodPost:
		// add to inbox in db?
		break
	case http.MethodDelete:
		// remove from inbox
		break
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleOutbox(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		// serve outbox? read the public posts
		break
	case http.MethodPost:
		// add to outbox in db? post publicly
		break
	case http.MethodDelete:
		// remove from outbox. delete public post
		break
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleCollections(writer http.ResponseWriter, request *http.Request) {

}
func HandleImg(writer http.ResponseWriter, request *http.Request) {

}

func HandleRoot(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("received root request")
	host := request.Host
	fmt.Println(host)
	var dir string
	if host == "portfolio.dontremember.me" {
		// Serve the directory for "portfolio.dontremember.me"
		dir = "./client/dist/portfolio"
	} else {
		// Serve a different directory for all other hosts
		dir = "./client/dist/main"
	}

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
