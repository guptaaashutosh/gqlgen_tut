package loaders

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type DbReader struct {
	DB *pgxpool.Pool
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	ChapterLoader     *dataloadgen.Loader[string, *model.Chapter]
	BookLoader        *dataloadgen.Loader[string, *model.Book]
	BookDetailsLoader *dataloadgen.Loader[string, *model.BookDetails]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(db *pgxpool.Pool) *Loaders {
	dl := &DbReader{DB: db}
	return &Loaders{
		ChapterLoader:     dataloadgen.NewLoader(dl.getChapters, dataloadgen.WithWait(time.Millisecond)),
		BookLoader:        dataloadgen.NewLoader[string, *model.Book](dl.getBooks, dataloadgen.WithWait(time.Millisecond)),
		BookDetailsLoader: dataloadgen.NewLoader(dl.getBookDetails, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
// func Middleware(db *pgxpool.Pool, next http.Handler) http.Handler {
// 	// return a middleware that injects the loader to the request context
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		loader := NewLoaders(db)
// 		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
// 		next.ServeHTTP(w, r)
// 	})
// }

func Middleware(db *pgxpool.Pool, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(db)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// for http.server
// func Middleware(db *pgxpool.Pool, next *handler.Server) http.Handler {
// 	// return a middleware that injects the loader to the request context
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		loader := NewLoaders(db)
// 		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
// 		next.ServeHTTP(w, r)
// 	})
// }

// ----- chapter -----
func (c DbReader) getChapters(ctx context.Context, chapterIds []string) ([]*model.Chapter, []error) {

	// Convert bookIds to a slice of empty interfaces
	var idInterfaces []interface{}
	for _, id := range chapterIds {
		idInterfaces = append(idInterfaces, id)
	}

	// Generate placeholders for the IN clause
	placeholders := make([]string, len(chapterIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Construct the SQL query with placeholders
	query := `SELECT cid, pages, duration FROM chapter WHERE bid IN (` +
		strings.Join(placeholders, ",") + `)`

	rows, err := c.DB.Query(ctx, query, idInterfaces...)

	if err != nil {
		return nil, []error{err}
	}
	var chapters []*model.Chapter
	errs := make([]error, 0, len(chapterIds))
	for rows.Next() {
		var chapter model.Chapter
		if err := rows.Scan(&chapter.Cid, &chapter.Pages, &chapter.Duration); err != nil {
			errs = append(errs, err)
			continue
		}
		chapters = append(chapters, &chapter)
	}
	return chapters, errs
}

// ---------------- for book -------------------
func (b DbReader) getBooks(ctx context.Context, bookIds []string) ([]*model.Book, []error) {
	// Convert bookIds to a slice of empty interfaces
	var idInterfaces []interface{}
	for _, id := range bookIds {
		idInterfaces = append(idInterfaces, id)
	}

	// Generate placeholders for the IN clause
	placeholders := make([]string, len(bookIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Construct the SQL query with placeholders
	query := `SELECT id, title, author, publication_year FROM book WHERE id IN (` +
		strings.Join(placeholders, ",") + `)`

	// Query the database
	// rows, err := b.DB.Query(ctx, `SELECT id, title, author, publication_year FROM book WHERE id IN ($1)`, idInterfaces...)
	rows, err := b.DB.Query(ctx, query, idInterfaces...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	// Process the query results
	var books []*model.Book
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear); err != nil {
			return nil, []error{err}
		}
		books = append(books, book)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, []error{err}
	}

	return books, nil
}

// for bookdetails
func (b DbReader) getBookDetails(ctx context.Context, bookIds []string) ([]*model.BookDetails, []error) {
	// Convert bookIds to a slice of empty interfaces
	var idInterfaces []interface{}
	for _, id := range bookIds {
		idInterfaces = append(idInterfaces, id)
	}

	// Generate placeholders for the IN clause
	placeholders := make([]string, len(bookIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Construct the SQL query with placeholders
	query := `SELECT id, title, author, publication_year FROM book WHERE id IN (` +
		strings.Join(placeholders, ",") + `)`

	// Query the database
	rows, err := b.DB.Query(ctx, query, idInterfaces...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	// Process the query results
	var booksDetails []*model.BookDetails
	for rows.Next() {
		bookDetails := &model.BookDetails{}
		if err := rows.Scan(&bookDetails.ID, &bookDetails.Title, &bookDetails.Author, &bookDetails.PublicationYear); err != nil {
			return nil, []error{err}
		}
		booksDetails = append(booksDetails, bookDetails)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, []error{err}
	}

	return booksDetails, nil
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// // GetUser returns single user by id efficiently
// func GetChapter(ctx context.Context, chapterId string) (*model.Chapter, error) {
// 	loaders := For(ctx)
// 	return loaders.ChapterLoader.Load(ctx, chapterId)
// }

// // GetUsers returns many users by ids efficiently
// func GetChapters(ctx context.Context, chapterIds []string) ([]*model.Chapter, error) {
// 	loaders := For(ctx)
// 	return loaders.ChapterLoader.LoadAll(ctx, chapterIds)
// }

// GetUsers returns many users by ids efficiently
// func GetChapters(ctx context.Context, bookIds []string) ([]*model.Chapter, error) {
// 	loaders := For(ctx)
// 	return loaders.ChapterLoader.LoadAll(ctx, bookIds)
// }

// GetUser returns single user by id efficiently
func GetChapter(ctx context.Context, bookId string) (*model.Chapter, error) {
	loaders := For(ctx)
	return loaders.ChapterLoader.Load(ctx, bookId)
}

func GetChapters(ctx context.Context, bookIds []string) ([]*model.Chapter, error) {
	loaders := For(ctx)
	return loaders.ChapterLoader.LoadAll(ctx, bookIds)
}

// GetUser returns single user by id efficiently
func GetBook(ctx context.Context, bookId string) (*model.Book, error) {
	loaders := For(ctx)
	return loaders.BookLoader.Load(ctx, bookId)
}

// GetUsers returns many users by ids efficiently
func GetBooks(ctx context.Context, bookIds []string) ([]*model.Book, error) {
	loaders := For(ctx)
	return loaders.BookLoader.LoadAll(ctx, bookIds)
}

func GetBookDetails(ctx context.Context, bookId string) (*model.BookDetails, error) {
	loaders := For(ctx)
	return loaders.BookDetailsLoader.Load(ctx, bookId)
}

// GetUsers returns many users by ids efficiently
func GetBooksDetails(ctx context.Context, bookIds []string) ([]*model.BookDetails, error) {
	loaders := For(ctx)
	return loaders.BookDetailsLoader.LoadAll(ctx, bookIds)
}

//GetBooksDetailsWithChapters
// func GetBooksDetailsWithChapters(ctx context.Context, bookIds []string) ([]*model.Book, error) {
// 	loaders := For(ctx)
// 	return loaders.BookLoader.LoadAll(ctx, bookIds)
// }
