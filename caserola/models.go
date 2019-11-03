// See page 110.
//!+

package caserola

import (
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Email      string `json:"email"`
	Pwd        string `json:"pwd"`
	Restaurant string `json:"restaurant"`
	UtcHHmm    string `json:"uctHH:mm"`
	NoDesert   bool   `json:"noDesert"`
}

func (cf *Config) GetUtcH() int {
	return cf.getUtcHHmm(0, 7)
}
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
