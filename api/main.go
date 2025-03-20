package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/compute/metadata"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type tweetData struct {
	Time      time.Time `json:"time"`
	Text      string    `json:"text"`
	Link      string    `json:"link"`
	CreatedAt string    `json:"createdAt"`
}

func init() {
	if !metadata.OnGCE() {
		if err := godotenv.Load("../.env.yaml"); err != nil {
			panic(err)
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	})
	mux.HandleFunc("GET /liked-tweets", getLikedTweetsHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mux)
}

func getLikedTweetsHandler(w http.ResponseWriter, r *http.Request) {
	earliestTime := r.URL.Query().Get("earliestTime")
	if earliestTime == "" {
		earliestTime = "now()"
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := `from(bucket: "twi_fav")
		|> range(start: 1970-01-01T00:00:00Z, stop: %s)
		|> filter(fn: (r) => r._measurement == "liked_tweet")
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> sort(columns: ["_time"], desc: true)
		|> limit(n: %d)`
	tweets, err := queryDB(r.Context(), fmt.Sprintf(q, earliestTime, limit))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tweets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func queryDB(ctx context.Context, query string) ([]tweetData, error) {
	client := influxdb2.NewClient(os.Getenv("INFLUXDB2_URL"), os.Getenv("INFLUXDB2_TOKEN"))
	defer client.Close()
	api := client.QueryAPI(os.Getenv("INFLUXDB2_ORG"))
	res, err := api.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	tweets := []tweetData{}
	for res.Next() {
		rec := res.Record()
		tweets = append(tweets, tweetData{
			Time:      rec.Time(),
			Text:      rec.ValueByKey("text").(string),
			Link:      rec.ValueByKey("link").(string),
			CreatedAt: toJST(rec.ValueByKey("createdAt").(string)),
		})
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	return tweets, nil
}

func toJST(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.In(jst).Format("2006/01/02 15:04")
}
