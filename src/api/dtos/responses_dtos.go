package dtos

//? ---------------------------------------------
//? All this dtos are for documentation porpurses
//? ---------------------------------------------

//* ------- BASE OK ----------

type RespOKDTO struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Data retrived!"`
}

//* -------- USERS ------------

type UserRespOKDTO struct {
	RespOKDTO
	Data UserDTO `json:"data"`
}

type UsersRespOKDTO struct {
	RespOKDTO
	Data []UserDTO `json:"data"`
}

//* -------- PRODUCTS ----------

type ProductRespOKDTO struct {
	RespOKDTO
	Data ProductDTO `json:"data"`
}

type ProductsRespOKDTO struct {
	RespOKDTO
	Data []ProductDTO `json:"data"`
}

//* -------- CARDS ----------

type CardRespOKDTO struct {
	RespOKDTO
	Data CardDTO `json:"data"`
}

type CardsRespOKDTO struct {
	RespOKDTO
	Data []CardDTO `json:"data"`
}

//* -------- ADDRESS ----------

type AddressRespOKDTO struct {
	RespOKDTO
	Data AddressDTO `json:"data"`
}

type AddressesRespOKDTO struct {
	RespOKDTO
	Data []AddressDTO `json:"data"`
}

//* ------ CATEGORIES ----------

type CategoryRespOKDTO struct {
	RespOKDTO
	Data CategoryDTO `json:"data"`
}

type CategoriesRespOKDTO struct {
	RespOKDTO
	Data []CategoryDTO `json:"data"`
}

// * ------ ORDERS ----------
type OrderRespOKDTO struct {
	RespOKDTO
	Data OrderDTO `json:"data"`
}

type DetailedOrderRespOKDTO struct {
	RespOKDTO
	Data DetailedOrderDTO `json:"data"`
}

type OrdersRespOKDTO struct {
	RespOKDTO
	Data []OrderDTO `json:"data"`
}

//* --------- AUTH -------------

type URLRespOKDTO struct {
	RespOKDTO
	Data string `json:"data" example:"https://google.com/oauth/pulse"`
}

type TokenRespOKDTO struct {
	RespOKDTO
	Data string `json:"data" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

type SignInRespOKDTO struct {
	RespOKDTO
	Data struct {
		RefreshToken string  `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxNDE0NTYzNiIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.Bl1Rpmk-BrbqtgJA6F9pTAuiOlaPLpdDQ7MJvZ7URSU"`
		User         UserDTO `json:"user"`
		AccessToken  string  `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	} `json:"data"`
}

//! --------- ERROR ------------

// RespErr represents a simple error response
type RespErrDTO struct {
	Status  string `json:"status" example:"failure"`
	Message string `json:"message" example:"something went wrong"`
}

// RespDetailErr represents a detailed error response
type DetailRespErrDTO struct {
	Status  string `json:"status" example:"failure"`
	Message string `json:"message" example:"something went wrong"`
	Detail  string `json:"detail,omitempty" example:"error querying the database"`
}

// RespValErr represents the validation error response
type ValidationRespErrDTO struct {
	Status  string `json:"status" example:"failure"`
	Message string `json:"message" example:"one or more field are invalid"`
	Errors  []struct {
		Field   string `json:"field" example:"Email"`
		Message string `json:"message" example:"invalid email"`
	} `json:"errors"`
}

type AuthRespErrDTO struct {
	Status  string `json:"status" example:"failure"`
	Message string `json:"message" example:"invalid access token"`
	Detail  string `json:"detail,omitempty" example:"the token is expired"`
}
