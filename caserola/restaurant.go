package caserola

import (
	"math/rand"
	"net/http"
)

type RestaurantCrawler func([]*http.Cookie) (*Menu, error)
type Restaurant interface {
	FeedMenu([]*http.Cookie) (*Menu, error)
	MakeLunch(*Menu) []*Product
}

func Shuffle(a []*Product) []*Product {
	swap := func(i, j int) {
		a[i], a[j] = a[j], a[i]
	}
	rand.Shuffle(len(a)-1, swap)
	return a
}

func defaultMakeLunch(menu *Menu) []*Product {
	res := make([]*Product, 0, 3)
	a := Shuffle(menu.Appeteazers)
	res = append(res, a[:1]...)
	m := Shuffle(menu.Mains)
	res = append(res, m[:1]...)
	d := Shuffle(menu.Deserts)
	res = append(res, d[:1]...)
	return res
}

func defaultCrawler(config RestaurantConfig) RestaurantCrawler {
	filter := func(sID int64) productFilter { return func(p *Product) bool { return p.SectionID == sID } }
	f := func(cookies []*http.Cookie) (*Menu, error) {
		products, err := crawlProducts(cookies, config.URL)
		if err != nil {
			return nil, err
		}
		res := Menu{
			Appeteazers: filterProducts(products, filter(config.AppeteazersSectionID)),
			Mains:       filterProducts(products, filter(config.MainsSectionID)),
			Deserts:     filterProducts(products, filter(config.DesertsSectionID)),
		}
		return &res, nil
	}
	return f
}
