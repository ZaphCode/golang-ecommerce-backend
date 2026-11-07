package utils

import "errors"

//* OAuth providers

const (
	GoogleProvider  = "google"
	DiscordProvider = "discord"
	GithubProvider  = "github"
)

func GetOAuthProviders() []string {
	return []string{GoogleProvider, DiscordProvider, GithubProvider}
}

//* User roles

const (
	UserRole      = "user"
	ModeratorRole = "moderator"
	AdminRole     = "admin"
)

func GetUserRoles() []string {
	return []string{UserRole, ModeratorRole, AdminRole}
}

//* Firestore collection names

const (
	UserColl  = "users"
	ProdColl  = "products"
	AddrColl  = "addresses"
	OrderColl = "orders"
	CategColl = "categories"
)

//* Errors

var (
	ErrNotFound = errors.New("resource not found")
)

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCancelled  = "cancelled"
	StatusShipped    = "shipped"
	StatusCompleted  = "completed"
)
