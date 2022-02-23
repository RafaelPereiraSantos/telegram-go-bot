package api

import (
	"strings"
	"telegram-go-bot/internal/application/model"
	"telegram-go-bot/internal/service"
)

type RedditAdp struct {
	redditIntegration *service.RedditIntegration
}

const userSymbol = "u"

func NewReddtAdp(redditIntegration *service.RedditIntegration) *RedditAdp {
	return &RedditAdp{
		redditIntegration: redditIntegration,
	}
}

func (adp *RedditAdp) RequestAccessToken(user, pass string) (*model.AccessToken, error) {
	return adp.redditIntegration.RequestAccessToken(user, pass)
}

func (adp *RedditAdp) FollowedUsers(accessToken model.AccessToken) ([]model.User, error) {
	subscriptionResponse, err := adp.redditIntegration.FollowedPages(accessToken)

	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0, len(subscriptionResponse.Data.Children)+1)

	for _, child := range subscriptionResponse.Data.Children {
		data := child.Data

		splittedName := strings.Split(data.DispalyNamePrefixed, "/")

		if splittedName[0] != userSymbol {
			continue
		}

		users = append(users, model.User{
			Name: splittedName[1],
		})
	}

	return users, nil
}

func (adp *RedditAdp) PostsFromUser(accessToken model.AccessToken, pageName string) ([]model.UserPost, error) {
	postResponse, err := adp.redditIntegration.PostsFromPage(accessToken, pageName)

	if err != nil {
		return nil, err
	}

	posts := make([]model.UserPost, 0, len(postResponse.Data.Children)+1)

	for _, child := range postResponse.Data.Children {
		data := child.Data
		post := model.UserPost{
			Title: data.Title,
			Image: model.PostImage{
				Url: data.LinkUrl,
			},
		}
		posts = append(posts, post)
	}

	return posts, nil
}
