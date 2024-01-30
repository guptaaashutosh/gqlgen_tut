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

type Loaders struct {
	ChapterLoader     *dataloadgen.Loader[string, *model.Chapter]
	BookLoader        *dataloadgen.Loader[string, *model.Book]
	BookDetailsLoader *dataloadgen.Loader[string, *model.BookDetails]
}

func NewLoaders(db *pgxpool.Pool) *Loaders {
	dl := &DbReader{DB: db}
	return &Loaders{
		ChapterLoader:     dataloadgen.NewLoader(dl.getChapters, dataloadgen.WithWait(time.Millisecond)),
		BookLoader:        dataloadgen.NewLoader[string, *model.Book](dl.getBooks, dataloadgen.WithWait(time.Millisecond)),
		BookDetailsLoader: dataloadgen.NewLoader(dl.getBookDetails, dataloadgen.WithWait(time.Millisecond)),
	}
}

func Middleware(db *pgxpool.Pool, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(db)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

func (c DbReader) getChapters(ctx context.Context, chapterIds []string) ([]*model.Chapter, []error) {
	var idInterfaces []interface{}
	for _, id := range chapterIds {
		idInterfaces = append(idInterfaces, id)
	}

	placeholders := make([]string, len(chapterIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

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

func (b DbReader) getBooks(ctx context.Context, bookIds []string) ([]*model.Book, []error) {

	var idInterfaces []interface{}
	for _, id := range bookIds {
		idInterfaces = append(idInterfaces, id)
	}

	placeholders := make([]string, len(bookIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := `SELECT id, title, author, publication_year FROM book WHERE id IN (` +
		strings.Join(placeholders, ",") + `)`

	rows, err := b.DB.Query(ctx, query, idInterfaces...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()


	var books []*model.Book
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear); err != nil {
			return nil, []error{err}
		}
		books = append(books, book)
	}

	
	if err := rows.Err(); err != nil {
		return nil, []error{err}
	}

	return books, nil
}


func (b DbReader) getBookDetails(ctx context.Context, bookIds []string) ([]*model.BookDetails, []error) {
	var idInterfaces []interface{}
	for _, id := range bookIds {
		idInterfaces = append(idInterfaces, id)
	}

	placeholders := make([]string, len(bookIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := `SELECT id, title, author, publication_year FROM book WHERE id IN (` +
		strings.Join(placeholders, ",") + `)`

	rows, err := b.DB.Query(ctx, query, idInterfaces...)
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	var booksDetails []*model.BookDetails
	for rows.Next() {
		bookDetails := &model.BookDetails{}
		if err := rows.Scan(&bookDetails.ID, &bookDetails.Title, &bookDetails.Author, &bookDetails.PublicationYear); err != nil {
			return nil, []error{err}
		}
		booksDetails = append(booksDetails, bookDetails)
	}

	if err := rows.Err(); err != nil {
		return nil, []error{err}
	}

	return booksDetails, nil
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}



func GetChapter(ctx context.Context, bookId string) (*model.Chapter, error) {
	loaders := For(ctx)
	return loaders.ChapterLoader.Load(ctx, bookId)
}

func GetChapters(ctx context.Context, bookIds []string) ([]*model.Chapter, error) {
	loaders := For(ctx)
	return loaders.ChapterLoader.LoadAll(ctx, bookIds)
}


func GetBook(ctx context.Context, bookId string) (*model.Book, error) {
	loaders := For(ctx)
	return loaders.BookLoader.Load(ctx, bookId)
}


func GetBooks(ctx context.Context, bookIds []string) ([]*model.Book, error) {
	loaders := For(ctx)
	return loaders.BookLoader.LoadAll(ctx, bookIds)
}

func GetBookDetails(ctx context.Context, bookId string) (*model.BookDetails, error) {
	loaders := For(ctx)
	return loaders.BookDetailsLoader.Load(ctx, bookId)
}


func GetBooksDetails(ctx context.Context, bookIds []string) ([]*model.BookDetails, error) {
	loaders := For(ctx)
	return loaders.BookDetailsLoader.LoadAll(ctx, bookIds)
}

