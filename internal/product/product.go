package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Price struct {
	Original int
	Final int
	DiscountPercentage int
	Currency string
}

type Discount struct {
	Sku string
	Categories []string
	DiscountPercentage int
}

type Product struct {
	Sku string `json:"sku" bson:"sku"`
	Name string `json:"name" bson:"name"`
	Category string `json:"category" bson:"category"`
	OriginalPrice int `json:"price" bson:"original_price"`
	Price Price `json:"price_struct"`
}

type Products struct {
	Data []Product `json:"products"`
}


type Collection []Product

func New(sku string, name string, category string, price Price) *Product {
	return &Product{Sku: sku, Name: name, Price: price, Category: category}
}

func (p *Product) DiscountApplies(d Discount) bool {
	if (d.Sku == p.Sku) {
		return true
	}
	if sliceContains(p.Category, d.Categories) {
		return true
	}
	return false
}

func (p *Product) ApplyDiscount(d Discount) bool {
	if d.DiscountPercentage < 0 || d.DiscountPercentage > 100 || p.Price.Original <= 0 {
		return false
	}
	if p.DiscountApplies(d) && d.DiscountPercentage > p.Price.DiscountPercentage {
		newDiscountPercentage := d.DiscountPercentage
		p.Price.DiscountPercentage = newDiscountPercentage
		p.Price.Final = int((0.01 * float64(100-newDiscountPercentage)) * float64(p.Price.Original))
		return true
	}
	return false
}

func Retrieve() *Products {
	var products Products

	cwd, _ := os.Getwd()
	jsonFile, _ := os.Open(cwd + "/assets/products.json")
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &products
}

func sliceContains(needle string, hay []string) bool {
	for _, element := range hay {
		if element == needle {
			return true
		}
	}
	return false
}

/*
func (p Price) MarshalJSON() ([]byte, error) {
	marshalledData, err := json.Marshal(p)
	return marshalledData, err
}
*/
