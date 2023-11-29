package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Item struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	ProductID string   `json:"productID,omitempty" bson:"productID,omitempty"`
	Quantity  int      `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Price     float64  `json:"price,omitempty" bson:"price,omitempty"`
	Product   *Product `json:"product,omitempty" bson:"product,omitempty"`
}

type Order struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID string    `json:"customerID,omitempty" bson:"customerID,omitempty"`
	Items      []*Item   `json:"items,omitempty" bson:"items,omitempty"`
	Total      float64   `json:"total,omitempty" bson:"total,omitempty"`
	Status     string    `json:"status,omitempty" bson:"status,omitempty"`
	Created    time.Time `json:"created,omitempty" bson:"created,omitempty"`
	Updated    time.Time `json:"updated,omitempty" bson:"updated,omitempty"`
	Customer   *Customer `json:"customer,omitempty" bson:"customer,omitempty"`
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	orders := []*Order{}
	collection := client.Database("store").Collection("orders")
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "customers"},
			{Key: "localField", Value: "customerID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "customer"},
		}},
	}
	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$customer"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}},
	}
	cur, err := collection.Aggregate(r.Context(), mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(r.Context())

	for cur.Next(r.Context()) {
		var order Order
		err := cur.Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, &order)
	}

	err = cur.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID := params["orderID"]
	order := Order{}
	collection := client.Database("store").Collection("orders")
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "customers"},
			{Key: "localField", Value: "customerID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "customer"},
		}},
	}
	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$customer"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}},
	}
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: orderID},
		}},
	}
	err := collection.FindOne(r.Context(), mongo.Pipeline{lookupStage, unwindStage, matchStage}).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	order := Order{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	order.ID = uuid.New().String()
	order.Created = time.Now()
	order.Updated = time.Now()
	for _, item := range order.Items {
		item.ID = uuid.New().String()
	}
	collection := client.Database("store").Collection("orders")
	_, err = collection.InsertOne(r.Context(), order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID := params["orderID"]
	order := Order{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	order.Updated = time.Now()
	collection := client.Database("store").Collection("orders")
	_, err = collection.UpdateOne(r.Context(), bson.M{"_id": orderID}, bson.M{"$set": order})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID := params["orderID"]
	collection := client.Database("store").Collection("orders")
	_, err := collection.DeleteOne(r.Context(), bson.M{"_id": orderID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
