package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Customer struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	customers := []*Customer{}
	collection := client.Database("store").Collection("customers")
	cur, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(r.Context())

	for cur.Next(r.Context()) {
		var customer Customer
		err := cur.Decode(&customer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, &customer)
	}

	err = cur.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID := params["customerID"]
	customer := Customer{}
	collection := client.Database("store").Collection("customers")
	err := collection.FindOne(r.Context(), bson.M{"_id": customerID}).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	customer := &Customer{}
	err := json.NewDecoder(r.Body).Decode(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	customer.ID = uuid.New().String()
	collection := client.Database("store").Collection("customers")
	_, err = collection.InsertOne(r.Context(), customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID := params["customerID"]
	customer := &Customer{}
	err := json.NewDecoder(r.Body).Decode(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	collection := client.Database("store").Collection("customers")
	_, err = collection.UpdateOne(r.Context(), bson.M{"_id": customerID}, bson.M{"$set": customer})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID := params["customerID"]
	collection := client.Database("store").Collection("customers")
	_, err := collection.DeleteOne(r.Context(), bson.M{"_id": customerID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
