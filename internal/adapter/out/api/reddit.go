package api

import (
	"telegram-go-bot/internal/application/model"
	"telegram-go-bot/internal/service"
)

type RedditAdp struct {
	redditIntegration *service.RedditIntegration
}

func NewReddtAdp(redditIntegration *service.RedditIntegration) *RedditAdp {
	return &RedditAdp{
		redditIntegration: redditIntegration,
	}
}

func (adp *RedditAdp) RequestAccessToken(user, pass string) (*model.AccessToken, error) {
	return adp.redditIntegration.RequestAccessToken(user, pass)
}

func (adp *RedditAdp) FollowedPages(accessToken model.AccessToken) (*model.SubscriptionsResponse, error) {
	return adp.redditIntegration.FollowedPages(accessToken)
}
