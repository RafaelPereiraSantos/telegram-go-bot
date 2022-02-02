package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"telegram-go-bot/internal/application/model"
)

type (
	RedditIntegration struct {
		appId    string
		appToken string
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

var (
	ErrTokenExpired = errors.New("Invalid authorization token")
)

func NewRedditIntegration(appId, appToken string) *RedditIntegration {
	return &RedditIntegration{
		appId:    appId,
		appToken: appToken,
	}
}

func (integration *RedditIntegration) RequestAccessToken(user, pass string) (*model.AccessToken, error) {
	fmt.Println("Requesting access token")

	client := &http.Client{}

	req, err := integration.buildOauth2Request(user, pass)

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

	return &model.AccessToken{
		User:        user,
		Value:       tokenResp.AccessToken,
		ExpiresIn:   tokenResp.ExpiresIn,
		RequestedAt: time.Now().Unix(),
	}, nil
}

func (integration *RedditIntegration) FollowedPages(accessToken model.AccessToken) (*model.SubscriptionsResponse, error) {
	if !isAccessTokenValid(accessToken) {
		return nil, ErrTokenExpired
	}

	url := redditOauthHost + mySubscriptionsPath
	req, err := buildGetWithHeaders(url, accessToken.User, integration.appId, accessToken.Value)

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

func isAccessTokenValid(token model.AccessToken) bool {
	expireAt := token.ExpiresIn + token.RequestedAt
	currentTime := time.Now().UnixMicro()

	return currentTime < expireAt
}

func (integration *RedditIntegration) buildOauth2Request(user, pass string) (*http.Request, error) {
	url := redditMainHost + accessTokenPath

	payload := bytes.NewBuffer([]byte(buildOauth2PayloadWithPassword(user, pass)))

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	req.SetBasicAuth(integration.appId, integration.appToken)
	addUserAgentHeader(&req.Header, user, integration.appId)

	return req, nil
}

func buildOauth2PayloadWithPassword(user, pass string) string {
	return fmt.Sprintf("grant_type=password&username=%s&password=%s", user, pass)
}

func buildGetWithHeaders(url, user, appId, authorization string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	addDefaultHeadersWithAuthentication(&req.Header, user, appId, authorization)

	return req, nil
}

func addDefaultHeadersWithAuthentication(header *http.Header, user, appId, authorization string) {
	header.Set("Authorization", fmt.Sprintf("bearer %s", authorization))
	header.Set("Accept", contentType)
	addUserAgentHeader(header, user, appId)
}

func addUserAgentHeader(header *http.Header, user, appId string) {
	header.Set("User-Agent", buildUserAgent(user, appId))
}

func buildUserAgent(user, appId string) string {
	return fmt.Sprintf(
		"%s:%s:%s (by /u/%s)",
		userAgent,
		appId,
		version,
		user,
	)
}
