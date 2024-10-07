package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const feedURL = "https://habr.com/ru/rss/best/daily/?fl=ru"

func TestParse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, feedURL, nil)
	assert.NoError(t, err, "client: could not create request")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "client: error making http request")

	b, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err, "read response body")

	var fi FeedInfo
	err = xml.Unmarshal(b, &fi)
	assert.NoError(t, err, "client: unmarshalling xml feed")

	feed, err := Parse(fi, 1)
	assert.NoError(t, err, "parsing RSS feed")
	assert.NotEmpty(t, feed, "check if decoded feed is empty")
}
