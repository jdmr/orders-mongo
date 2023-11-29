package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Orders API")

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers", createCustomer).Methods("POST")
	r.HandleFunc("/customers/{customerID}", getCustomer).Methods("GET")
	r.HandleFunc("/customers/{customerID}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{customerID}", deleteCustomer).Methods("DELETE")

	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{productID}", getProduct).Methods("GET")
	r.HandleFunc("/products/{productID}", updateProduct).Methods("PUT")
	r.HandleFunc("/products/{productID}", deleteProduct).Methods("DELETE")

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Listening on port 8080")
	log.Fatal(srv.ListenAndServe())
}
