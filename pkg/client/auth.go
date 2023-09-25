package client

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/yanmhlv/pcloud/pkg/util"
)

type AuthResult struct {
	Result
	Auth string `json:"auth"`
}

// SetToken client; https://docs.pcloud.com/methods/intro/authentication.html
func (c *pCloudClient) SetToken(tokenStr string) {
	c.Auth = &tokenStr
}

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

	result := AuthResult{}

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return err
	}

	if result.Result.Result != 0 {
		return errors.New(result.Error)
	}

	c.SetToken(result.Auth)
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
