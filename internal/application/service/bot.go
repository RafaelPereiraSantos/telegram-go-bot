package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"telegram-go-bot/internal/application/port/out"

	"telegram-go-bot/internal/application/model"
)

const (
	helpCommand         = "/help"
	updateCommand       = "/updates"
	authenticateCommand = "/authenticateme"
	userPosts           = "/postsfrom"
	unloggedHelpText    = "I am a bot, you need to be logged to play with me type /authenticateme [user] [passscode]"
	loggedHelpText      = "/help to get more information about commands.\n/updates to get updates from your social media.\n/authenticateme to authenticate, remember that a session last for only a hour.\n/postsfrom to get posts from a user."
	genericErrorReply   = "Sorry, can't do that now"
)

type Bot struct {
	scraper        out.SocialMediaScraper
	userRepository out.UserRepository
}

func NewBot(scraper out.SocialMediaScraper, userRepository out.UserRepository) *Bot {
	return &Bot{
		scraper:        scraper,
		userRepository: userRepository,
	}
}

func (bot *Bot) ReceiveMessage(message model.Message) string {
	chatId := strconv.Itoa(int(message.ChatId))
	user := message.User
	content := message.Text

	msg := fmt.Sprintf("Message [%s] Received From [%s]", content, user)
	fmt.Println(msg)

	cmd, params := stripCommandAndParams(content)

	switch cmd {
	case helpCommand:
		return bot.replyWithHelp(chatId)
	case authenticateCommand:
		if len(params) < 2 {
			return genericErrorReply
		}
		return bot.authorizeUserAndReply(chatId, params[0], params[1])
	case updateCommand:
		return bot.replyWithListOfFollowedUsers(chatId)
	case userPosts:
		if len(params) < 1 {
			return genericErrorReply
		}
		return bot.replyWithUserPosts(chatId, params[0])
	}

	return "stop bothering me"
}

func (bot *Bot) replyWithHelp(chatId string) string {
	token, err := bot.retrieveUserAuthorization(chatId)

	if token == nil || err != nil {
		return unloggedHelpText
	}

	return loggedHelpText
}

func (bot *Bot) authorizeUserAndReply(chatId, user, pass string) string {
	_, err := bot.authorizeUserAndSaveToken(chatId, user, pass)

	if err != nil {
		return genericErrorReply
	}

	return "Good to see you back sir"
}

func (bot *Bot) replyWithListOfFollowedUsers(chatId string) string {
	token, err := bot.retrieveUserAuthorization(chatId)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] authorization due to [%s]\n", chatId, err.Error())
		return genericErrorReply
	}

	msg, err := bot.retrieveListOfFollowedUsers(*token)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] updates due to [%s]\n", chatId, err.Error())
		return genericErrorReply
	}

	return msg
}

func (bot *Bot) replyWithUserPosts(chatId, user string) string {
	token, err := bot.retrieveUserAuthorization(chatId)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] authorization due to [%s]\n", chatId, err.Error())
		return genericErrorReply
	}

	msg, err := bot.retrievePostsFromUser(*token, user)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] updates due to [%s]\n", chatId, err.Error())
		return genericErrorReply
	}

	return msg
}

func (bot *Bot) retrieveUserAuthorization(chatId string) (*model.AccessToken, error) {
	accesToken, err := bot.userRepository.RetrieveAccessToken(chatId)

	if err != nil {
		fmt.Printf("Not possible to retrive [%s] token\n", chatId)
		return nil, err
	}

	return accesToken, nil
}

func (bot *Bot) authorizeUserAndSaveToken(chatId, user, pass string) (*model.AccessToken, error) {
	accessToken, err := bot.scraper.RequestAccessToken(user, pass)

	if err != nil {
		fmt.Printf("Not possible to authorize [%s]\n", chatId)
		return nil, err
	}

	err = bot.userRepository.SaveAccessToken(accessToken, chatId)

	if err != nil {
		fmt.Printf("Not possible to save user authorization [%s]\n", chatId)
		return nil, err
	}

	return accessToken, nil
}

func (bot *Bot) retrieveListOfFollowedUsers(accessToken model.AccessToken) (string, error) {
	scraper := bot.scraper

	users, err := scraper.FollowedUsers(accessToken)

	if err != nil {
		return "", err
	}

	var msgBuffer bytes.Buffer

	for _, u := range users {
		msgBuffer.WriteString(u.Name)
		msgBuffer.WriteString("\n")
	}

	return msgBuffer.String(), nil
}

func (bot *Bot) retrievePostsFromUser(accessToken model.AccessToken, userName string) (string, error) {
	scraper := bot.scraper

	posts, err := scraper.PostsFromUser(accessToken, userName)

	if err != nil {
		return "", err
	}

	var msgBuffer bytes.Buffer

	for _, p := range posts {
		msgBuffer.WriteString(p.Title)
		msgBuffer.WriteString("\n")
	}

	return msgBuffer.String(), nil
}

func stripCommandAndParams(message string) (string, []string) {
	words := strings.Split(message, " ")
	size := len(words)
	firstWord := words[0]
	cmd := ""

	if strings.HasPrefix(firstWord, "/") {
		cmd = firstWord
	}

	if size == 1 {
		return cmd, make([]string, 0)
	}

	return cmd, words[1:size]
}
