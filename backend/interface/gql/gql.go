package gql

import (
	interface_pkg "backend/interface"
	auth_gql "backend/interface/gql/auth"
	product_gql "backend/interface/gql/product"
	user_gql "backend/interface/gql/user"

	_ "backend/docs"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/vektah/gqlparser/v2/ast"
)

var logger = logging.MustGetLogger("gql")

func SetupGql(ginEngine *gin.Engine, commonDependencies interface_pkg.CommonDependency) {
	// auth
	setupGraphQLHandler(
		ginEngine,
		"/auth/graphql",
		"/auth/graphql/playground",
		auth_gql.NewExecutableSchema(
			auth_gql.Config{
				Resolvers: auth_gql.NewResolver(commonDependencies.AuthUcase),
			},
		),
		"GraphQL Auth Playground",
	)

	// user
	setupGraphQLHandler(
		ginEngine,
		"/users/graphql",
		"/users/graphql/playground",
		user_gql.NewExecutableSchema(
			user_gql.Config{
				Resolvers: user_gql.NewResolver(commonDependencies.UserUcase),
			},
		),
		"GraphQL User Playground",
	)

	// product
	setupGraphQLHandler(
		ginEngine,
		"/products/graphql",
		"/products/graphql/playground",
		product_gql.NewExecutableSchema(
			product_gql.Config{
				Resolvers: product_gql.NewResolver(commonDependencies.ProductUcase),
			},
		),
		"GraphQL Product Playground",
	)
}

func setupGraphQLHandler(
	ginEngine *gin.Engine,
	endpoint string,
	playgroundEndpoint string,
	schema graphql.ExecutableSchema,
	playgroundTitle string,
) {
	handler := handler.New(schema)
	handler.AddTransport(transport.Options{})
	handler.AddTransport(transport.GET{})
	handler.AddTransport(transport.POST{})
	handler.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	handler.Use(extension.Introspection{})
	handler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	ginEngine.POST(endpoint, func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})
	ginEngine.GET(playgroundEndpoint, func(c *gin.Context) {
		playground.Handler(playgroundTitle, endpoint).ServeHTTP(c.Writer, c.Request)
	})
}
