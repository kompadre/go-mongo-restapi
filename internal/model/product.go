package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Price struct {
	Original           int    `json:"original"`
	Final              int    `json:"final"`
	DiscountPercentage int    `json:"discount_percentage,omitempty"`
	Currency           string `json:"currency"`
}

type Product struct {
	Sku           string `json:"sku" bson:"sku"`
	Name          string `json:"name" bson:"name"`
	Category      string `json:"category" bson:"category"`
	OriginalPrice int    `json:"price" bson:"original_price"`
	Price         Price  `json:"price_struct,omitempty"`
}

type Products struct {
	Data []Product `json:"products" bson:"products"`
}

func (ps *Products) ApplyDiscounts(ds []Discount) {
	for _, d := range ds {
		for k, _ := range ps.Data {
			ps.Data[k].ApplyDiscount(d)
		}
	}
}

type Collection []Product

func New(sku string, name string, category string, price Price) *Product {
	return &Product{Sku: sku, Name: name, Price: price, Category: category}
}

// Making sure the response matches what spec says
func (p Product) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	value, _ := json.Marshal(p.Sku)
	buffer.WriteString(`"sku":` + string(value) + `,`)
	value, _ = json.Marshal(p.Name)
	buffer.WriteString(`"name":` + string(value) + `,`)
	value, _ = json.Marshal(p.Category)
	buffer.WriteString(`"category":` + string(value) + `,`)
	value, _ = json.Marshal(p.Price)
	buffer.WriteString(`"price":` + string(value))
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}


func (p *Product) DiscountApplies(d Discount) bool {
	if d.Sku == p.Sku {
		return true
	}
	if sliceContains(p.Category, d.Categories) {
		return true
	}
	return false
}

func (p *Product) ApplyDiscount(d Discount) bool {
	if d.DiscountPercentage <= p.Price.DiscountPercentage || !p.DiscountApplies(d) {
		return false
	}
	fmt.Println("Applying discount!")
	newDiscountPercentage := d.DiscountPercentage
	newPrice := int(float32(p.OriginalPrice) * 0.01 * float32(100-newDiscountPercentage))
	fmt.Println(newPrice)
	price := p.Price
	price.DiscountPercentage = newDiscountPercentage
	price.Final = newPrice
	p.Price = price
	return true
}

func sliceContains(needle string, hay []string) bool {
	for _, element := range hay {
		if element == needle {
			return true
		}
	}
	return false
}
