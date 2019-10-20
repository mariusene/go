package caserola

import "net/http"

type Rawdiacorp struct{}

var config = RestaurantConfig{
	ID:                   "rawdiacorp",
	URL:                  "https://corporate.caserola.ro/restaurant/rawdiacorp",
	AppeteazersSectionID: 512,
	MainsSectionID:       513,
	DesertsSectionID:     515,
}

func (*Rawdiacorp) FeedMenu(cookies []*http.Cookie) (*Menu, error) {
	return defaultCrawler(config)(cookies)
}
func (*Rawdiacorp) MakeLunch(menu *Menu) []*Product {
	return defaultMakeLunch(menu)
}
