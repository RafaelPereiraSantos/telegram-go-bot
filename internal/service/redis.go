package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"telegram-go-bot/internal/application/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisIntegration struct {
	client  *redis.Client
	context context.Context
}

const (
	accessTokenSuffix = "_token"
	attemptsSuffix    = "_attemps"
	validateKey       = "true"
)

var attemptExpirationTime int64 = 60

func NewRedisIntegration(address, password string, db int) (*RedisIntegration, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	if rdb == nil {
		return nil, fmt.Errorf("was not possible to connect with redis")
	}

	return &RedisIntegration{
		client:  rdb,
		context: context.Background(),
	}, nil
}

func (db *RedisIntegration) SaveAccessToken(accessToken *model.AccessToken, userId string) error {
	acccessTokeyKey := buildAccessTokenKey(userId)
	value, err := json.Marshal(&accessToken)

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible marshal access token due to %s", err.Error()))
		return err
	}

	err = db.setKeyValue(acccessTokeyKey, string(value), accessToken.ExpiresIn)

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to save access token due to %s", err.Error()))
	}

	return err
}

func (db *RedisIntegration) RetrieveAccessToken(userId string) (*model.AccessToken, error) {
	acccessTokeyKey := buildAccessTokenKey(userId)

	value, err := db.getValue(acccessTokeyKey)

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to retrieve access token due to %s", err.Error()))
		return nil, err
	}

	var accessToken model.AccessToken
	err = json.Unmarshal([]byte(value), &accessToken)

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to unmarshal access token due to %s", err.Error()))
		return nil, err
	}

	return &accessToken, nil
}

func (db *RedisIntegration) RetrieveUserLoginAttempt(userId string) (int, error) {
	key := buildAttemptKey(userId)

	value, err := db.getValue(key)

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to retrieve [%s] attempt due to [%s]", userId, err.Error()))
	}

	attemp, err := strconv.Atoi(value)

	if err != nil {
		fmt.Println(fmt.Sprintf("Attemp of [%s] bad format [%s]", userId, err.Error()))
		return 0, err
	}

	return attemp, nil
}

func (db *RedisIntegration) IncreaseUserLoginAttempt(userId string) error {
	key := buildAttemptKey(userId)

	currentAttempStr, err := db.getValue(key)

	errorHandler := func(err error) {
		fmt.Println(fmt.Sprintf("Not possible to save [%s] attempt due to [%s]", userId, err.Error()))
	}

	if err != nil {
		if err != redis.Nil {
			errorHandler(err)
			return err
		} else {
			err := db.setKeyValue(key, "1", attemptExpirationTime)

			if err != nil {
				errorHandler(err)
				return err
			}

			return nil
		}
	}

	currentAttemp, _ := strconv.Atoi(currentAttempStr)
	currentAttemp += 1

	err = db.setKeyValue(key, strconv.Itoa(currentAttemp), attemptExpirationTime)

	if err != nil {
		errorHandler(err)
		return err
	}

	return nil
}

func (db *RedisIntegration) ResetUserLoginAttemp(userId string) error {
	err := db.client.Expire(db.context, buildAttemptKey(userId), 0).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to reset [%s] attempt due to [%s]", userId, err.Error()))
	}
	return err
}

func (db *RedisIntegration) setKeyValue(k, v string, expTimeSeconds int64) error {
	cli := db.client
	expiration := time.Second * time.Duration(expTimeSeconds)
	err := cli.Set(db.context, k, v, expiration).Err()
	return err
}

func (db *RedisIntegration) getValue(k string) (string, error) {
	cli := db.client
	value, err := cli.Get(db.context, k).Result()
	return value, err
}

func buildAttemptKey(userId string) string {
	return userId + attemptsSuffix
}

func buildAccessTokenKey(userId string) string {
	return userId + accessTokenSuffix
}
