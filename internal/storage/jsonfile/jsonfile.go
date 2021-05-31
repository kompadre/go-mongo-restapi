package jsonfile

import (
	"encoding/json"
	"go-mongo-restapi/internal/model"
	"io/ioutil"
	"log"
	"os"
)

func RetrieveProducts() *model.Products {
	var products model.Products
	cwd, _ := os.Getwd()
	jsonFile, err := os.Open(cwd + "/assets/products.json")
	if err != nil {
		log.Panicf("Couldnt open %s for reading: %v", cwd + "/assets/products.json", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Panicf("Error retrieving products from json file %s: %v, %v", cwd + "/assets/products.json", err, jsonFile)
	}
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		log.Panicf("Error unmarshalling products: %v", err)
	}
	log.Printf("%v", products)
	return &products
}
