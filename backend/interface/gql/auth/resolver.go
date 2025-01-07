package auth_gql

import ucase "backend/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authUcase ucase.IAuthUcase
}

func NewResolver(authUcase ucase.IAuthUcase) *Resolver {
	return &Resolver{
		authUcase: authUcase,
	}
}
