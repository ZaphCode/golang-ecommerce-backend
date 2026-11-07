package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
)

var (
	githubTokenUrl      = "https://github.com/login/oauth/access_token"
	githubUserUrl       = "https://api.github.com/user"
	githubUserEmailsUrl = "https://api.github.com/user/emails"
)

type GithubToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

type GithubUser struct {
	Login      string    `json:"login"`
	ID         int       `json:"id"`
	NodeID     string    `json:"node_id"`
	AvatarURL  string    `json:"avatar_url"`
	GravatarID string    `json:"gravatar_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Bio        string    `json:"bio"`
	Followers  int       `json:"followers"`
	Following  int       `json:"following"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Verified   bool      `json:"verified"` // Custom
}

func (u GithubUser) AdaptToUser() (user domain.User) {
	user.Username = u.Login
	user.Email = u.Email
	user.Password = utils.RandomString(20)
	user.ImageUrl = u.AvatarURL
	user.VerifiedEmail = u.Verified
	return
}

func NewGithubOAuthService() OAuthService {
	return &githubOAuthServiceImpl{}
}

type githubOAuthServiceImpl struct{}

func (s githubOAuthServiceImpl) GetOAuthUrl() string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		config.Get().OAuth.Github.ClientID,
		config.Get().Api.ServerHost+"/api/auth/github/callback",
	)
}

func (s githubOAuthServiceImpl) GetOAuthUser(code string) (*domain.User, error) {
	tokens, err := s.getGitHubTokens(code)

	if err != nil {
		return nil, err
	}

	githubUser, err := s.getGithubUser(tokens)

	if err != nil {
		return nil, err
	}

	user := githubUser.AdaptToUser()

	return &user, nil
}

func (s githubOAuthServiceImpl) getGitHubTokens(code string) (*GithubToken, error) {
	cfg := config.Get()

	body := fmt.Sprintf(`{
		"client_id": "%s", 
		"client_secret": "%s",
		"code": "%s"
	}`, cfg.OAuth.Github.ClientID, cfg.OAuth.Github.ClientSecret, code)

	bodyBytes := bytes.NewBufferString(body)

	req, err := http.NewRequest(http.MethodPost, githubTokenUrl, bodyBytes)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching to github api. %w", err)
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	tokens := GithubToken{}

	if err := json.Unmarshal(resBody, &tokens); err != nil || tokens.AccessToken == "" {
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	return &tokens, nil
}

func (s githubOAuthServiceImpl) getGithubUser(token *GithubToken) (*GithubUser, error) {
	req, err := http.NewRequest(http.MethodGet, githubUserUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

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
		return nil, fmt.Errorf("error fetching to github api")
	}

	user := GithubUser{}

	if err := json.Unmarshal(resBody, &user); err != nil || user.Login == "" {
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	if user.Email == "" {
		email, err := s.getGithubUserEmail(token)

		if err != nil {
			return nil, err
		}

		user.Email = email.Email
		user.Verified = email.Verified
	}

	return &user, nil
}

func (s githubOAuthServiceImpl) getGithubUserEmail(token *GithubToken) (*GithubEmail, error) {
	req, err := http.NewRequest(http.MethodGet, githubUserEmailsUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Println(string(resBody))
		return nil, fmt.Errorf("error fetching to github api (emails)")
	}

	var emails []GithubEmail

	if err := json.Unmarshal(resBody, &emails); err != nil || len(emails) <= 0 {
		return nil, fmt.Errorf("error parsing emails. %w", err)
	}

	for _, email := range emails {
		if email.Primary && email.Email != "" {
			return &email, nil
		}
	}

	return nil, fmt.Errorf("no emails to display")
}
