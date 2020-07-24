package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type products struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Type        string             `bson:"type,omitempty"`
	OwnerName   string             `bson:"ownerName,omitempty"`
	Size        string             `bson:"size,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Overview    string             `bson:"overview,omitempty"`
	Pictures    string             `bson:"pictures,omitempty"`
	Price       string             `bson:"price,omitempty"`
	SeaPort     string             `bson:"seaPort,omitempty"`
	Location    string             `bson:"location,omitempty"`
	City        string             `bson:"city,omitempty"`
	PhoneNumber string             `bson:"phoneNumber,omitempty"`
	Kitchen     bool               `bson:"kitchen,omitempty"`
	Wifi        bool               `bson:"wifi,omitempty"`
	V           int32              `bson:"__v,omitempty"`
}

func connectDB() *mongo.Collection {

	clientOptions := options.Client().ApplyURI("mongodb+srv://<username>:<password>@cluster0.fdbop.mongodb.net/<Db_Name>?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("Hangiliman").Collection("products")

	return collection
}

type filterDatas struct {
	City    string `json:"city"`
	SeaPort string `json:"seaPort"`
}

func allendpoints(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		fmt.Println("GET called")
		res.Write([]byte(`{"message": "This method not available"}`))
	case "POST":
		fmt.Println("POST called")
		var filter filterDatas
		err := json.NewDecoder(req.Body).Decode(&filter)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		collection := connectDB()

		var cur *mongo.Cursor
		city := fmt.Sprintf("%v", filter.City)
		seaport := fmt.Sprintf("%v", filter.SeaPort)
		if len(city) > 0 && len(seaport) > 0 {
			cur, _ = collection.Find(context.TODO(), bson.M{"city": city, "seaPort": seaport})
		} else if len(city) > 0 {
			cur, _ = collection.Find(context.TODO(), bson.M{"city": city})
		} else {
			res.Write([]byte(`{"code":"404", "message":"There is no filter"}`))
			return
		}

		var prods []products

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var product products

			err := cur.Decode(&product)
			if err != nil {
				log.Fatal(err)
			}
			prods = append(prods, product)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(res).Encode(prods)
	default:
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"message": "status not found"}`))
	}
}
func home(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`{"code":"404", "message":"not found homepage :)"}`))
}
func main() {
	port := os.Getenv("PORT")
	fmt.Println("The Endpoint is running...")
	router := mux.NewRouter()
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/getProducts", allendpoints).Methods("GET")
	router.HandleFunc("/getProducts", allendpoints).Methods("POST")
	fmt.Println("Server listening on ", port)
	http.ListenAndServe(":"+port, router)
}
