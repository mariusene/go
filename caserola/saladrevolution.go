package caserola

import "net/http"

//Saladrevolution ...
type Saladrevolution struct{}

//FeedMenu ...
func (*Saladrevolution) FeedMenu(cookies []*http.Cookie) (*Menu, error) {
	config := RestaurantConfig{
		URL:                  "https://corporate.caserola.ro/restaurant/saladrevolution",
		AppeteazersSectionID: []int64{499},
		MainsSectionID:       []int64{496},
		DesertsSectionID:     []int64{502},
	}
	return defaultCrawler(config)(cookies)
}

//MakeLunch ...
func (*Saladrevolution) MakeLunch(menu *Menu, noDesert bool) []*Product {
	return MakeLunchByRandom(menu, noDesert)
}
