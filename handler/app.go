package handler

import (
	"go-crud/auth"
	"go-crud/database"
	"go-crud/middleware"
	"go-crud/repository"
	"go-crud/service"
	"log"
	"net/http"
	"os"

	// _ "annisa-salon/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartApp() {
	db, err := database.InitDb()
	if err != nil {
		log.Fatal("Eror Db Connection")
	}

	secretKey := os.Getenv("SECRET_KEY")

	router := gin.Default()

	router.Static("/static", "./static")

	router.LoadHTMLGlob("templates/*")

	// Render register page
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "data_produk.html", nil)
	})

	router.GET("/input-produk", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form_input_produk.html", nil)
	})

	router.GET("/lihat-data", func(c *gin.Context) {
		c.HTML(http.StatusOK, "lihat_produk.html", nil)
	})

	router.GET("/edit-produk", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form_edit_produk.html", nil)
	})

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Origin , Accept , X-Requested-With , Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization"},
		AllowMethods:    []string{"POST, OPTIONS, GET, PUT, DELETE"},
	}))

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := auth.NewUserAuthService()
	authService.SetSecretKey(secretKey)
	userHandler := NewUserHandler(userService, authService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := NewProductHandler(productService)

	router.POST("/products", middleware.AuthMiddleware(authService, userService), productHandler.CreateProductHandler)
	router.GET("/products/:id", productHandler.GetProductByIDHandler)
	router.GET("/products", productHandler.GetAllProductsHandler)
	router.PUT("/products/:id", middleware.AuthMiddleware(authService, userService), productHandler.UpdateProductHandler)
	router.DELETE("/products/:id", middleware.AuthMiddleware(authService, userService), productHandler.DeleteProductHandler)

	user := router.Group("api/user")
	user.POST("/register", userHandler.RegisterUser)
	user.POST("/login", userHandler.Login)

	router.Run(":8080")

}
