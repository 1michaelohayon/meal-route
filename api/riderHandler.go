package api

import (
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderHandler struct {
	store db.Store
}

func NewRiderHandler(s db.Store) *RiderHandler {
	return &RiderHandler{
		store: s,
	}
}

func (h *RiderHandler) GetAll(ctx *fiber.Ctx) error {
	riders, err := h.store.Rider.Get(ctx.Context())
	if err != nil {
		if err == mongo.ErrNilDocument {
			return ctx.JSON([]types.Rider{})
		}
		return err
	}
	return ctx.JSON(riders)
}

func (h *RiderHandler) Add(ctx *fiber.Ctx) error {
	fpID := ctx.Params("foodPorviderID")
	fp, err := h.store.Fp.GetById(ctx.Context(), fpID)
	if err != nil {
		return ctx.Status(404).JSON(map[string]string{"Add error": "food provider not found"})
	}

	var params types.Rider
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
		//TODO better parse erros
	}

	//TODO VALIDATION

	params.FoodProviderID = fp.ID
	insertedRider, err := h.store.Rider.Insert(ctx.Context(), &params)
	if err != nil {
		return err
	}
	if err := h.store.Fp.PutRider(ctx.Context(), fp.ID, insertedRider.ID); err != nil {
		return err
	}

	return ctx.JSON(insertedRider)
}

func (h *RiderHandler) GetOne(ctx *fiber.Ctx) error {
	fpID := ctx.Params("foodPorviderID")
	riderID := ctx.Params("id")
	if _, err := h.store.Fp.GetById(ctx.Context(), fpID); err != nil {
		return ctx.Status(404).JSON(map[string]string{"Rider->GetOne error": "food provider not found"})
	}

	rider, err := h.store.Rider.GetById(ctx.Context(), riderID)
	if err != nil {
		return fiber.ErrNotFound
	}
	return ctx.JSON(rider)
}

/* not ready */
func (h *RiderHandler) AssignUser(ctx *fiber.Ctx) error {
	//get user from context

	fpID := ctx.Params("foodPorviderID")
	//riderID := ctx.Params("id")
	if _, err := h.store.Fp.GetById(ctx.Context(), fpID); err != nil {
		return ctx.Status(404).JSON(map[string]string{"Rider->GetOne error": "food provider not found"})
	}
	return nil
}

/*
TODO:
	validation for adding riders, maybe in a middleware
	better erros


*/
