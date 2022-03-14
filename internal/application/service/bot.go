package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"telegram-go-bot/internal/application/port/out"
	"telegram-go-bot/internal/utils/ziputils"

	"telegram-go-bot/internal/application/model"
)

const (
	tempFolder          = "temp"
	helpCommand         = "/help"
	updateCommand       = "/updates"
	authenticateCommand = "/authenticateme"
	userPosts           = "/postsfrom"
	unloggedHelpText    = "I am a bot, you need to be logged to play with me type /authenticateme [user] [passscode]"
	loggedHelpText      = "/help to get more information about commands.\n/updates to get updates from your social media.\n/authenticateme to authenticate, remember that a session last for only a hour.\n/postsfrom to get posts from a user."
	genericReply        = "I am not so sure about what I should do."
	genericErrorReply   = "I think something went wrong"
	genericSuccessReply = "Ok"
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

func (bot *Bot) ReceiveMessage(message model.ReceivedMessage) []model.ReplyMessage {
	chatId := strconv.Itoa(int(message.ChatId))
	user := message.User
	content := message.Text

	msg := fmt.Sprintf("Message [%s] Received From [%s]", content, user)
	fmt.Println(msg)

	cmd, params := stripCommandAndParams(content)

	reply := make([]model.ReplyMessage, 0)

	switch cmd {
	case helpCommand:
		reply = append(reply, model.ReplyMessage{
			Text: bot.helpText(chatId),
		})
	case authenticateCommand:
		if len(params) >= 2 {
			err := bot.authorizeAndSaveUserCredentialLocally(chatId, params[0], params[1])

			msg := genericSuccessReply
			if err != nil {
				msg = genericErrorReply
			}

			reply = append(reply, model.ReplyMessage{
				Text: msg,
			})
		}
	case updateCommand:
		reply = append(reply, model.ReplyMessage{
			Text: bot.listOfUsersAsSingleText(chatId),
		})
	case userPosts:
		if len(params) >= 1 {
			userPosts, err := bot.checkAuthorizationAndRetrievePostsFromUser(chatId, params[0])
			if err == nil {
				urls := make([]string, 0, len(userPosts)+1)
				for i, p := range userPosts {
					if i > 3 {
						break
					}
					urls = append(urls, p.Image.Url)
				}

				zipFileName, err := ziputils.DownloadFilesAndCompress(urls)

				if err != nil {
					msg := fmt.Sprintf("Unable to compress files due to [%v]", err)
					fmt.Println(msg)

					reply = append(reply, model.ReplyMessage{
						Text: genericErrorReply,
					})
				} else {
					reply = append(reply, model.ReplyMessage{
						Image: &model.ReplyLocalImage{
							FileName: zipFileName,
							FilePath: zipFileName,
						},
					})
				}

			}
		}
	}

	if len(reply) == 0 {
		reply = append(reply, model.ReplyMessage{
			Text: genericReply,
		})
	}

	return reply
}

func (bot *Bot) helpText(chatId string) string {
	token, err := bot.retrieveUserAuthorizationFromLocal(chatId)

	if token == nil || err != nil {
		return unloggedHelpText
	}

	return loggedHelpText
}

func (bot *Bot) authorizeAndSaveUserCredentialLocally(chatId, user, pass string) error {
	accessToken, err := bot.retrieveUserAuthorizationFromRemote(chatId, user, pass)

	if err != nil {
		return err
	}

	err = bot.userRepository.SaveAccessToken(accessToken, chatId)

	if err != nil {
		fmt.Printf("Not possible to save user authorization [%s]\n", chatId)
		return err
	}

	return err
}

func (bot *Bot) listOfUsersAsSingleText(chatId string) string {
	token, err := bot.retrieveUserAuthorizationFromLocal(chatId)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] authorization due to [%s]\n", chatId, err.Error())
		return genericReply
	}

	msg, err := bot.retrieveListOfFollowedUsers(*token)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] updates due to [%s]\n", chatId, err.Error())
		return genericReply
	}

	return msg
}

func (bot *Bot) checkAuthorizationAndRetrievePostsFromUser(chatId, user string) ([]model.UserPost, error) {
	token, err := bot.retrieveUserAuthorizationFromLocal(chatId)

	if err != nil {
		fmt.Printf("Not possible to retrieve user [%s] authorization due to [%s]\n", chatId, err.Error())
		return nil, err
	}

	msg, err := bot.retrievePostsFromUser(*token, user)

	if err != nil {
		fmt.Printf("Not possible to retrieve [%s]'s posts due to [%s]\n", chatId, err.Error())
		return nil, err
	}

	return msg, nil
}

func (bot *Bot) retrieveUserAuthorizationFromLocal(chatId string) (*model.AccessToken, error) {
	accesToken, err := bot.userRepository.RetrieveAccessToken(chatId)

	if err != nil {
		fmt.Printf("Not possible to retrive [%s] token\n", chatId)
		return nil, err
	}

	return accesToken, nil
}

func (bot *Bot) retrieveUserAuthorizationFromRemote(chatId, user, pass string) (*model.AccessToken, error) {
	accessToken, err := bot.scraper.RequestAccessToken(user, pass)

	if err != nil {
		fmt.Printf("Not possible to authorize [%s]\n", chatId)
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

func (bot *Bot) retrievePostsFromUser(accessToken model.AccessToken, userName string) ([]model.UserPost, error) {
	scraper := bot.scraper

	posts, err := scraper.PostsFromUser(accessToken, userName)

	if err != nil {
		return nil, err
	}

	return posts, nil
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
