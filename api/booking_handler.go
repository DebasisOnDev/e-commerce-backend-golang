package api

import (
	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	bookingStore db.BookingStore
}

func NewBookingHandler(bookingStore db.BookingStore) *BookingHandler {
	return &BookingHandler{
		bookingStore: bookingStore,
	}
}

func (b *BookingHandler) HandleGetBookingById(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := b.bookingStore.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(booking)
}

func (b *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := b.bookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (b *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := b.bookingStore.GetBookingById(c.Context(), id)
	if err != nil {
		return c.JSON(fiber.Map{"error": "invalid user"})
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return c.JSON(fiber.Map{"error": "unauthorized"})
	}
	if booking.UserID != user.ID {
		return c.JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := b.bookingStore.UpdateBooking(c.Context(), id, bson.M{"cancelled": true}); err != nil {
		return c.JSON(fiber.Map{"error": "invalid booking"})
	}
	return c.JSON(fiber.Map{"message": "updated"})

}
