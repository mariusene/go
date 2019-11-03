package caserola

import (
	"math/rand"
	"net/http"
)

type RestaurantCrawler func([]*http.Cookie) (*Menu, error)
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
	filter := func(sID []int64) productFilter {

		mapS := make(map[int64]bool)
		for _, id := range sID {
			mapS[id] = true
		}

		return func(p *Product) bool {
			_, found := mapS[p.SectionID]
			return found
		}
	}
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
