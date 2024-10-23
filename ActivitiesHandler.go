package main

import (
	"Tubes/db"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

func HandleInboxJson(input string, actorId string, client *db.MongoClient) error {
	isValid, mapped := checkAndUnmarshalToMap(input)
	if !isValid {
		return errors.New("json is not an activity stream")
	}
	// decide what to do based on type
	// deduplicate activity
	// handle any server side effects
	// add activity to user's ordered collection
	switch mapped["type"] {
	case "Accept":
		fmt.Println("Handling Accept")
	case "Add":
		fmt.Println("Handling Add")
	case "Announce":
		fmt.Println("Handling Announce")
	case "Block":
		fmt.Println("Handling Block")
	case "Create":
		fmt.Println("Handling Create")
	case "Delete":
		fmt.Println("Handling Delete")
	case "Dislike":
		fmt.Println("Handling Dislike")
	case "Follow":
		fmt.Println("Handling Follow")
	case "Ignore":
		fmt.Println("Handling Ignore")
	case "Invite":
		fmt.Println("Handling Invite")
	case "Join":
		fmt.Println("Handling Join")
	case "Leave":
		fmt.Println("Handling Leave")
	case "Like":
		fmt.Println("Handling Like")
	case "Move":
		fmt.Println("Handling Move")
	case "Reject":
		fmt.Println("Handling Reject")
	case "Remove":
		fmt.Println("Handling Remove")
	case "TentativeAccept":
		fmt.Println("Handling TentativeAccept")
	case "Undo":
		fmt.Println("Handling Undo")
	case "Update":
		fmt.Println("Handling Update")
	case "TentativeReject":
	case "Question":
	case "Offer":
	case "Listen":
	case "Flag":
	case "Arrive":
	case "Read":
	case "Travel":
	case "View":
	default:
		fmt.Println("Unhandled action")
	}

	return nil
}

func HandleOutboxJson(input string) error {
	isValid, mapped := checkAndUnmarshalToMap(input)
	if !isValid {
		return errors.New("json is not an activity stream")
	}
	// decide what to do based on type
	// deduplicate activity
	// handle any server side effects
	// send to other server if necessary
	// add activity to user's outbox ordered collection
	switch mapped["type"] {
	case "Accept":
		fmt.Println("Handling Accept")
	case "Add":
		fmt.Println("Handling Add")
	case "Announce":
		fmt.Println("Handling Announce")
	case "Block":
		fmt.Println("Handling Block")
	case "Create":
		fmt.Println("Handling Create")
	case "Delete":
		fmt.Println("Handling Delete")
	case "Dislike":
		fmt.Println("Handling Dislike")
	case "Follow":
		fmt.Println("Handling Follow")
	case "Ignore":
		fmt.Println("Handling Ignore")
	case "Invite":
		fmt.Println("Handling Invite")
	case "Join":
		fmt.Println("Handling Join")
	case "Leave":
		fmt.Println("Handling Leave")
	case "Like":
		fmt.Println("Handling Like")
	case "Move":
		fmt.Println("Handling Move")
	case "Reject":
		fmt.Println("Handling Reject")
	case "Remove":
		fmt.Println("Handling Remove")
	case "TentativeAccept":
		fmt.Println("Handling TentativeAccept")
	case "Undo":
		fmt.Println("Handling Undo")
	case "Update":
		fmt.Println("Handling Update")
	case "TentativeReject":
	case "Question":
	case "Offer":
	case "Listen":
	case "Flag":
	case "Arrive":
	case "Read":
	case "Travel":
	case "View":
	default:
		fmt.Println("Unhandled action")
	}

	return nil
}

func checkAndUnmarshalToMap(input string) (bool, map[string]string) {
	var handleObj map[string]string
	json.Unmarshal([]byte(input), &handleObj)
	//check to see if this is the right kind of json
	if handleObj["@context"] != "https://www.w3.org/ns/activitystreams" {
		//maybe context was an array so try that
		var context []string
		json.Unmarshal([]byte(handleObj["@context"]), &context)
		if !slices.Contains(context, "https://www.w3.org/ns/activitystreams") {
			return false, nil
		}
	}
	return true, handleObj
}

func getObject(objectJson string) map[string]string {
	var obj map[string]string
	json.Unmarshal([]byte(objectJson), &obj)
	return obj
}
