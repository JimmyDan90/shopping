package api

import (
	"github.com/gin-gonic/gin"
	"log"
	cartApi "shopping/api/cart"
	categoryApi "shopping/api/category"
	orderApi "shopping/api/order"
	productApi "shopping/api/product"
	userApi "shopping/api/user"
	"shopping/config"
	"shopping/domain/cart"
	"shopping/domain/category"
	"shopping/domain/order"
	"shopping/domain/product"
	"shopping/domain/user"
	"shopping/utils/database_handler"
	"shopping/utils/middleware"
)

// Database 结构体
type Database struct {
	categoryRepository    *category.Repository
	userRepository        *user.Repository
	productRepository     *product.Repository
	cartRepository        *cart.Repository
	cartItemRepository    *cart.ItemRepository
	orderRepository       *order.Repository
	orderedItemRepository *order.OrderedItemRepository
}

var AppConfig = &config.Configuration{}

// CreateDBs 根据配置文件创建数据库
func CreateDBs() *Database {
	cfgFile := "./config/config.yaml"
	conf, err := config.GetAllConfigValues(cfgFile)
	AppConfig = conf
	if err != nil {
		return nil
	}
	if err != nil {
		log.Fatalf("读取配置文件失败，%v", err.Error())
	}
	db := database_handler.NewMySqlDb(AppConfig.DatabaseSettings.DatabaseURI)
	return &Database{
		categoryRepository:    category.NewCategoryRepository(db),
		cartRepository:        cart.NewCartRepository(db),
		userRepository:        user.NewUserRepository(db),
		productRepository:     product.NewProductRepository(db),
		cartItemRepository:    cart.NewCartItemRepository(db),
		orderRepository:       order.NewOrderRepository(db),
		orderedItemRepository: order.NewOrderedItemRepository(db),
	}
}

func RegisterHandlers(r *gin.Engine) {
	dbs := *CreateDBs()
	RegisterUserHandlers(r, dbs)
	RegisterCategoryHandlers(r, dbs)
	RegisterCartHandlers(r, dbs)
	RegisterProductHandlers(r, dbs)
	RegisterOrderHandlers(r, dbs)
}

// RegisterUserHandlers 注册用户控制器
func RegisterUserHandlers(r *gin.Engine, dbs Database) {
	userService := user.NewUserService(*dbs.userRepository)
	userController := userApi.NewUserController(userService, AppConfig)
	userGroup := r.Group("/user")
	userGroup.POST("", userController.CreateUser)
	userGroup.POST("/login", userController.Login)
}

// RegisterCategoryHandlers 注册分类控制器
func RegisterCategoryHandlers(r *gin.Engine, dbs Database) {
	categoryService := category.NewCategoryService(*dbs.categoryRepository)
	categoryController := categoryApi.NewCategoryController(categoryService)
	categoryGroup := r.Group("/category")
	categoryGroup.POST(
		"", middleware.AuthAdminMiddleware(AppConfig.JwtSettings.SecretKey), categoryController.CreateCategory)
	categoryGroup.GET("", categoryController.GetCategories)
	categoryGroup.POST("/upload", middleware.AuthAdminMiddleware(AppConfig.JwtSettings.SecretKey), categoryController.BulkCreateCategory)
}

// RegisterCartHandlers 注册购物车控制器
func RegisterCartHandlers(r *gin.Engine, dbs Database) {
	cartService := cart.NewService(*dbs.cartRepository, *dbs.cartItemRepository, *dbs.productRepository)
	cartController := cartApi.NewCartController(cartService)
	cartGroup := r.Group("/cart", middleware.AuthUserMiddleware(AppConfig.JwtSettings.SecretKey))
	cartGroup.POST("/item", cartController.AddItem)
	cartGroup.PATCH("/item", cartController.UpdateItem)
	cartGroup.GET("/", cartController.GetCart)
}

// RegisterProductHandlers 注册商品控制器
func RegisterProductHandlers(r *gin.Engine, dbs Database) {
	productService := product.NewService(*dbs.productRepository)
	productController := productApi.NewProductController(*productService)
	productGroup := r.Group("/product")
	productGroup.GET("", productController.GetProducts)
	productGroup.POST("", middleware.AuthAdminMiddleware(AppConfig.JwtSettings.SecretKey), productController.CreateProduct)
	productGroup.DELETE("", middleware.AuthAdminMiddleware(AppConfig.JwtSettings.SecretKey), productController.DeleteProduct)
	productGroup.PATCH("", middleware.AuthAdminMiddleware(AppConfig.JwtSettings.SecretKey), productController.UpdateProduct)
}

// RegisterOrderHandlers 注册订单控制器
func RegisterOrderHandlers(r *gin.Engine, dbs Database) {
	orderService := order.NewService(
		*dbs.orderRepository, *dbs.orderedItemRepository, *dbs.productRepository, *dbs.cartRepository, *dbs.cartItemRepository)
	productController := orderApi.NewOrderController(orderService)
	orderGroup := r.Group("/order", middleware.AuthUserMiddleware(AppConfig.JwtSettings.SecretKey))
	orderGroup.POST("", productController.CompleteOrder)
	orderGroup.DELETE("", productController.CancelOrder)
	orderGroup.GET("", productController.GetOrders)
}
