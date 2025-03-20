package webhook

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type tweet struct {
	Text      string
	Link      string
	CreatedAt time.Time
}

func init() {
	functions.HTTP("SaveLikedTweet", SaveLikedTweet)
}

func SaveLikedTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	createdAt, err := time.Parse("January 2, 2006 at 03:04PM", r.FormValue("createdAt"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Error parsing createdAt", http.StatusBadRequest)
		return
	}
	tweet := tweet{
		Text:      r.FormValue("text"),
		Link:      r.FormValue("link"),
		CreatedAt: createdAt,
	}
	fmt.Printf("%#v\n", tweet)
	if err := writeDB(r.Context(), tweet); err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Error writing to DB", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeDB(ctx context.Context, tweet tweet) error {
	client := influxdb2.NewClient(os.Getenv("INFLUXDB2_URL"), os.Getenv("INFLUXDB2_TOKEN"))
	defer client.Close()
	api := client.WriteAPIBlocking(os.Getenv("INFLUXDB2_ORG"), os.Getenv("INFLUXDB2_BUCKET"))
	p := influxdb2.NewPointWithMeasurement("liked_tweet").
		AddField("text", tweet.Text).
		AddField("link", tweet.Link).
		AddField("createdAt", tweet.CreatedAt).
		SetTime(time.Now())
	return api.WritePoint(ctx, p)
}
