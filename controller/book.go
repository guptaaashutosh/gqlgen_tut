package controller

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	loaders "github.com/guptaaashutosh/gqlgen-prac/data_loader"
	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetBooks(ctx context.Context, DB *pgxpool.Pool) ([]*model.Book, error) {
	var books []*model.Book
	rows, err := DB.Query(ctx, "select id, title, author, publication_year from book")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		var book model.Book
		rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear)
		books = append(books, &book)
	}
	return books, nil
}

func GetBooksUsingLimitOffset(ctx context.Context, DB *pgxpool.Pool, limit int, offset int) ([]*model.Book, error)  {
	var books []*model.Book
	rows, err := DB.Query(ctx, "select id, title, author, publication_year from book limit $1 offset $2",limit,offset)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		var book model.Book
		rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear)
		books = append(books, &book)
	}
	return books, nil
}

func GetBook(ctx context.Context, id string) (*model.Book, error) {
	return loaders.GetBook(ctx,id)
}

func UpdateBook(ctx context.Context, DB *pgxpool.Pool, input model.UpdateBook) (*model.UpdatedBook, error) {
	book := &model.Book{
		Title:           input.Title,
		Author:          input.Author,
		ID:              input.ID,
		PublicationYear: input.PublicationYear,
	}
	var updatedBook model.UpdatedBook
	rows := DB.QueryRow(ctx, "UPDATE book set title=$2,author=$3,publication_year=$4 WHERE id=$1 RETURNING id,title,author,publication_year", book.ID, book.Title, book.Author, book.PublicationYear)
	rows.Scan(&updatedBook.ID, &updatedBook.Title, &updatedBook.Author, &updatedBook.PublicationYear)
	return &updatedBook, nil
}

func DeleteBook(ctx context.Context, DB *pgxpool.Pool, id string) (*model.DeletedBook, error) {
	var deletedId model.DeletedBook
	rows, err := DB.Query(ctx, "DELETE FROM book where id=$1 RETURNING id", id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&deletedId.ID)
	}
	return &deletedId, nil
}

func CreateBook(ctx context.Context, DB *pgxpool.Pool, input model.NewBook) (*model.Book, error) {
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
	book := &model.Book{
		Title:           input.Title,
		Author:          input.Author,
		ID:              fmt.Sprintf("T%d", randNumber),
		PublicationYear: input.PublicationYear,
	}
	_, err := DB.Exec(ctx, "INSERT INTO book (id,title,author,publication_year) VALUES ($1,$2,$3,$4)", book.ID, book.Title, book.Author, book.PublicationYear)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return book, nil
}



func BookWithChapter(ctx context.Context, DB *pgxpool.Pool, id string) (*model.BookDetails, error) {
	return loaders.GetBookDetails(ctx,id)
}

//BooksWithChapters
func BooksWithChapters(ctx context.Context, bookIds []string) ([]*model.BookDetails, error) {
	return loaders.GetBooksDetails(ctx,bookIds)
}


func Chapter(ctx context.Context, DB *pgxpool.Pool, obj *model.BookDetails) ([]*model.Chapter, error) {
	rows, err := DB.Query(ctx, "select cid, pages, duration from chapter where bid=$1", obj.ID)
	if err != nil {
		return nil, err
	}
	var chapters []*model.Chapter
	for rows.Next() {
		var chapter model.Chapter
		rows.Scan(&chapter.Cid, &chapter.Pages, &chapter.Duration)
		chapters = append(chapters, &chapter)
	}
	return chapters, nil
	// return loaders.GetChapters()
}

func GetBooksWithIds(ctx context.Context, bookIds []string) ([]*model.Book, error) {
	return loaders.GetBooks(ctx,bookIds)
}

