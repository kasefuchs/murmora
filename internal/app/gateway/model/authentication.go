// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package model

type AuthenticationLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type AuthenticationRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=32,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type AuthenticationTokenResponse struct {
	Token string `json:"token"`
}
