// Copyright © 2017 Jack Zampolin <jack.zampolin@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"sort"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TwitterClient is a twitter client
type TwitterClient struct {
	Client *twitter.Client
	Tweets []twitter.Tweet
}

type Tweets []twitter.Tweet

// NewTwitterClient does the things
func NewTwitterClient() *TwitterClient {
	config := oauth1.NewConfig(viper.GetString("consumerKey"), viper.GetString("consumerSecret"))
	token := oauth1.NewToken(viper.GetString("accessToken"), viper.GetString("accessSecret"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	s, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: viper.GetString("search"),
	})
	if err != nil {
		panic(err)
	}
	sort.Sort(Tweets(s.Statuses))
	if err != nil {
		panic(err)
	}
	return &TwitterClient{
		Client: client,
		Tweets: s.Statuses,
	}
}

func (t *TwitterClient) makeStream() *twitter.Stream {
	params := &twitter.StreamFilterParams{
		Track:         []string{viper.GetString("search")},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := t.Client.Streams.Filter(params)
	if err != nil {
		panic(err)
	}
	return stream
}

func (t *TwitterClient) readStream() {
	demux := twitter.NewSwitchDemux()
	stream := t.makeStream()
	demux.Tweet = func(tweet *twitter.Tweet) {
		log.Println(tweet.User.Name, tweet.Text)
		t.addTweet(*tweet)
	}
	log.Println("Handling Channel")
	go demux.HandleChan(stream.Messages)
}

func (t *TwitterClient) addTweet(tweet twitter.Tweet) {
	t.Tweets = append(t.Tweets, tweet)
	sort.Sort(Tweets(t.Tweets))
	t.Tweets = t.Tweets[:len(t.Tweets)-1]
}

func (t *TwitterClient) handleTwitter(w http.ResponseWriter, r *http.Request) {
	tweets, err := json.Marshal(t.Tweets)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(tweets)
}

func (t Tweets) Len() int {
	return len(t)
}
func (t Tweets) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t Tweets) Less(i, j int) bool {
	return t[i].CreatedAt > t[j].CreatedAt
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		t := NewTwitterClient()
		t.readStream()
		mux := http.NewServeMux()
		mux.HandleFunc("/", t.handleTwitter)
		handler := cors.Default().Handler(mux)
		log.Println("Server listening on port", viper.GetString("port"))
		http.ListenAndServe(viper.GetString("port"), handler)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
