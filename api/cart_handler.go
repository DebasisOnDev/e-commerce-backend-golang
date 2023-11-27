package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type CartHandler struct {
	cartStore    db.CartStore
	productStore db.ProductStore
	bookingStore db.BookingStore
}

func NewCartHandler(cartStore db.CartStore, productStore db.ProductStore, bookingStore db.BookingStore) *CartHandler {
	return &CartHandler{
		cartStore:    cartStore,
		productStore: productStore,
		bookingStore: bookingStore,
	}
}

func (p *CartHandler) HandleAddToCart(c *fiber.Ctx) error {
	id := c.Params("id")
	if p == nil {
		return errors.New("CartHandler is nil")
	}

	if p.productStore == nil {
		return errors.New("ProductStore is nil")
	}
	product, err := p.productStore.GetProductByID(c.Context(), id)
	if err != nil {
		return fmt.Errorf("GetProductByID error: %v", err)
	}
	item, err := types.CreateCartItemfromProduct(*product)
	if err != nil {
		return fmt.Errorf("CreateCartItemfromProduct error: %v", err)
	}

	cartItem, err := p.cartStore.InsertToCart(c.Context(), item)
	if err != nil {
		return err
	}
	if cartItem == nil {
		return c.JSON(map[string]string{"error": "item already in cart"})
	}

	return c.JSON(map[string]string{"added to cart": id})
}

func (p *CartHandler) HandleGetFromCart(c *fiber.Ctx) error {
	items, err := p.cartStore.GetFromCart(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(items)
}

func (p *CartHandler) HandleRemoveFromCart(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := p.cartStore.RemoveFromCart(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(map[string]string{"cart": "item removed from cart"})
}

func calculateTotalPrice(cartitem []*types.CartItem) float64 {
	var totalPrice float64
	for _, item := range cartitem {
		totalPrice += item.Price
	}
	return totalPrice
}

func ProductsInCart(cartitem []*types.CartItem) []primitive.ObjectID {
	var cartItemIds []primitive.ObjectID
	for _, item := range cartitem {
		cartItemIds = append(cartItemIds, item.ItemId)
	}
	return cartItemIds
}

func (p *CartHandler) HandleBookProductFromCart(c *fiber.Ctx) error {

	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	cartItems, err := p.cartStore.GetFromCart(c.Context())
	if err != nil {
		return err
	}

	totalPrice := calculateTotalPrice(cartItems)

	orderDate := time.Now()
	expextedDeliveryDate := orderDate.Add(time.Hour * 24 * 7)

	booking := types.Booking{
		UserID:               user.ID,
		Products:             ProductsInCart(cartItems),
		TotalPrice:           totalPrice,
		OrderDate:            orderDate,
		ExpectedDeliveryDate: expextedDeliveryDate,
		Cancelled:            false,
	}
	book, err := p.bookingStore.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	if err := p.HandleBookingProduct(c, book); err != nil {
		fmt.Println(err)
		return c.JSON(fiber.Map{"error": "cant handle product store", "details": err.Error()})
	}
	if err := p.cartStore.EmptyTheCart(c.Context(), bson.M{}); err != nil {
		return c.JSON(fiber.Map{"error": "items still in cart"})
	}

	return c.JSON(book)
}

func (p *CartHandler) HandleBookingProduct(contx *fiber.Ctx, booking *types.Booking) error {
	fmt.Println("entering the function")
	for _, item := range booking.Products {
		pid := item.Hex()
		fmt.Println(pid)
		product, err := p.productStore.GetProductByID(context.TODO(), string(pid))
		if err != nil {
			return err
		}
		count := product.Quantity - 1
		if count <= 0 {
			return contx.JSON(fiber.Map{"error": "can not update product in store"})
		}
		if err := p.productStore.UpdateProduct(contx.Context(), pid, bson.M{"quantity": count}); err != nil {
			return contx.JSON(fiber.Map{"error": "product is not updated "})
		}

	}
	return nil
}
