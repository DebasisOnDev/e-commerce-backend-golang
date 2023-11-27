package api

import (
	"fmt"

	"github.com/DebasisOnDev/E-Commerce-Backend/db"
	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productStore db.ProductStore
}

func NewProducthandler(productStore db.ProductStore) *ProductHandler {
	return &ProductHandler{
		productStore: productStore,
	}
}

func (p *ProductHandler) HandleGetProducts(c *fiber.Ctx) error {
	products, err := p.productStore.GetProducts(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(products)
}
func (p *ProductHandler) HandleGetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := p.productStore.GetProductByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(product)
}

func (p *ProductHandler) HandleInsertProduct(c *fiber.Ctx) error {
	var prod types.CreateProductParams
	if err := c.BodyParser(&prod); err != nil {
		return err
	}

	errors := prod.Validate()
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}
	product, err := types.CreateProductFromParams(prod)
	if err != nil {
		fmt.Println("cant create product")
		return err

	}
	newprod, err := p.productStore.InsertProduct(c.Context(), product)
	if err != nil {
		fmt.Println("cant insert product")
		return err
	}
	return c.JSON(newprod)
}

func (p *ProductHandler) HandleDeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := p.productStore.DeleteProduct(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(map[string]string{"Deleted": id})
}
