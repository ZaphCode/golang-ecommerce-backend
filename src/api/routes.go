package api

import (
	addressHandler "github.com/ZaphCode/clean-arch/src/api/handlers/address"
	authHandler "github.com/ZaphCode/clean-arch/src/api/handlers/auth"
	cardHandler "github.com/ZaphCode/clean-arch/src/api/handlers/card"
	categoryHandler "github.com/ZaphCode/clean-arch/src/api/handlers/category"
	orderHandler "github.com/ZaphCode/clean-arch/src/api/handlers/order"
	productHandler "github.com/ZaphCode/clean-arch/src/api/handlers/product"
	userHandler "github.com/ZaphCode/clean-arch/src/api/handlers/user"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/utils"
)

func (s *Server) CreateAuthRoutes(
	authHdlr *authHandler.AuthHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/auth")
	r.Get("/:provider/url", authHdlr.GetOAuthUrl)
	r.Get("/:provider/callback", authHdlr.SignInWihOAuth)
	r.Get("/refresh", authHdlr.RefreshToken)
	r.Get("/me", authMdlw.AuthRequired, authHdlr.GetAuthUser)
	r.Get("/signout", authHdlr.SignOut)
	r.Post("/signin", authHdlr.SignIn)
	r.Post("/signup", authHdlr.SignUp)
}

func (s *Server) CreateUserRoutes(
	usrHdlr *userHandler.UserHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/user")
	r.Get("/all", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), usrHdlr.GetUsers)
	r.Get("/get/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), usrHdlr.GetUser)
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), usrHdlr.CreateUser)
	r.Put("/update/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), usrHdlr.UpdateUser)
	r.Delete("/delete/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), usrHdlr.DeleteUser)
}

func (s *Server) CreateProductRoutes(
	prodHdlr *productHandler.ProductHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/product")
	r.Get("/all", prodHdlr.GetProducts)
	r.Get("/get/:id", prodHdlr.GetProduct)
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.CreateProduct)
	r.Put("/update/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.UpdateProduct)
	r.Delete("/delete/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.AdminRole), prodHdlr.DeleteProduct)
}

func (s *Server) CreateCategoryRoutes(
	catHdlr *categoryHandler.CategoryHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/category")
	r.Get("/all", catHdlr.GetCategories)
	r.Post("/create", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), catHdlr.CreateCategory)
	r.Delete("/delete/:id", authMdlw.AuthRequired, authMdlw.RoleRequired(utils.ModeratorRole), catHdlr.DeleteCategory)
}

func (s *Server) CreateAddressesRoutes(
	addrHdlr *addressHandler.AddressHandler,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/address")
	r.Get("/list", authMdlw.AuthRequired, addrHdlr.GetUserAddress)
	r.Post("/create", authMdlw.AuthRequired, addrHdlr.CreateAddress)
	r.Put("/update/:id", authMdlw.AuthRequired, addrHdlr.UpdateAddress)
	r.Delete("/delete/:id", authMdlw.AuthRequired, addrHdlr.DeleteAddress)
}

func (s *Server) CreateCardRoutes(
	cardHdlr *cardHandler.CardHandler,
	paymMdlw *middlewares.PaymentMiddleware,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/card")
	r.Get("/list", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, cardHdlr.GetUserCards)
	r.Post("/save", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, cardHdlr.SaveUserCard)
	r.Delete("/remove/:id", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, cardHdlr.RemoveUserCard)
}

func (s *Server) CreateOrderRoutes(
	ordHdlr *orderHandler.OrderHandler,
	paymMdlw *middlewares.PaymentMiddleware,
	authMdlw *middlewares.AuthMiddleware,
) {
	r := s.app.Group("/api/order")
	r.Get("/list", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, ordHdlr.GetOrders)
	r.Post("/new", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, ordHdlr.CreateOrder)
	r.Get("/get/:id", authMdlw.AuthRequired, paymMdlw.CustomerIDRequired, ordHdlr.GetOrderById)
}
