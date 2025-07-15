// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/internal/app/gateway/model"
)

type Authentication interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type authenticationHandler struct {
	validate             *validator.Validate
	authenticationClient authentication.AuthenticationServiceClient
}

var _ Authentication = (*authenticationHandler)(nil)

func NewAuthentication(
	validate *validator.Validate,
	authenticationClient authentication.AuthenticationServiceClient,
) Authentication {
	return &authenticationHandler{
		validate:             validate,
		authenticationClient: authenticationClient,
	}
}

func (h *authenticationHandler) Login(ctx *fiber.Ctx) error {
	request := new(model.AuthenticationLoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tokenData, err := h.authenticationClient.Login(ctx.Context(), &authentication.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&model.AuthenticationTokenResponse{Token: tokenData.Token.Value})
}

func (h *authenticationHandler) Register(ctx *fiber.Ctx) error {
	request := new(model.AuthenticationRegisterRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tokenData, err := h.authenticationClient.Register(ctx.Context(), &authentication.RegisterRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(&model.AuthenticationTokenResponse{Token: tokenData.Token.Value})
}
