package jsonclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type JsonClient[Req any, Res any] struct {
	conn    *retryablehttp.Client
	headers map[string]string
}

func NewJsonClient[Req any, Res any]() *JsonClient[Req, Res] {
	httpcon := retryablehttp.NewClient()
	httpcon.RetryMax = 3
	httpcon.RetryWaitMax = time.Second * 30
	ret := &JsonClient[Req, Res]{
		conn: httpcon,
		// standard headers for every request
		headers: make(map[string]string),
	}
	ret.headers["Content-Type"] = "application/json"
	return ret
}

func (c *JsonClient[Req, Res]) SetTimeout(timeout time.Duration) {
	c.conn.RetryWaitMax = timeout
}

func (c *JsonClient[Req, Res]) SetRetries(retries int) {
	c.conn.RetryMax = retries
}

func (c *JsonClient[Req, Res]) AddHeader(key string, val string) {
	c.headers[key] = val
}

func (c *JsonClient[Req, Res]) Get(requrl string) (response *Res, err error) {
	u, err := url.Parse(requrl)
	if err != nil {
		return nil, fmt.Errorf("invalid url %s: %w", requrl, err)
	}
	surl := u.String()

	r, err := retryablehttp.NewRequest("GET", surl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request for %q: %v", surl, err)
	}

	for k, v := range c.headers {
		r.Header.Add(k, v)
	}

	res, err := c.conn.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	bodybytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}
	body := string(bodybytes)

	response = new(Res)
	err = json.Unmarshal([]byte(body), response)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v; raw:%s", err, body)
	}

	return response, nil
}

func (c *JsonClient[Req, Res]) Post(requrl string, request *Req) (response *Res, err error) {
	req_json, err := json.MarshalIndent(*request, "", " ")
	if err != nil {
		return nil, fmt.Errorf("error json-encoding request : %v", err)
	}

	r, err := retryablehttp.NewRequest("POST", requrl, bytes.NewBuffer(req_json))
	if err != nil {
		return nil, fmt.Errorf("error creating http request for %q: %v", requrl, err)
	}

	for k, v := range c.headers {
		r.Header.Add(k, v)
	}

	res, err := c.conn.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error posting request: %v", err)
	}
	defer res.Body.Close()

	bodybytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}
	body := string(bodybytes)

	response = new(Res)
	err = json.Unmarshal([]byte(body), response)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v; raw:%s", err, body)
	}

	return response, nil
}

func (c *JsonClient[Req, Res]) PostStream(requrl string, request *Req) (response *Res, err error) {
	req_json, err := json.MarshalIndent(*request, "", " ")
	if err != nil {
		return nil, fmt.Errorf("error json-encoding request : %v", err)
	}

	r, err := retryablehttp.NewRequest("POST", requrl, bytes.NewBuffer(req_json))
	if err != nil {
		return nil, fmt.Errorf("error creating http request for %q: %v", requrl, err)
	}

	for k, v := range c.headers {
		r.Header.Add(k, v)
	}

	res, err := c.conn.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error posting request: %v", err)
	}
	defer res.Body.Close()

	bodybytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}
	body := string(bodybytes)

	response = new(Res)
	err = json.Unmarshal([]byte(body), response)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v; raw:%s", err, body)
	}

	return response, nil
}
