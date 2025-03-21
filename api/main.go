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
	"github.com/influxdata/influxdb-client-go/v2/api/query"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
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
		fmt.Fprint(w, "<h1>Hello, twi-fav-api!</h1>")
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
		fmt.Fprintf(os.Stderr, "'limit' query parameter is invalid: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `from(bucket: "twi_fav")
		|> range(start: 0, stop: %s)
		|> filter(fn: (r) => r._measurement == "liked_tweet")
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> sort(columns: ["_time"], desc: true)
		|> limit(n: %d)`
	tweets, err := queryDB(r.Context(), fmt.Sprintf(query, earliestTime, limit))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tweets); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode response: %v", err)
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
			Text:      stringByKey(rec, "text"),
			Link:      stringByKey(rec, "link"),
			CreatedAt: toJST(stringByKey(rec, "createdAt")),
		})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return tweets, nil
}

func stringByKey(rec *query.FluxRecord, key string) string {
	v, ok := rec.ValueByKey(key).(string)
	return lo.Ternary(ok, v, "")
}

func toJST(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.In(jst).Format("2006/01/02 15:04")
}
