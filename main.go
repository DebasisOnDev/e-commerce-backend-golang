package main

import (
	"context"
	"log"
	"os"

	"github.com/DebasisOnDev/E-Commerce-Backend/api"
	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	productStore := db.NewMongoProductStore(client)
	cartStore := db.NewMongoCartStore(client)
	bookingStore := db.NewMongoBookingStore(client, cartStore)
	userStore := db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	productHandler := api.NewProducthandler(db.NewMongoProductStore(client))
	cartHandler := api.NewCartHandler(db.NewMongoCartStore(client), productStore, bookingStore)
	bookingHandler := api.NewBookingHandler(db.NewMongoBookingStore(client, cartStore))
	authHandler := api.NewAuthHandler(userStore)
	app := fiber.New()
	app.Use(logger.New())
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	auth := app.Group("/api/auth")
	admin := apiv1.Group("/admin", api.AdminAuth)

	//auth handler and auth routes

	auth.Post("/register", authHandler.HandleRegisterUser)
	auth.Post("/login", authHandler.HandleLogInUser)
	auth.Get("/logout", authHandler.HandleLogOutUser)

	//user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//product handlers

	apiv1.Get("/product", productHandler.HandleGetProducts)
	apiv1.Get("/product/:id", productHandler.HandleGetProductByID)
	apiv1.Post("/product", productHandler.HandleInsertProduct)
	apiv1.Delete("/product/:id", productHandler.HandleDeleteProduct)

	//cart handlers
	apiv1.Get("/product/:id/add", cartHandler.HandleAddToCart)
	apiv1.Get("/cart", cartHandler.HandleGetFromCart)
	apiv1.Delete("/cart/:id", cartHandler.HandleRemoveFromCart)

	//booking handlers
	apiv1.Get("/cart/book", cartHandler.HandleBookProductFromCart)
	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBookingById)
	apiv1.Get("/bookings/:id/cancel", bookingHandler.HandleCancelBooking)

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	app.Listen(listenAddress)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
