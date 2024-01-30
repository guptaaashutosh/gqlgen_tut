package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"

	loaders "github.com/guptaaashutosh/gqlgen-prac/data_loader"
	dbsetup "github.com/guptaaashutosh/gqlgen-prac/dbSetup"
	"github.com/guptaaashutosh/gqlgen-prac/graph"
	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
	"github.com/guptaaashutosh/gqlgen-prac/middleware"
	"github.com/guptaaashutosh/gqlgen-prac/utils"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	db := dbsetup.ConnectDB()

	router.Use(middleware.UserMiddleware())

	resolver := &graph.Resolver{
		DB: db,
	}

	resolverConfig := graph.Config{Resolvers: resolver}

	// authentication
	resolverConfig.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		authToken := ctx.Value("token")
		if authToken == "" {
			return nil, errors.New("auth token is empty")
		}
		//validate token
		tokenError := utils.VerifyToken(fmt.Sprintf("%v", authToken))
		if tokenError != nil {
			return nil, tokenError
		}
		return next(ctx)
	}
	// authorization
	resolverConfig.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		authRole := ctx.Value("role")
		if authRole == "" {
			return nil, errors.New("role is empty")
		}
		if authRole != "admin" {
			return nil, errors.New("Unauthorized")
		}
		return next(ctx)
	}

	es := graph.NewExecutableSchema(resolverConfig)
	var srv http.Handler = handler.NewDefaultServer(es) // by default NewDefaultServer have transport connection with websocket

	// adding dataloader middleware
	srv = loaders.Middleware(db, srv)

	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
