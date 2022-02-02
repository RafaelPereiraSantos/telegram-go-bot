package out

import "telegram-go-bot/internal/application/model"

type SocialMediaScraper interface {
	RequestAccessToken(user, pass string) (*model.AccessToken, error)
	FollowedPages(accessToken model.AccessToken) (*model.SubscriptionsResponse, error)
}
