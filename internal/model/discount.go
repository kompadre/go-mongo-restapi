package model

type Discount struct {
	Sku                string   `json:"sku,omitempty"`
	Categories         []string `json:"categories,omitempty"`
	DiscountPercentage int      `json:"discount"`
}

func CurrentDiscounts() []Discount {
	return []Discount{
	{
		Sku:                "",
		Categories:         []string{"boots"},
		DiscountPercentage: 30,
	},
	{
		Sku:                "000003",
		Categories:         nil,
		DiscountPercentage: 15,
	}}
}

