package graph

import (
	// "github.com/guptaaashutosh/gqlgen-prac/graph/model"

	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	books       []*model.Book
	book        *model.Book
	updatedBook *model.UpdateBook
	DB          *pgxpool.Pool
	users       []*model.User
	user        *model.User
	// TotalBooksMap map[string]model.TotalBooks
	// BooksMap map[string][]model.Book
}












// for pagination
// func NewResolver() graph.Config {
// 	const nBooks = 20
// 	const nBooksPerBook = 100
// 	r := Resolver{}
// 	r.TotalBooksMap = make(map[string]model.TotalBooks, nBooks)
// 	r.BooksMap = make(map[string][]model.Book, nBooksPerBook)

// 	// rand.Seed(time.Now().UnixNano())
// 	for i := 0; i < nBooks; i++ {
// 	  id := strconv.Itoa(i + 1)
// 	  mockChatRoom := model.TotalBooks{
// 		ID: id,
// 		Author: fmt.Sprintf("author %d", i),
// 		Title: fmt.Sprintf("tile %d", i),
// 		PublicationYear: 2000+i,
// 	  }
// 	  r.TotalBooksMap[id] = mockChatRoom
// 	  r.BooksMap[id] = make([]model.Book, nBooksPerBook)

// 	  // Generate messages for the ChatRoom
// 	  for k := 0; k < nBooksPerBook; k++ {
// 		id := strconv.Itoa(k + 1)
// 		text := fmt.Sprintf("Message %d", k)

// 		mockMessage := model.Book{
// 	  ID: id,
// 	  Author: &text,
// 		}

// 		r.BooksMap[id][k] = mockMessage
// 	  }
// 	}

// 	return graph.Config{
// 	  Resolvers: &r,
// 	}
//   }
