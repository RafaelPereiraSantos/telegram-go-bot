package out

import "telegram-go-bot/internal/application/model"

type SocialMediaScraper interface {
	RequestAccessToken(user, pass string) (*model.AccessToken, error)
	FollowedUsers(accessToken model.AccessToken) ([]model.User, error)
	PostsFromUser(accessToken model.AccessToken, pageName string) ([]model.UserPost, error)
}
