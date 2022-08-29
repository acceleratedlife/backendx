package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

// schedule a function to be called every minute
// returns channel to stop repeating
// to set up
//  runEveryMinute(func(time2 time.Time) {
//	  lgr.Printf("INFO Tick2 at", time2)
//   })
// to stop:   done <- true

func runEveryMinute(db *bolt.DB) (done chan bool) {
	ticker := time.NewTicker(4 * time.Second)
	done = make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				coinGecko(db)
			}
		}
	}()

	return done
}

func coinGecko(db *bolt.DB) (err error) {
	u, _ := url.ParseRequestURI("https://api.coingecko.com/api/v3/simple/price")
	q := u.Query()
	q.Set("ids", KeyCoins)
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		lgr.Printf(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lgr.Printf(err.Error())
		return
	}

	var decodedResp CoinGecko

	err = json.Unmarshal(body, &decodedResp)
	if err != nil {
		lgr.Printf(err.Error())
		return
	}
	fmt.Println(decodedResp) // to parse out your value
	fmt.Println(decodedResp.Bitcoin.Usd)

	return

}
