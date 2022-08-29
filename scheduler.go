package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
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
	ticker := time.NewTicker(60 * time.Second)
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

	var decodedResp map[string]map[string]float32
	err = json.Unmarshal(body, &decodedResp)
	if err != nil {
		lgr.Printf(err.Error())
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cryptosBucket, err := tx.CreateBucketIfNotExists([]byte(KeyCryptos))
		if err != nil {
			lgr.Printf(err.Error())
			return err
		}
		var cryptoInfo openapi.CryptoCb

		for k, v := range decodedResp {
			cryptoInfo.Usd = v["usd"]
			cryptoInfo.UpdatedAt = time.Now().Truncate(time.Second)
			marshal, err := json.Marshal(cryptoInfo)
			if err != nil {
				lgr.Printf(err.Error())
				return err
			}

			err = cryptosBucket.Put([]byte(k), marshal)
			if err != nil {
				lgr.Printf(err.Error())
				return err
			}
		}

		lgr.Printf("Updated Cryptos")

		return nil
	})

	return

}
