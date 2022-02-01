package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisIntegration struct {
	client  *redis.Client
	context context.Context
}

const (
	adminSuffix    = "_admin"
	attemptsSuffix = "_attemps"
	validateKey    = "true"
)

var adminExpirationTime = time.Minute * 10
var attemptExpirationTime = time.Minute

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

func (db *RedisIntegration) SaveUserAsAdmin(account string) error {
	fmt.Println(fmt.Sprintf("Saving user [%s] as admin", account))

	err := db.setKeyValue(buildAdminKey(account), validateKey, adminExpirationTime)

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to save [%s] as admin due to [%s]", account, err.Error()))
		return err
	}

	return nil
}

func (db *RedisIntegration) isAdminUser(account string) (bool, error) {
	fmt.Println(fmt.Sprintf("Checking if [%s] is admin", account))

	value, err := db.getValue(buildAdminKey(account))

	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to check [%s] privileges due to [%s]", account, err.Error()))
		return false, err
	}

	return value == validateKey, nil
}

func (db *RedisIntegration) IncreaseUserLoginAttempt(account string) error {
	accountKey := buildAttemptKey(account)

	currentAttempStr, err := db.getValue(accountKey)

	errorHandler := func(err error) {
		fmt.Println(fmt.Sprintf("Not possible to save [%s] attempt due to [%s]", account, err.Error()))
	}

	if err != nil {
		if err != redis.Nil {
			errorHandler(err)
			return err
		} else {
			err := db.setKeyValue(accountKey, "1", attemptExpirationTime)

			if err != nil {
				errorHandler(err)
				return err
			}

			return nil
		}
	}

	currentAttemp, _ := strconv.Atoi(currentAttempStr)
	currentAttemp += 1

	err = db.setKeyValue(accountKey, strconv.Itoa(currentAttemp), attemptExpirationTime)

	if err != nil {
		errorHandler(err)
		return err
	}

	return nil
}

func (db *RedisIntegration) ResetUserLoginAttemp(account string) error {
	err := db.client.Expire(db.context, buildAttemptKey(account), 0).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Not possible to reset [%s] attempt due to [%s]", account, err.Error()))
	}
	return err
}

func (db *RedisIntegration) setKeyValue(k, v string, expTime time.Duration) error {
	cli := db.client
	err := cli.Set(db.context, k, v, expTime).Err()
	return err
}

func (db *RedisIntegration) getValue(k string) (string, error) {
	cli := db.client
	value, err := cli.Get(db.context, k).Result()
	return value, err
}

func buildAttemptKey(account string) string {
	return account + attemptsSuffix
}

func buildAdminKey(account string) string {
	return account + adminSuffix
}
