package product

import (
	"reflect"
	"testing"
)

func getTestProduct() Product {
	originalPrice := 89000
	sku := "000001"
	name := "BV Lean leather ankle boots"
	category := "boots"
	price := Price{
		Original:           originalPrice,
		Final:              originalPrice,
		DiscountPercentage: 0,
		Currency:           "EUR",
	}
	return Product{
		Sku:        sku,
		Name:       name,
		Category: category,
		price:      price,
	}
}

func TestNew(t *testing.T) {
	want := getTestProduct()
	originalPrice := 89000
	sku := "000001"
	name := "BV Lean leather ankle boots"
	category := "boots"
	price := Price{
		Original:           originalPrice,
		Final:              originalPrice,
		DiscountPercentage: 0,
		Currency:           "EUR",
	}
	got := New(sku, name, category, price)
	if !reflect.DeepEqual(got, &want) {
		t.Errorf("NewProduct failed, got %q want %q", got, want)
	}
}

func TestProduct_DiscountApplies(t *testing.T) {
	p := getTestProduct()
	d := Discount{
		Sku:                "999999",
		Categories:         []string{"ble", "blo"},
		DiscountPercentage: 0,
	}
	if p.DiscountApplies(d) {
		t.Errorf("Discount %q applies where %q it shouldn't", p, d)
	}
	d.Sku = "000001"
	if !p.DiscountApplies(d) {
		t.Errorf("Discount %q doesn't apply where %q it should", p, d)
	}
	d.Sku = "999999"
	d.Categories = []string{"hats", "boots"}
	if !p.DiscountApplies(d) {
		t.Errorf("Discount %q doesn't apply where %q it should", p, d)
	}

}

func TestProduct_ApplyDiscount(t *testing.T) {
	p := getTestProduct()
	d := Discount{
		Sku:                "000001",
		Categories:         []string{"ble", "blo"},
		DiscountPercentage: 0,
	}
	if p.ApplyDiscount(d) {
		t.Errorf("Discount %q was applied where it shouldn't %q", d, p)
	}

	oldDiscount := p.Price.DiscountPercentage
	oldPrice    := p.Price.Final
	d.DiscountPercentage = 10
	if !p.ApplyDiscount(d) || oldDiscount == p.Price.DiscountPercentage || oldPrice == p.Price.Final {
		t.Errorf("Discount %q wasn't applied while it should have %q", d, p)
	}
}

/*
func TestProduct_ToJson(t *testing.T) {
	p := getTestProduct()
	want := `{"Sku":"000001","Name":"BV Lean leather ankle boots","Categories":["boots"],"Price":{"Original":89000,"Final":89000,"DiscountPercentage":0,"Currency":"EUR"}}`
	get, _ := p.ToJson()
	if want != get {
		t.Errorf("JSON generated %v doesn't match wanted %v", get, want)
	}
}
 */