package caserola

import "net/http"

type Rawdiacorp struct{}

func (*Rawdiacorp) FeedMenu(cookies []*http.Cookie) (*Menu, error) {
	config := RestaurantConfig{
		URL:                  "https://corporate.caserola.ro/restaurant/rawdiacorp",
		AppeteazersSectionID: []int64{512},
		MainsSectionID:       []int64{513},
		DesertsSectionID:     []int64{515},
	}
	return defaultCrawler(config)(cookies)
}
func (*Rawdiacorp) MakeLunch(menu *Menu) []*Product {
	return MakeLunchByRandom(menu)
}
