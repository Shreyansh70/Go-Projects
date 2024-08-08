package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title  string `bson:"title"`
    Year   string    `bson:"year"`
    Genre  string `bson:"genre"`
}

var client *mongo.Client
var clientOnce sync.Once

func GetMongoClient() (*mongo.Client, error) {
    var err error
    clientOnce.Do(func() {
        clientOptions := options.Client().ApplyURI("mongodb+srv://Shreyansh:Mdash%409987@mernatoz.c9nsqws.mongodb.net/?retryWrites=true&w=majority&appName=MERNAtoZ")
        client, err = mongo.Connect(context.TODO(), clientOptions)
        if err != nil {
            log.Fatal(err)
        }

        // Check the connection
        err = client.Ping(context.TODO(), nil)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Connected to MongoDB!")
    })

    return client, err
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books" , getAllBooks).Methods("GET")
	r.HandleFunc("/books/{title}" , getSpecificBook).Methods("GET")
	r.HandleFunc("/books" , insertBook).Methods("POST")
	r.HandleFunc("/books/{title}" , deleteBook).Methods("DELETE")

	http.ListenAndServe(":8000" , r)
}

func deleteBook(w http.ResponseWriter, r *http.Request){
	fmt.Println("Delete Functionality Called...")
	client, err := GetMongoClient()
	if err != nil {
        log.Fatal(err)
    }
	collection := client.Database("TestDB").Collection("Books")

	params := mux.Vars(r)
	filter := bson.D{{"title", params["title"]}}

	_,err = collection.DeleteOne(context.TODO() , filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w , "Record Deleted Successfully")
}

func insertBook(w http.ResponseWriter, r *http.Request){
	client, err := GetMongoClient()
	if err != nil {
        log.Fatal(err)
    }
	collection := client.Database("TestDB").Collection("Books")

	var book Book

	_ = json.NewDecoder(r.Body).Decode(&book)

	_,err = collection.InsertOne(context.TODO() , book)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w , "Document inserted successfully")
}

func getSpecificBook(w http.ResponseWriter, r *http.Request){
	client, err := GetMongoClient()
	if err != nil {
        log.Fatal(err)
    }

	var book Book
	collection := client.Database("TestDB").Collection("Books")
	params := mux.Vars(r)
	filter := bson.D{{"title", params["title"]}}
	err = collection.FindOne(context.TODO() , filter).Decode(&book)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-type" , "application/json")
	json.NewEncoder(w).Encode(book)
}

func getAllBooks(w http.ResponseWriter, r *http.Request){
	client, err := GetMongoClient()
	if err != nil {
        log.Fatal(err)
    }
	collection := client.Database("TestDB").Collection("Books")

	cursor, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(context.TODO())

    var movies []Book
    for cursor.Next(context.TODO()) {
        var movie Book
        err := cursor.Decode(&movie)
        if err != nil {
            log.Fatal(err)
        }
        movies = append(movies, movie)
    }

    if err := cursor.Err(); err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}