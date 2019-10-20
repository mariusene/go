package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"lunch/caserola"
)

//!-
var restaurants = map[string]caserola.Restaurant{"rawdiacorp": &caserola.Rawdiacorp{}}

func main() {
	var cf caserola.Config
	r, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("Panicking at reading config.json :(\n")
		r.Close()
		panic(err)
	}
	json.NewDecoder(r).Decode(&cf)
	r.Close()
	loginStr := fmt.Sprintf("email=%s&password=%s", cf.Email, cf.Pwd)
	ok, cookies, err := login(loginStr)
	if ok == false {
		fmt.Printf("Oops! I cannot log in.:(\n")
	}

	orders, err := caserola.FeedOrders(cookies)
	if err != nil {
		fmt.Printf("Oops! I could not read the orders so I will not place any today :(\n")
	} else {
		if orders.DidIOrderToday() {
			fmt.Printf("I see you've already ordered. Job well done!\n")
		} else {
			fmt.Printf("You did not order. I'm preparing the order for you!\n")
		}
	}

	if r, found := restaurants[cf.Restaurant]; found {
		if menu, err := r.FeedMenu(cookies); err == nil {
			lunch := r.MakeLunch(menu)
			if ok, _ := caserola.PlaceOrder(lunch, cookies); ok {
				fmt.Printf("w00t w00t w00t w00t! Your lunch is set! Enjoy!:)\n")
			}
		} else {
			fmt.Printf("Oops! I could not read the restaurant menu so I will not place any today :(\n")
		}
	} else {
		fmt.Printf("Sorry! You're preferred restaurant is not yet implemented!\n")
	}
}

func login(loginStr string) (bool, []*http.Cookie, error) {
	client := &http.Client{}
	requestBody := bytes.NewBufferString(loginStr)
	req, _ := http.NewRequest("POST", "https://corporate.caserola.ro/api/login", requestBody)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Oops! I could not login! :(\n")
		return false, nil, err
	}
	defer resp.Body.Close()
	if resp.Status != "200" {
		return false, nil, err
	}
	return true, resp.Cookies(), nil
}
