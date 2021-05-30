package model

type Discount struct {
	Sku                string   `json:"sku,omitempty"`
	Categories         []string `json:"categories,omitempty"`
	DiscountPercentage int      `json:"discount"`
}
