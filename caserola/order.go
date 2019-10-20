package caserola

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
)

type Orders []*Order

const orderHistoryURL = "https://corporate.caserola.ro/api/order/me"

func (ors Orders) DidIOrderToday() bool {
	n := time.Now().UTC()
	for _, o := range ors {
		if d := o.DatePlaced; d.Year() == n.Year() && d.YearDay() == n.YearDay() && o.State != "cancelled" {
			return true
		}
	}
	return false
}

func FeedOrders(cookies []*http.Cookie) (Orders, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", orderHistoryURL, nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var oh OrderFeed
	if err := json.NewDecoder(resp.Body).Decode(&oh); err != nil {
		return nil, err
	}

	return oh.Items, nil
}

func PlaceOrder(products []*Product, cookies []*http.Cookie) (bool, error) {
	buff := strOrder(products)
	body := bytes.NewBufferString(buff.String())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://corporate.caserola.ro/api/order", body)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Oops! I could not place the order! :(\n")
		return false, err
	}
	defer resp.Body.Close()
	if resp.Status != "200" {
		return false, nil
	}
	return true, nil
}

func strOrder(products []*Product) bytes.Buffer {
	// type orderItem struct {
	// 	p         *caserola.Product
	// 	delimeter string
	// }
	type order struct {
		Items   []*Product
		LastIdx int
	}
	// var ps []orderItem
	// for idx, p := range products {
	// 	ps = append(ps, orderItem{p, idx == len(products) -1 ? ",": "."})
	// }
	todaysOrder := order{products, len(products) - 1}
	data, err := ioutil.ReadFile("order.template")
	orderTpl := template.Must(template.New("order").Delims("(", ")").Parse(string(data)))
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = orderTpl.Execute(&buf, todaysOrder)
	if err != nil {
		panic(err)
	}
	return buf
}

//!-
