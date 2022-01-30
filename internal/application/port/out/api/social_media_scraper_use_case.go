package api

import "github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"

type SocialMediaScraper interface {
	FollowedPages() (*model.SubscriptionsResponse, error)
}
