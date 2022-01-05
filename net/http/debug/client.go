package debug

import (
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/syb-devs/goth/log"
)

// New allocates and returns a debug.Client
func New(base *http.Client) *Client {
	return &Client{Client: &http.Client{}}
}

// Client wraps an http.Client and logs request/responses for debugging purposes
type Client struct {
	*http.Client
}

// Do executes the given HTTP Request, logging it and the corresponding response
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	dumpRequest(req)
	res, err := c.Client.Do(req)
	dumpResponse(res)
	return res, err

}

// Get performs a GET Request to the given URL, logging the response
func (c *Client) Get(url string) (resp *http.Response, err error) {
	log.Debugf("Client request GET: %s", url)
	res, err := c.Client.Get(url)
	dumpResponse(res)
	return res, err
}

// Post executes a POST Request, logging it and the corresponding response
func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	log.Debug("Client.Post")
	return nil, nil
}

func dumpRequest(req *http.Request) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("request dump: %s", dump)
}

func dumpResponse(res *http.Response) {
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("response dump: %s", dump)
}
