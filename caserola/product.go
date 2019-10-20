package caserola

import (
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

type productFilter func(*Product) bool

func crawlProducts(cookies []*http.Cookie, url string) ([]*Product, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return findProducts(nil, doc), nil
}

func findProducts(products []*Product, n *html.Node) []*Product {
	if n.Type == html.ElementNode && n.Data == "div" {
		var p Product
		addMe := false
		for _, a := range n.Attr {
			switch a.Key {
			case "itemtype":
				addMe = a.Val == "http://schema.org/Product"
			case "data-id":
				p.ID, _ = strconv.ParseInt(a.Val, 10, 64)
			case "data-restaurant-id":
				p.RestaurantID, _ = strconv.ParseInt(a.Val, 10, 64)
			case "data-section-id":
				p.SectionID, _ = strconv.ParseInt(a.Val, 10, 64)
			default:
				break
			}
		}
		if addMe {
			products = append(products, &p)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		products = findProducts(products, c)
	}
	return products
}

func filterProducts(products []*Product, filters ...productFilter) []*Product {
	filteredProducts := make([]*Product, 0, len(products))
	for _, r := range products {
		keep := true
		for _, f := range filters {
			if !f(r) {
				keep = false
				break
			}
		}

		if keep {
			filteredProducts = append(filteredProducts, r)
		}
	}

	return filteredProducts
}
