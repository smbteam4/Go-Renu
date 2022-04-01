package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"smbteam4/common"
	"smbteam4/models"
	"smbteam4/mongodb"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBCONN string = "mongodb://localhost:27017"

var (
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	fmt.Println("smbteam4 - Rest API v1.0")
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/v1/add-survivors", createNewsurvivors).Methods("POST")
	myRouter.HandleFunc("/api/v1/update-location", updateSurvivorsLocation).Methods("POST")
	myRouter.HandleFunc("/api/v1/infected", updateInfection).Methods("POST")

	myRouter.HandleFunc("/api/v1/percentage/{spec}", percentageSpecification)
	myRouter.HandleFunc("/api/v1/survivors/{spec}", listSurvivors)
	myRouter.HandleFunc("/api/v1/robots/all", listAllRobots)

	port := common.GetenvData("PORT")
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}
func createNewsurvivors(w http.ResponseWriter, r *http.Request) {

	var err error
	client, ctx, cancel, err = mongodb.MongoDBconnect(DBCONN)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(client, ctx, cancel)

	var survivor_data models.Survivor
	reqBody, _ := ioutil.ReadAll(r.Body)
	err1 := json.Unmarshal(reqBody, &survivor_data)

	//check if survivor exist in DB
	filter := bson.M{
		"Id": bson.M{"$eq": survivor_data.Id},
	}
	//Fetch total document count
	count, err := mongodb.MongoDBCountDocuments(client, ctx, common.GetenvData("DB_NAME"), common.GetenvData("COLLECTION_NAME"), filter)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}
	if count > 0 {
		response := models.Response{
			Message: "Survivor already exist",
			Success: false,
			Data:    nil,
			Error:   err,
		}
		fmt.Fprintf(w, "%+v", response)
	} else {

		var document interface{}

		document = bson.D{
			{"Id", survivor_data.Id},
			{"Name", survivor_data.Name},
			{"Age", survivor_data.Age},
			{"Gender", survivor_data.Gender},
			{"Location", survivor_data.Location},
			{"Resource", survivor_data.Resource},
			{"Status", 0},
		}
		insertResult, err := mongodb.MongoDBinsertOne(client, ctx, common.GetenvData("DB_NAME"),
			common.GetenvData("COLLECTION_NAME"), document)
		if err != nil {
			fmt.Fprintf(w, "%+v", err1)
		}
		response := models.Response{
			Message: "Survivor data inserted successfully",
			Success: true,
			Data:    insertResult.InsertedID,
			Error:   err,
		}
		fmt.Fprintf(w, "%+v", response)
	}

}

// Function to modify survivors Location
func updateSurvivorsLocation(w http.ResponseWriter, r *http.Request) {

	client, ctx, cancel, err := mongodb.MongoDBconnect(DBCONN)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(client, ctx, cancel)

	var survivor models.Survivor
	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &survivor)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	filter := bson.M{
		"Id": bson.M{"$eq": survivor.Id},
	}
	update := bson.M{
		"$set": bson.M{"Location": survivor.Location},
	}
	result, err := mongodb.MongoDBUpdateOne(client, ctx, common.GetenvData("DB_NAME"), common.GetenvData("COLLECTION_NAME"), filter, update)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}
	response := models.Response{
		Message: "Survivor Location Modified successfully",
		Success: true,
		Data:    result.ModifiedCount,
		Error:   err,
	}
	fmt.Fprintf(w, "%+v", response)
}

//Function to Update data of infected survivor
func updateInfection(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := mongodb.MongoDBconnect(DBCONN)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(client, ctx, cancel)

	var survivor models.Survivor
	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &survivor)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	//Update status
	filter := bson.M{
		"Id": bson.M{"$eq": survivor.Id},
	}
	update := bson.M{
		"$inc": bson.M{"Status": 1},
	}
	result, err := mongodb.MongoDBUpdateOne(client, ctx, common.GetenvData("DB_NAME"), common.GetenvData("COLLECTION_NAME"), filter, update)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	response := models.Response{
		Message: "Survivor data Modified successfully",
		Success: true,
		Data:    result.ModifiedCount,
		Error:   err,
	}
	fmt.Fprintf(w, "%+v", response)
}

// Function to find Percentage of infected/non-infected survivors
func percentageSpecification(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := mongodb.MongoDBconnect(DBCONN)
	if err != nil {
		panic(err)
	}
	defer mongodb.MongoDBclose(client, ctx, cancel)

	var total_count, count int64
	filter := bson.M{
		"Id": bson.M{"$ne": ""},
	}
	//Fetch total document count
	total_count, err = mongodb.MongoDBCountDocuments(client, ctx, common.GetenvData("DB_NAME"), common.GetenvData("COLLECTION_NAME"), filter)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}

	vars := mux.Vars(r)
	key := vars["spec"]

	//check the request specification
	switch key {
	case "infected":
		filter = bson.M{
			"Status": bson.M{"$gt": 2},
		}
	case "non-infected":
		filter = bson.M{
			"Status": bson.M{"$lt": 3},
		}
	default:
		fmt.Fprintf(w, "%+v", "404 page not found")
	}
	count, err = mongodb.MongoDBCountDocuments(client, ctx, common.GetenvData("DB_NAME"), common.GetenvData("COLLECTION_NAME"), filter)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
	}
	// Calculate Percentage
	percentage := int((float64(count) / float64(total_count)) * 100)
	percentageStr := fmt.Sprintf("%d", percentage)

	response := models.Response{
		Message: "Percentage",
		Success: true,
		Data:    percentageStr + "%",
		Error:   err,
	}
	fmt.Fprintf(w, "%+v", response)

}

//Function to List infected/non-infected survivors
func listSurvivors(w http.ResponseWriter, r *http.Request) {

	client, ctx, cancel, err := mongodb.MongoDBconnect(DBCONN)
	if err != nil {
		panic(err)
	}

	// Free the resource when main function is  returned
	defer mongodb.MongoDBclose(client, ctx, cancel)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}
	vars := mux.Vars(r)
	key := vars["spec"]
	switch key {
	case "infected":
		filter = bson.M{
			"Status": bson.M{"$gt": 2},
		}
	case "non-infected":
		filter = bson.M{
			"Status": bson.M{"$lt": 3},
		}
	default:
		fmt.Fprintf(w, "%+v", "404 page not found")
	}
	option = bson.D{{"_id", 0}}
	cursor, err := mongodb.MongoDBquery(client, ctx, common.GetenvData("DB_NAME"), common.GetenvData("COLLECTION_NAME"), filter, option)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	for _, doc := range results {
		fmt.Println(doc)
		fmt.Fprintf(w, "%+v\n", doc)
	}

}

// Function to Connect to the Robot CPU system
func listAllRobots(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", common.GetenvData("ROBO_CPU_URL"), nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	response := models.Response{
		Message: "Connect to the Robot CPU system",
		Success: true,
		Data:    string(bodyBytes),
		Error:   err,
	}

	fmt.Fprintf(w, "%+v", response)
}
