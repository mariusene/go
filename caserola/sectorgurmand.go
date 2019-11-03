package caserola

import "net/http"

//Sectorgurmand ...
type Sectorgurmand struct{}

//FeedMenu ...
func (*Sectorgurmand) FeedMenu(cookies []*http.Cookie) (*Menu, error) {
	config := RestaurantConfig{
		URL:                  "https://corporate.caserola.ro/restaurant/sectorgurmand",
		AppeteazersSectionID: []int64{504},
		MainsSectionID:       []int64{505},
		DesertsSectionID:     []int64{506},
	}
	return defaultCrawler(config)(cookies)
}

//MakeLunch ...
func (*Sectorgurmand) MakeLunch(menu *Menu, noDesert bool) []*Product {
	return MakeLunchByRandom(menu, noDesert)
}
