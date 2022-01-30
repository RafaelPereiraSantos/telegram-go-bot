package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"
)

type (
	RedditIntegration struct {
		appId       string
		appToken    string
		user        RedditUser
		accessToken *model.AccessToken
	}

	RedditUser struct {
		UserName string
		Password string
	}
)

const (
	userAgent           = "golang-bot"
	version             = "0.0.1"
	redditMainHost      = "https://www.reddit.com"
	redditOauthHost     = "https://oauth.reddit.com"
	accessTokenPath     = "/api/v1/access_token"
	mySubscriptionsPath = "/subreddits/mine/subscriber"
	contentType         = "application/json"
)

func NewRedditIntegration(appId, appToken string, user RedditUser) *RedditIntegration {
	return &RedditIntegration{
		appId:    appId,
		appToken: appToken,
		user:     user,
	}
}

func (integration *RedditIntegration) FollowedPages() (*model.SubscriptionsResponse, error) {
	if !integration.isAccessTokenValid() {
		err := integration.updateAccessToken()

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	req, err := integration.buildGetWithHeaders(redditOauthHost + mySubscriptionsPath)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	rawBody := resp.Body

	defer rawBody.Close()

	body, err := io.ReadAll(rawBody)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var subResp model.SubscriptionsResponse

	err = json.Unmarshal(body, &subResp)

	fmt.Println(subResp)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &subResp, nil
}

func (integration *RedditIntegration) isAccessTokenValid() bool {
	accessToken := integration.accessToken

	if accessToken == nil {
		return false
	}

	expireAt := accessToken.ExpiresIn + accessToken.RequestedAt
	currentTime := time.Now().UnixMicro()

	return currentTime < expireAt
}

func (integration *RedditIntegration) updateAccessToken() error {
	resp, err := integration.RequestAccessToken()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	integration.accessToken = &model.AccessToken{
		Value:       resp.AccessToken,
		ExpiresIn:   resp.ExpiresIn,
		RequestedAt: time.Now().UnixMicro(),
	}

	return nil
}

func (integration *RedditIntegration) RequestAccessToken() (*model.AccessTokenResponse, error) {
	fmt.Println("Requesting access token")

	client := &http.Client{}

	req, err := integration.buildOauth2Request()

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	rawBody := resp.Body

	defer rawBody.Close()

	body, err := io.ReadAll(rawBody)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var tokenResp model.AccessTokenResponse

	err = json.Unmarshal(body, &tokenResp)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tokenResp, nil
}

func (integration *RedditIntegration) buildGetWithHeaders(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	integration.addDefaultHeaders(&req.Header)

	return req, nil
}

func (integration *RedditIntegration) addDefaultHeaders(header *http.Header) {
	value := fmt.Sprintf("bearer %s", integration.accessToken.Value)
	header.Set("Authorization", value)
	header.Set("Accept", contentType)

	integration.addUserAgentHeader(header)
}

func (integration *RedditIntegration) buildOauth2Request() (*http.Request, error) {
	url := redditMainHost + accessTokenPath

	payload := bytes.NewBuffer([]byte(integration.buildOauth2PayloadWithPassword()))

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	req.SetBasicAuth(integration.appId, integration.appToken)
	integration.addUserAgentHeader(&req.Header)

	return req, nil
}

func (integration *RedditIntegration) addUserAgentHeader(header *http.Header) {
	header.Set("User-Agent", integration.buildUserAgent())
}

func (integration *RedditIntegration) buildOauth2PayloadWithPassword() string {
	return fmt.Sprintf("grant_type=password&username=%s&password=%s", integration.user.UserName, integration.user.Password)
}

func (integration *RedditIntegration) buildUserAgent() string {
	return fmt.Sprintf(
		"%s:%s:%s (by /u/%s)",
		userAgent,
		integration.appId,
		version,
		integration.user,
	)
}
