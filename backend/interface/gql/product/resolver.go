package product_gql

import ucase "backend/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	productUcase ucase.IProductUcase
}
