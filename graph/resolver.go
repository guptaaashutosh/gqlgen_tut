package graph

import (
	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resolver struct {
	books       []*model.Book
	book        *model.Book
	updatedBook *model.UpdateBook
	DB          *pgxpool.Pool
	users       []*model.User
	user        *model.User
}

