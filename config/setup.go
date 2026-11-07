package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type api struct {
	Port               string `json:"port"`
	ClientOrigin       string `json:"client_origin"`
	ServerHost         string `json:"server_host"`
	VerificationSecret string `json:"verification_secret"`
	ChangepassSecret   string `json:"changepass_secret"`
	AccessTokenSecret  string `json:"access_token_secret"`
	RefreshTokenSecret string `json:"refresh_token_secret"`
	AccessTokenHeader  string `json:"access_token_header"`
	RefreshTokenHeader string `json:"refresh_token_header"`
	RefreshTokenCookie string `json:"refresh_token_cookie"`
}

type oauthServices struct {
	Google  oauthCredentials `json:"google"`
	Discord oauthCredentials `json:"discord"`
	Github  oauthCredentials `json:"github"`
}

type oauthCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type storage struct {
	Bucket       string `json:"bucket"`
	UploadFolder string `json:"upload_folder"`
}

type smtp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type stripe struct {
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
}

type Config struct {
	Api     api           `json:"api"`
	OAuth   oauthServices `json:"oauth"`
	Smtp    smtp          `json:"smtp"`
	Storage storage       `json:"storage"`
	Stripe  stripe        `json:"stripe"`
}

var (
	config      Config
	firebaseApp *firebase.App = nil
)

func MustLoadConfig(path string) {
	file, err := os.ReadFile(path + "/app_config.json")

	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(file, &config); err != nil {
		panic(err)
	}

	if config.Api.Port == "" || config.Api.ServerHost == "" {
		panic("json config not readed")
	}
}

func MustLoadFirebaseConfig(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		panic("cannot find firebase configuration file")
	}

	ctx := context.Background()

	opt := option.WithCredentialsFile(path + "/firebase_config.json")

	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		panic(fmt.Errorf("error connecting to firebase: %s", err))
	}

	firebaseApp = app
}

func Get() Config {
	return config
}

func GetFirebaseApp() *firebase.App {
	return firebaseApp
}

func IsProduction() bool {
	return os.Getenv("NODE_ENV") == "production"
}
