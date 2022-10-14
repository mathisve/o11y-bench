// this is horrible hacked together go
// will improve later

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
)

var baseURL = "http://af76109e6c3334215adf8cd1932e6fc3-1235338525.us-east-1.elb.amazonaws.com:16686/api/traces?end=%s&start=%s&limit=%s&lookback=%s&maxDuration&minDuration&service=%s"

// services generated by omnition/synthetic-load-generator
var services = []string{
	"adservice",
	"cartservice",
	"checkoutservice",
	"currencyservice",
	"emailservice",
	"frontend",
	"paymentservice",
	"productcatalogservice",
	"recommendationservice",
	"shippingservice"}

// lookback options (later converted to hours)
var lookback = []int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
}

// struct used for composing unique URLs
// generated in generateURL
type URLComposition struct {
	Start    string
	End      string
	Limit    string
	Lookback string
	Service  string
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	// amount of unique requests
	var total = 500

	// amount of goroutines
	var concurrency = 4

	// slice that holds all the timings for the requests
	var times []float64
	var timesLock sync.Mutex

	// int that holds amount of time the request came back with a status code other than 200
	var not200count int64
	var not200countLock sync.Mutex

	var wg sync.WaitGroup

	log.Println("Starting requests ...")

	for i := 0; i < concurrency; i += 1 {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			for j := 0; j < int(math.Round(float64(total/concurrency))); j += 1 {

				comp := generateURL()
				url := formatURL(comp)

				var totalTime int64
				start := time.Now().UnixMilli()

				resp, err := http.Get(url)
				if err != nil {
					log.Panic(err)
				}

				if resp.StatusCode != http.StatusOK {

					totalTime = time.Now().UnixMilli() - start

					not200countLock.Lock()
					not200count += 1
					not200countLock.Unlock()

					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Println(err)
					}

					log.Println(url)
					log.Println(string(body))
				} else {
					totalTime = time.Now().UnixMilli() - start
				}

				timesLock.Lock()
				times = append(times, float64(totalTime))
				timesLock.Unlock()

				log.Println(resp.Status)
				log.Printf("%+v\n", comp)
				log.Println(url)
				log.Println(totalTime, " ms")
				fmt.Println("------------")

				time.Sleep(time.Millisecond * 200)
			}

			wg.Done()
		}(&wg)
	}

	wg.Wait()

	log.Println(times)
	log.Printf("not 200 code: %d/%d", not200count, total)

	mean, _ := stats.Mean(times)
	fmt.Printf("mean -> %f", mean)

	var percentiles = []float64{10, 25, 50, 75, 90, 95, 99}

	for _, percentile := range percentiles {

		p, err := stats.Percentile(times, percentile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%f -> %f\n", percentile, p)

	}
}

func generateURL() URLComposition {
	comp := URLComposition{}

	now := time.Now()

	// add one so loopback is always bigger than 0
	lookback := rand.Intn(len(lookback)) + 1

	// add 3 0's to match with Jaegers query timestamp
	// i'm sure there is a cleaner way to do this
	// we subtract the lookback time by adding the negative of lookback
	comp.Start = fmt.Sprintf("%d000", int(now.Add(time.Hour*time.Duration(-lookback)).UnixMilli()))
	comp.End = fmt.Sprintf("%d000", int(now.UnixMilli()))
	comp.Lookback = fmt.Sprintf("%dh", lookback)

	// pick random service
	comp.Service = services[rand.Intn(len(services))]

	// pick random limit under 20
	comp.Limit = fmt.Sprint(rand.Intn(20))

	return comp
}

func formatURL(comp URLComposition) (url string) {
	return fmt.Sprintf(baseURL, comp.End, comp.Start, comp.Limit, comp.Lookback, comp.Service)
}
