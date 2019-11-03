package caserola

import (
	"math/rand"
	"net/http"
)

var (
	capPriceDesert     = 15.00
	capPriceAppeteazer = 15.00
	capPriceMain       = 35.00
)

//RestaurantCrawler ...
type RestaurantCrawler func([]*http.Cookie) (*Menu, error)

//Restaurant ...
type Restaurant interface {
	FeedMenu([]*http.Cookie) (*Menu, error)
	MakeLunch(*Menu, bool) []*Product
}

func shuffle(a []*Product) []*Product {
	swap := func(i, j int) {
		a[i], a[j] = a[j], a[i]
	}
	rand.Shuffle(len(a)-1, swap)
	return a
}

func randomProduct(a []*Product) (*Product, bool) {
	if n := len(a); n != 0 {
		r := rand.Intn(n)
		return a[r], true
	}
	return nil, false
}

//MakeLunchByRandom ...
func MakeLunchByRandom(menu *Menu, noDesert bool) []*Product {
	res := make([]*Product, 0, 3)
	if p, ok := randomProduct(menu.Appeteazers); ok {
		res = append(res, p)
	}
	if p, ok := randomProduct(menu.Mains); ok {
		res = append(res, p)
	}
	if noDesert {
		return res
	}
	if p, ok := randomProduct(menu.Deserts); ok {
		res = append(res, p)
	}
	return res
}

//MakeLunchByShuffle ...
func MakeLunchByShuffle(menu *Menu, noDesert bool) []*Product {
	res := make([]*Product, 0, 3)
	a := shuffle(menu.Appeteazers)
	res = append(res, a[:1]...)
	m := shuffle(menu.Mains)
	res = append(res, m[:1]...)
	if noDesert {
		return res
	}
	d := shuffle(menu.Deserts)
	res = append(res, d[:1]...)
	return res
}

func defaultCrawler(config RestaurantConfig) RestaurantCrawler {
	sFilter := func(sID []int64) productFilter {

		mapS := make(map[int64]bool)
		for _, id := range sID {
			mapS[id] = true
		}

		return func(p *Product) bool {
			_, found := mapS[p.SectionID]
			return found
		}
	}
	capFilter := func(cap float64) productFilter {
		return func(p *Product) bool {
			return p.Price <= cap
		}
	}
	f := func(cookies []*http.Cookie) (*Menu, error) {
		products, err := crawlProducts(cookies, config.URL)
		if err != nil {
			return nil, err
		}
		res := Menu{
			Appeteazers: filterProducts(products, sFilter(config.AppeteazersSectionID), capFilter(capPriceAppeteazer)),
			Mains:       filterProducts(products, sFilter(config.MainsSectionID), capFilter(capPriceMain)),
			Deserts:     filterProducts(products, sFilter(config.DesertsSectionID), capFilter(capPriceDesert)),
		}
		return &res, nil
	}
	return f
}
