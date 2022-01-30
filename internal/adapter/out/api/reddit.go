package api

import (
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/service"
)

type RedditAdp struct {
	redditIntegration *service.RedditIntegration
}

func NewReddtAdp(redditIntegration *service.RedditIntegration) *RedditAdp {
	return &RedditAdp{
		redditIntegration: redditIntegration,
	}
}

func (adp *RedditAdp) FollowedPages() (*model.SubscriptionsResponse, error) {
	return adp.redditIntegration.FollowedPages()
}
