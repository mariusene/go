package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"lunch/caserola"
)

//!-
type cookieChan chan []*http.Cookie
type restaurantLunch struct {
	restaurantKey string
	products      []*caserola.Product
}
type lunchChan chan restaurantLunch

var (
	cf           caserola.Config
	restaurants  = map[string]caserola.Restaurant{"saladrevolution": &caserola.Saladrevolution{}, "rawdiacorp": &caserola.Rawdiacorp{}}
	clients      = make(map[string]cookieChan)
	lunchFeed    = make(lunchChan, len(restaurants))
	placeOder    = make(chan bool)
	yearDayOrder [366]int
	messages     = make(chan string, 100)
)

func main() {
	go printer()
	r, err := os.Open("config.json")
	defer r.Close()
	if err != nil {
		panic(err)
	}
	json.NewDecoder(r).Decode(&cf)

	if _, found := restaurants[cf.Restaurant]; !found {
		fmt.Println("Sorry! You're preferred restaurant is not yet implemented!")
		os.Exit(0)
	}

	clients["history"] = make(cookieChan)
	clients["lunch"] = make(cookieChan)
	for key := range restaurants {
		clients[key] = make(cookieChan)
	}

	go spinner(100 * time.Millisecond)
	go printer()
	go areYouAlive()

	off := []int{1, 1, 1, 1, 1, 3, 2}
	for {
		messages <- fmt.Sprintf("-----------------------------")
		n := time.Now().UTC()
		triggerTime := time.Date(n.Year(), n.Month(), n.Day()+off[int(n.Weekday())], cf.UtcH, cf.UtcM, 0, 0, n.Location())
		untilTomorrow := triggerTime.Sub(n)

		if d := yearDayOrder[n.YearDay()]; d == 1 {
			messages <- fmt.Sprintf("You have lunch for today, so I go to sleep until tomorrow!")
			time.Sleep(untilTomorrow)
		}

		if d := n.Weekday(); d == 6 || d == 0 {
			messages <- fmt.Sprintf("Is weekend, so I will go to sleep until Monday!")
			time.Sleep(untilTomorrow)
		}

		if n.Hour() >= cf.UtcH && n.Minute() >= cf.UtcM {
			messages <- fmt.Sprintf("Time to order:%v!", time.Now().Format("15:04:05"))
			go loginAsMe()
			go checkMyOrders()
			for key := range restaurants {
				go buildRestaurant(key)
			}
			go makeMeLunch()
		}

		time.Sleep(5 * time.Minute)
	}
}

func printer() {
	for msg := range messages {
		fmt.Println(msg)
	}
}

func makeMeLunch() {
	cookies := <-clients["lunch"]
	do := <-placeOder
	messages <- fmt.Sprintf("Waiting for the restaurant feeds.")
	for lunch := range lunchFeed {
		messages <- fmt.Sprintf("I got the feed for:%s", lunch.restaurantKey)
		if do && lunch.restaurantKey == cf.Restaurant {
			messages <- fmt.Sprintf("It's your favorite one. I'm ordering.")
			if ok, _ := caserola.PlaceOrder(lunch.products, cookies); ok {
				messages <- fmt.Sprintf("w00t w00t w00t w00t! Check your email. I hope you like what I've ordered for you!;)\n")
				yearDayOrder[time.Now().UTC().YearDay()] = 1
			}
		}
	}
	messages <- fmt.Sprintf("%v", cookies)
}

func buildRestaurant(key string) {
	cookies := <-clients[key]
	messages <- fmt.Sprintf("Building restaurant: %s", key)
	menu, err := restaurants[key].FeedMenu(cookies)
	if err != nil {
		messages <- fmt.Sprintf("Oops! I could not read the restaurant:{%s} menu so I will not place any today :(", key)
	}
	messages <- fmt.Sprintf("Restaurant:%s has %d-Appeteazers, %d-Mains and %d-Deserts", key, len(menu.Appeteazers), len(menu.Mains), len(menu.Deserts))
	lunch := restaurants[key].MakeLunch(menu)
	lunchFeed <- restaurantLunch{key, lunch}
}

func checkMyOrders() {
	cookies := <-clients["history"]
	orders, err := caserola.FeedOrders(cookies)
	if err != nil {
		messages <- fmt.Sprintf("Oops! I could not read the orders so I will not place any today :(")
		placeOder <- false
	} else {
		if orders.DidIOrderToday() {
			messages <- fmt.Sprintf("I see you've already ordered. Job well done!")
			yearDayOrder[time.Now().UTC().YearDay()] = 1
			placeOder <- false
		} else {
			messages <- fmt.Sprintf("You did not order. I'm going to order for you!")
			placeOder <- true
		}
	}
}

func loginAsMe() {
	loginStr := fmt.Sprintf("email=%s&password=%s", cf.Email, cf.Pwd)

	client := &http.Client{}
	requestBody := bytes.NewBufferString(loginStr)
	req, _ := http.NewRequest("POST", "https://corporate.caserola.ro/api/login", requestBody)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Oops! I could not login! :(\n %v\n", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Oops! I could not login! :(\n %v\n", err))
	}

	var cookies = resp.Cookies()
	for _, val := range clients {
		val <- cookies
	}
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func areYouAlive() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		messages <- fmt.Sprintf("I'm alive and waiting for next time to order!")
	}
}
