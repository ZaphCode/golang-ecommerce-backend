package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ZaphCode/clean-arch/config"
	_ "github.com/ZaphCode/clean-arch/docs" // Swagger docs
	"github.com/ZaphCode/clean-arch/src/api"
	addressHandler "github.com/ZaphCode/clean-arch/src/api/handlers/address"
	authHandler "github.com/ZaphCode/clean-arch/src/api/handlers/auth"
	cardHandler "github.com/ZaphCode/clean-arch/src/api/handlers/card"
	categoryHandler "github.com/ZaphCode/clean-arch/src/api/handlers/category"
	orderHandler "github.com/ZaphCode/clean-arch/src/api/handlers/order"
	productHandler "github.com/ZaphCode/clean-arch/src/api/handlers/product"
	userHandler "github.com/ZaphCode/clean-arch/src/api/handlers/user"
	"github.com/ZaphCode/clean-arch/src/api/middlewares"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/address"
	"github.com/ZaphCode/clean-arch/src/repositories/category"
	"github.com/ZaphCode/clean-arch/src/repositories/order"
	"github.com/ZaphCode/clean-arch/src/repositories/product"
	"github.com/ZaphCode/clean-arch/src/repositories/user"
	"github.com/ZaphCode/clean-arch/src/services/auth"
	"github.com/ZaphCode/clean-arch/src/services/core"
	"github.com/ZaphCode/clean-arch/src/services/email"
	"github.com/ZaphCode/clean-arch/src/services/payment"
	"github.com/ZaphCode/clean-arch/src/services/validation"
	"github.com/ZaphCode/clean-arch/src/utils"
)

func init() {
	config.MustLoadConfig("./config")
	config.MustLoadFirebaseConfig("./config")
}

// @title Go E-commerce API
// @version 1.0
// @description E-commerce REST API built with Go and Fiber using Clean Architecture.
// @contact.name ZaphCode
// @contact.url https://github.com/ZaphCode
// @host localhost:9000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name X-Access-Token

func main() {
	cfg := config.Get()
	server := api.New()

	setServerConfiguration(server, cfg)

	go server.InitBackgroundTasks()

	if err := server.Start(":" + cfg.Api.Port); err != nil {
		log.Fatal(err)
	}
}

func setServerConfiguration(server *api.Server, cfg config.Config) {
	//* Repos
	var (
		userRepo domain.UserRepository
		prodRepo domain.ProductRepository
		catRepo  domain.CategoryRepository
		addrRepo domain.AddressRepository
		ordRepo  domain.OrderRepository
	)

	if len(os.Args) > 1 && os.Args[1] == "dev" {
		//* Development
		fmt.Println("DEV MODE")
		//userRepo = user.NewMemoryUserRepository(utils.UserAdmin, utils.UserExp1, utils.UserExp2)
		userRepo = user.NewJSONUserRepository("db/users.json")
		//prodRepo = product.NewMemoryProductRepository(utils.ProductExpToDev3, utils.ProductExpToDev2, utils.ProductExpToDev1)
		prodRepo = product.NewJSONProductRepository("db/products.json")
		//catRepo = category.NewMemoryCategoryRepository(utils.CategoryExp1, utils.CategoryExp2, utils.CategoryExp3)
		catRepo = category.NewJSONCategoryRepository("db/categories.json")
		//addrRepo = address.NewMemoryAddressRepository(utils.AddrExp1, utils.AddrExp2)
		addrRepo = address.NewJSONAddressRepository("db/addresses.json")
		//ordRepo = order.NewMemoryOrderRepository()
		ordRepo = order.NewJSONOrderRepository("db/orders.json")
	} else {
		//* Production
		client := utils.GetFirestoreClient(config.GetFirebaseApp())
		userRepo = user.NewFirestoreUserRepository(client, utils.UserColl)
		prodRepo = product.NewFirestoreProductRepository(client, utils.ProdColl)
		catRepo = category.NewFirestoreCategoryRepository(client, utils.CategColl)
		addrRepo = address.NewFirestoreAddressRepository(client, utils.AddrColl)
		ordRepo = order.NewFirestoreOrderRepository(client, utils.OrderColl)
	}

	//* Services
	userSvc := core.NewUserService(userRepo)
	prodSvc := core.NewProductService(prodRepo, catRepo)
	catSvc := core.NewCategoryService(catRepo, prodRepo)
	addrSvc := core.NewAddressService(addrRepo, userRepo)
	ordSvc := core.NewOrderService(ordRepo, addrRepo, prodRepo)
	pmSvc := payment.NewStripePaymentService(cfg.Stripe.SecretKey, userRepo)
	emailSvc := email.NewSmtpEmailService()
	vldSvc := validation.NewValidationService()
	jwtSvc := auth.NewJWTService()

	//* Middlewares
	authMdlw := middlewares.NewAuthMiddleware(jwtSvc)
	paymMdlw := middlewares.NewPaymentMiddleware(pmSvc)

	// Handlers
	usrHdlr := userHandler.NewUserHandler(userSvc, vldSvc)
	addrHdlr := addressHandler.NewAddressHandler(userSvc, addrSvc, vldSvc)
	authHdlr := authHandler.NewAuthHandler(userSvc, emailSvc, jwtSvc, vldSvc)
	prodHdlr := productHandler.NewProductHandler(prodSvc, catSvc, vldSvc)
	catHdlr := categoryHandler.NewCategoryHandler(prodSvc, catSvc, vldSvc)
	cardHdlr := cardHandler.NewCardHandler(userSvc, pmSvc, vldSvc)
	ordHdlr := orderHandler.NewOrderHandler(userSvc, ordSvc, prodSvc, pmSvc, vldSvc)

	//* Setup
	server.SetGlobalMiddlewares()

	//* Routes
	server.CreateAuthRoutes(authHdlr, authMdlw)
	server.CreateUserRoutes(usrHdlr, authMdlw)
	server.CreateProductRoutes(prodHdlr, authMdlw)
	server.CreateCategoryRoutes(catHdlr, authMdlw)
	server.CreateAddressesRoutes(addrHdlr, authMdlw)
	server.CreateCardRoutes(cardHdlr, paymMdlw, authMdlw)
	server.CreateOrderRoutes(ordHdlr, paymMdlw, authMdlw)
}
