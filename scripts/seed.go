package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/DebasisOnDev/E-Commerce-Backend/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("database connected")

	userStore := db.NewMongoUserStore(client)
	productStore := db.NewMongoProductStore(client)
	cartStore := db.NewMongoCartStore(client)

	product := fixtures.AddProduct(productStore, "biskut", "dark biscuit", 2.5, 5)
	fmt.Println(product)
	cartItem, err := fixtures.AddProductToCart(cartStore, 10, 4.1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cartItem)

	for i := 1; i <= 5; i++ {
		user := fixtures.AddUser(userStore, fmt.Sprintf("John%d", i), fmt.Sprintf("Doe%d", i), false)

		fmt.Println(user)
	}

	adminUser := fixtures.AddUser(userStore, "admin", "user", true)
	fmt.Println(adminUser)

}
