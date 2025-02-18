package auth_gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.62

import (
	"backend/domain/dto"
	"context"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input RegisterInput) (*RegisterResult, error) {
	// convert payload
	dto := dto.RegisterUserReq{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	}

	// process
	data, err := r.authUcase.Register(ctx, dto)
	if err != nil {
		return nil, err
	}

	// convert resp
	return &RegisterResult{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	// convert payload
	dto := dto.LoginReq{
		Username: input.Username,
		Password: input.Password,
	}

	// process
	data, err := r.authUcase.Login(ctx, dto)
	if err != nil {
		return nil, err
	}

	// convert resp
	return &LoginResult{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenResult, error) {
	// convert payload
	dto := dto.RefreshTokenReq{
		RefreshToken: input.RefreshToken,
	}

	// process
	data, err := r.authUcase.RefreshToken(ctx, dto)
	if err != nil {
		return nil, err
	}

	// convert resp
	return &RefreshTokenResult{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}, nil
}

// CheckToken is the resolver for the checkToken field.
func (r *queryResolver) CheckToken(ctx context.Context, input CheckTokenInput) (*CheckTokenResult, error) {
	// convert payload
	dto := dto.CheckTokenReq{
		AccessToken: input.AccessToken,
	}

	// process
	data, err := r.authUcase.CheckToken(dto)
	if err != nil {
		return nil, err
	}

	// convert resp
	return &CheckTokenResult{
		UUID:     data.UUID,
		Username: data.Username,
		Role:     data.Role.String(),
		Email:    data.Email,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
