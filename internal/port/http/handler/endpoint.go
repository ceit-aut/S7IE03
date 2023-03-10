package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ceit-aut/S7IE03/internal/model"
	"github.com/ceit-aut/S7IE03/internal/port/http/request"
	"github.com/ceit-aut/S7IE03/internal/port/http/response"

	"github.com/gofiber/fiber/v2"
)

var (
	errMaximumEndpoints = errors.New("you have reached the maximum number of endpoints")
	errEmptyAddress     = errors.New("address cannot be empty")
	errSaveEndpoint     = errors.New("failed to save endpoint")
	errRemoveEndpoint   = errors.New("failed to remove endpoint")
)

// RegisterEndpoint will add an endpoint to application.
// takes a endpoint info request.
func (h *Handler) RegisterEndpoint(ctx *fiber.Ctx) error {
	// get all endpoints to find the length
	eps := h.Repositories.Endpoints.GetUserEndpoints(ctx.Locals("username").(string))
	if len(eps) > h.Limit {
		return errMaximumEndpoints
	}

	// create a user request
	var (
		id      string
		userReq request.EndpointInfo
	)

	// parse user request
	if err := ctx.BodyParser(&userReq); err != nil {
		log.Printf("parse error:\n\t%v\n", err)

		return errParsingFailed
	}

	// check address
	if userReq.Address == "" {
		return errEmptyAddress
	}

	// creating a new endpoint model
	e := model.Endpoint{
		Username:    ctx.Locals("username").(string),
		Url:         userReq.Address,
		Threshold:   h.Threshold,
		FailedTimes: 0,
		CreateTime:  time.Now(),
	}

	// save endpoint into database
	id, err := h.Repositories.Endpoints.Insert(e)
	if err != nil {
		log.Printf("failed to save into databse:\n\t%v\n", err)

		return errSaveEndpoint
	}

	return ctx.SendString(id)
}

// GetAllEndpoints for a user.
func (h *Handler) GetAllEndpoints(ctx *fiber.Ctx) error {
	// create endpoints list
	var endpoints []response.EndpointResponse

	// get all endpoints of a user
	list := h.Repositories.Endpoints.GetUserEndpoints(ctx.Locals("username").(string))

	// create responses
	for _, item := range list {
		er := response.EndpointResponse{
			Id:        item.ID.Hex(),
			Address:   item.Url,
			CreatedAt: item.CreateTime,
		}

		endpoints = append(endpoints, er)
	}

	return ctx.JSON(endpoints)
}

// GetEndpointStatus will return one endpoint status.
func (h *Handler) GetEndpointStatus(ctx *fiber.Ctx) error {
	// create list of requests
	var requests []response.EndpointRequest

	// get one endpoint requests
	list := h.Repositories.Requests.GetAll(ctx.Locals("id").(string))

	// generate requests response
	for _, item := range list {
		er := response.EndpointRequest{
			Status: item.Code,
			Time:   item.CreateTime,
		}

		requests = append(requests, er)
	}

	return ctx.JSON(requests)
}

// RemoveUserEndpoint will remove and endpoint.
func (h *Handler) RemoveUserEndpoint(ctx *fiber.Ctx) error {
	if err := h.Repositories.Endpoints.Delete(ctx.Locals("id").(string)); err != nil {
		log.Printf("failed to remove endpoint:\n\t%v\n", err)

		return errRemoveEndpoint
	}

	// remove all requests for endpoint
	go func() {
		_ = h.Repositories.Requests.Delete(ctx.Locals("id").(string))
	}()

	return ctx.SendStatus(http.StatusNoContent)
}
