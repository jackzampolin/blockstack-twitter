# Blockstack Twitter

This is a service that takes a query string and caches the 15 most recent tweets from that query for easy access from a front end.


There is also a sample systemd unit file for easy and safe deployment!

A sample config file (`~/.blockstack-twitter.yaml`)

```yaml
search:         "Verifying my Blockstack"
consumerKey:    "yourTwitterConsumerKey"
consumerSecret: "yourTwitterConsumerSecret"
accessToken:    "yourTwitterAccessToken"
accessSecret:   "yourTwitterAccessSecret"
port:           ":8080"
```
