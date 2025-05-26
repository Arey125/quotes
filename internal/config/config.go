package config

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"quotes/internal/users"
)

type config struct {
	Port          int
	Db            string
	Oauth         users.OauthConfig
}

func getString(name string) string {
	res := os.Getenv(name)
	if res == "" {
		panic(fmt.Sprintf("%s environment variable is not set", name))
	}
	return res
}

func getInt(name string) int {
	stringRes := getString(name)
	res, err := strconv.Atoi(stringRes)
	if err != nil {
		panic(fmt.Sprintf("cannot parse %s environment variable as an integer", name))
	}
	return res
}

func Get() config {
	port := getInt("PORT")
	db := getString("DB")

	id := getString("CLIENT_ID")
	secret := getString("CLIENT_SECRET")
	callbackUrl := getString("CLIENT_CALLBACK_URL")

    sessionSecret := getString("SESSION_SECRET")

	return config{
		Port: port,
		Db:   db,
		Oauth: users.OauthConfig{
			Id:          id,
			Secret:      secret,
			CallbackUrl: callbackUrl,
            SessionSecret: sessionSecret,
		},
	}
}
