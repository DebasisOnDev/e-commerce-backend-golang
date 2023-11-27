package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/DebasisOnDev/E-Commerce-Backend/types"
)

func AddUser(userStore db.UserStore, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserFromParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddProduct(productStore db.ProductStore, n, desc string, price float64, quant int) *types.Product {
	product, err := types.CreateProductFromParams(types.CreateProductParams{
		Name:        n,
		Quantity:    quant,
		Price:       price,
		Description: desc,
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedProduct, err := productStore.InsertProduct(context.TODO(), product)
	if err != nil {
		log.Fatal(err)
	}
	return insertedProduct
}

func AddProductToCart(cartStore db.CartStore, count int, price float64) (*types.CartItem, error) {
	cartItem, err := types.CreateCartItemfromProduct(types.Product{
		Price:    price,
		Quantity: count,
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedCartProd, err := cartStore.InsertToCart(context.TODO(), cartItem)
	if err != nil {
		log.Fatal(err)
	}
	return insertedCartProd, nil
}
