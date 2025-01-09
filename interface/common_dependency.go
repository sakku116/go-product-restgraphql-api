package interface_pkg

import ucase "backend/usecase"

type CommonDependency struct {
	AuthUcase    ucase.IAuthUcase
	UserUcase    ucase.IUserUcase
	ProductUcase ucase.IProductUcase
}
