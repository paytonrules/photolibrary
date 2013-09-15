package jobs

import (
	"encoding/json"
	"github.com/paytonrules/photolibrary/thumbnailrequest"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GenerateThumbnailsSuite struct{}

var _ = Suite(&GenerateThumbnailsSuite{})

func (s *GenerateThumbnailsSuite) TestGenerateThumbnailsMakesPostRequest(c *C) {
	var method, host, path string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		host = r.Host
		path = r.URL.Path
	}))
	defer ts.Close()

	job := GenerateThumbnails{Server: ts.URL,
		Duration:  1,
		Directory: "root"}

	job.Run()

	c.Assert(ts.URL, Matches, "http://"+host)
	c.Assert(method, Equals, "POST")
	c.Assert(path, Equals, "/generateThumbnails")
}

func (s *GenerateThumbnailsSuite) TestSendsJsonContentType(c *C) {
	var contentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header["Content-Type"][0]
	}))
	defer ts.Close()

	job := GenerateThumbnails{Server: ts.URL,
		Duration:  1,
		Directory: "root"}

	job.Run()

	c.Assert(contentType, Equals, "text/json")
}

func (s *GenerateThumbnailsSuite) TestSendsJsonOfDirectoryAndDuration(c *C) {
	var req thumbnailrequest.Request
	var err interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
	}))
	defer ts.Close()

	job := GenerateThumbnails{Server: ts.URL,
		Duration:  1,
		Directory: "root"}

	job.Run()

	c.Assert(err, IsNil)
	c.Assert(req.Duration, Equals, 1)
	c.Assert(req.Directory, Equals, "root")
}
