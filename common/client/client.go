package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	url string

	token string
}

func New(url string) *Client {
	return &Client{
		url: strings.TrimRight(url, "/"),
	}
}

func (c *Client) Request(path string, data interface{}, out interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal data failed, %v", err)
	}

	resp, err := http.Post(c.url+path, "application/json", bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("http post failed, %v", err)
	}

	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed, %v", err)
	}

	if out == nil {
		out = struct{}{}
	}

	ret := struct {
		Error string      `json:"error"`
		Data  interface{} `json:"data"`
	}{
		Data: out,
	}

	err = json.Unmarshal(buf, &ret)
	if err != nil {
		return fmt.Errorf("json unmarshal body failed, %v", err)
	}

	if len(ret.Error) > 0 {
		return fmt.Errorf(ret.Error)
	}

	return nil
}

func (c *Client) Token() string {
	return c.token
}

func (c *Client) SetToken(token string) {
	c.token = token
}
