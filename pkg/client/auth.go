package client

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/yanmhlv/pcloud/pkg/util"
)

// Login client; https://docs.pcloud.com/methods/intro/authentication.html
func (c *pCloudClient) Login(username string, password string) error {
	values := url.Values{
		"getauth":  {"1"},
		"username": {username},
		"password": {password},
	}

	buf, err := ConvertToBuffer(c.Client.Get(util.UrlBuilder("userinfo", values)))
	if err != nil {
		return err
	}

	result := struct {
		Auth   string `json:"auth"`
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	c.Auth = &result.Auth
	return nil
}

// Logout client; https://docs.pcloud.com/methods/auth/logout.html
func (c *pCloudClient) Logout() error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	if err := CheckResult(c.Client.Get(util.UrlBuilder("logout", values))); err != nil {
		return err
	}

	c.Auth = nil
	return nil
}
