// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package gateway

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/gateway/handler"
)

var validate = validator.New()

func SetupRoutes(
	router fiber.Router,
	userClient user.UserServiceClient,
	tokenClient token.TokenServiceClient,
	sessionClient session.SessionServiceClient,
	authenticationClient authentication.AuthenticationServiceClient,
) {
	authenticationGroup := router.Group("/authentication")
	{
		authenticationHandler := handler.NewAuthentication(validate, authenticationClient)

		authenticationGroup.Post("/login", authenticationHandler.Login)
		authenticationGroup.Post("/register", authenticationHandler.Register)
	}
}
