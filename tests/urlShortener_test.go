package tests

import (
	"net/url"
	"testing"
	"url-shortener/internals/http-server/handlers/url/save"
	"url-shortener/internals/lib/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	// "url-shortener/internal/http-server/handlers/url/save"
	// "url-shortener/internal/lib/random"
	// "url-shortener/internals/http-server/handlers/url/save"
	// "url-shortener/internals/lib/random"
)

const (
	host = "localhost:8081"
)

func TestURLShortener_HappyPath(t *testing.T) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, url.String())

	e.POST("/url/").WithJSON(save.Request{
		URL:   gofakeit.URL(),
		Alias: random.NewRandomString(10),
	}).
		WithBasicAuth("myname", "mypass").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}
