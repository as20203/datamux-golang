package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func downlink(w http.ResponseWriter, r *http.Request) {
	var loradataentry downlinkdata

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &loradataentry); err != nil {
		fmt.Println("Unable to unmarshall body", err)
		w.WriteHeader(422)
		return
	}

	response := processDownlinkData(loradataentry)
	if response == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
	}
	return
}

func saveLoraData(w http.ResponseWriter, r *http.Request) {
	var loradataentry loradata

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &loradataentry); err != nil {
		fmt.Println("Unable to unmarshall body", err)
		w.WriteHeader(422)
		return
	}

	response := processLoraData(loradataentry)
	if response == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
	}
	return
}

func showDevices(w http.ResponseWriter, r *http.Request) {
	uri := "mongodb://localhost:27017/datamux"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	var devices []bson.M
	datamuxDatabase := client.Database("datamux")
	devicesCollection := datamuxDatabase.Collection("devices")
	filterCursor, err := devicesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = filterCursor.All(ctx, &devices); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(devices)
	}
}

func showDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uri := "mongodb://localhost:27017/datamux"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	var device bson.M
	datamuxDatabase := client.Database("datamux")
	devicesCollection := datamuxDatabase.Collection("devices")

	if err = devicesCollection.FindOne(ctx, bson.M{"deviceeui": params["id"]}).Decode(&device); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	}
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uri := "mongodb://localhost:27017/datamux"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	datamuxDatabase := client.Database("datamux")
	devicesCollection := datamuxDatabase.Collection("devices")
	deviceResult, deletionError := devicesCollection.DeleteOne(ctx, bson.M{
		"deviceeui": params["id"],
	})
	if deletionError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(deletionError.Error()))
	} else if deviceResult.DeletedCount == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`No device found`))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deviceResult)
	}
}

func addDevice(w http.ResponseWriter, r *http.Request) {
	var dev device
	_ = json.NewDecoder(r.Body).Decode(&dev)
	uri := "mongodb://localhost:27017/datamux"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")
	datamuxDatabase := client.Database("datamux")
	devicesCollection := datamuxDatabase.Collection("devices")
	deviceResult, insertionError := devicesCollection.InsertOne(ctx, bson.D{
		{Key: "deviceeui", Value: dev.Deviceeui},
		{Key: "devicetype", Value: dev.Devicetype},
		{Key: "endpointtype", Value: dev.Endpointtype},
		{Key: "endpointdest", Value: dev.Endpointdest},
		{Key: "access_token", Value: dev.Customer},
		{Key: "incl_radio", Value: dev.InclRadio},
		{Key: "raw_data", Value: dev.RawData},
		{Key: "customer", Value: dev.Customer},
	})
	mod := mongo.IndexModel{
		Keys: bson.M{
			"deviceeui": 1, // index in ascending order
		}, Options: options.Index().SetUnique(true),
	}
	_, indexError := devicesCollection.Indexes().CreateOne(ctx, mod)
	if indexError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(indexError.Error()))
	}
	if insertionError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(insertionError.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deviceResult)
	}

}
