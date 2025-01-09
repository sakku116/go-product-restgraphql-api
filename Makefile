gqlgen:
	@echo "running gqlgen to generate code..."
	@cd interface/gql/auth && gqlgen generate
	@cd interface/gql/product && gqlgen generate
	@cd interface/gql/user && gqlgen generate

genswagger:
	@echo "running swag init to generate swagger files..."
	@swag init