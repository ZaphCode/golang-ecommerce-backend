package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
)

var (
	discordOAuthTokenUrl = "https://discord.com/api/oauth2/token"
	discordOAuthUserUrl  = "https://discord.com/api/users/@me"
)

// Custom types
type DiscordTokens struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordUser struct {
	ID               string      `json:"id"`
	Username         string      `json:"username"`
	Avatar           string      `json:"avatar"`
	AvatarDecoration interface{} `json:"avatar_decoration"`
	Discriminator    string      `json:"discriminator"`
	PublicFlags      int         `json:"public_flags"`
	Flags            int         `json:"flags"`
	Banner           interface{} `json:"banner"`
	BannerColor      interface{} `json:"banner_color"`
	AccentColor      interface{} `json:"accent_color"`
	Locale           string      `json:"locale"`
	MfaEnabled       bool        `json:"mfa_enabled"`
	PremiumType      int         `json:"premium_type"`
	Email            string      `json:"email"`
	Verified         bool        `json:"verified"`
}

func (u DiscordUser) AdaptToUser() (user domain.User) {
	user.Username = u.Username
	user.Email = u.Email
	if u.Avatar != "" {
		url := fmt.Sprintf(
			"https://cdn.discordapp.com/avatars/%s/%s.png",
			u.ID, u.Avatar,
		)
		user.ImageUrl = url
	} else {
		user.ImageUrl = "https://cdn.discordapp.com/embed/avatars/0.png"
	}
	user.Password = utils.RandomString(20)
	user.VerifiedEmail = u.Verified
	return
}

// Constructor
func NewDiscordOAuthService() OAuthService {
	return &discordOAuthServiceImpl{}
}

// Implementation
type discordOAuthServiceImpl struct{}

func (s discordOAuthServiceImpl) GetOAuthUrl() string {
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=",
		config.Get().OAuth.Discord.ClientID,
		config.Get().Api.ServerHost+"/api/auth/discord/callback",
	) + "identify"
}

func (s discordOAuthServiceImpl) GetOAuthUser(code string) (*domain.User, error) {
	tokens, err := s.getDiscordTokens(code)

	if err != nil {
		return nil, err
	}

	discordUser, err := s.getDiscordUser(tokens)

	if err != nil {
		return nil, err
	}

	user := discordUser.AdaptToUser()

	return &user, nil
}

func (s discordOAuthServiceImpl) getDiscordTokens(code string) (*DiscordTokens, error) {
	form := url.Values{}

	form.Add("code", code)
	form.Add("client_id", config.Get().OAuth.Discord.ClientID)
	form.Add("client_secret", config.Get().OAuth.Discord.ClientSecret)
	form.Add("redirect_uri", config.Get().Api.ServerHost+"/api/auth/discord/callback")
	form.Add("grant_type", "authorization_code")
	form.Add("scope", "identify")

	req, err := http.NewRequest(http.MethodPost, discordOAuthTokenUrl, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Println(string(resBody))
		return nil, fmt.Errorf("error fetching to discord api")
	}

	tokens := DiscordTokens{}

	if err := json.Unmarshal(resBody, &tokens); err != nil || tokens.AccessToken == "" {
		return nil, fmt.Errorf("error getting tokens: %w", err)
	}

	return &tokens, nil
}

func (s discordOAuthServiceImpl) getDiscordUser(tokens *DiscordTokens) (*DiscordUser, error) {
	req, err := http.NewRequest(http.MethodGet, discordOAuthUserUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching to google api | %w", err)
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	user := DiscordUser{}

	if err := json.Unmarshal(resBody, &user); err != nil || user.Email == "" {
		return nil, fmt.Errorf("error parsing user. %w", err)
	}

	return &user, nil
}
