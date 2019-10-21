// See page 110.
//!+

package caserola

import "time"

type Config struct {
	Email      string `json:"email"`
	Pwd        string `json:"pwd"`
	Restaurant string `json:"restaurant"`
	UtcH       int    `json:"utcH"`
	UtcM       int    `json:"utcT"`
}

type OrderFeed struct {
	Items []*Order `json:"data"`
}

type Order struct {
	ID         int64     `json:"id"`
	State      string    `json:"state"`
	DatePlaced time.Time `json:"datePlaced"`
}

type Product struct {
	ID           int64 `json:"id"`
	RestaurantID int64 `json:"restaurant-id"`
	SectionID    int64 `json:"section-id"`
}

type Menu struct {
	Appeteazers []*Product
	Mains       []*Product
	Deserts     []*Product
}

type RestaurantConfig struct {
	URL                  string
	AppeteazersSectionID []int64
	MainsSectionID       []int64
	DesertsSectionID     []int64
}

//!-
