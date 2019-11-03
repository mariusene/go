// See page 110.
//!+

package caserola

import (
	"strconv"
	"strings"
	"time"
)

//Config ...
type Config struct {
	Email      string `json:"email"`
	Pwd        string `json:"pwd"`
	Restaurant string `json:"restaurant"`
	UtcHHmm    string `json:"uctHH:mm"`
	NoDesert   bool   `json:"noDesert"`
}

//GetUtcH ...
func (cf *Config) GetUtcH() int {
	return cf.getUtcHHmm(0, 7)
}

//GetUtcM ...
func (cf *Config) GetUtcM() int {
	return cf.getUtcHHmm(1, 30)
}

func (cf *Config) getUtcHHmm(idx, defaultV int) int {
	s := strings.Split(cf.UtcHHmm, ":")
	if len(s) == idx {
		return defaultV
	}

	if v, ok := strconv.Atoi(s[idx]); ok == nil {
		return v
	}
	return defaultV
}

//OrderFeed ...
type OrderFeed struct {
	Items []*Order `json:"data"`
}

//Order ...
type Order struct {
	ID         int64     `json:"id"`
	State      string    `json:"state"`
	DatePlaced time.Time `json:"datePlaced"`
}

//Product ...
type Product struct {
	ID           int64 `json:"id"`
	RestaurantID int64 `json:"restaurant-id"`
	SectionID    int64 `json:"section-id"`
	Price        float64
}

//Menu ...
type Menu struct {
	Appeteazers []*Product
	Mains       []*Product
	Deserts     []*Product
}

//RestaurantConfig ...
type RestaurantConfig struct {
	URL                  string
	AppeteazersSectionID []int64
	MainsSectionID       []int64
	DesertsSectionID     []int64
}

//!-
