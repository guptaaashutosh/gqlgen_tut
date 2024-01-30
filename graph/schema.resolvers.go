package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/guptaaashutosh/gqlgen-prac/controller"
	"github.com/guptaaashutosh/gqlgen-prac/graph/model"
)
func (r *bookDetailsResolver) Chapter(ctx context.Context, obj *model.BookDetails) ([]*model.Chapter, error) {
	return controller.Chapter(ctx, r.DB, obj)
}
func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	return controller.CreateBook(ctx, r.DB, input)
}
func (r *mutationResolver) DeleteBook(ctx context.Context, id string) (*model.DeletedBook, error) {
	return controller.DeleteBook(ctx, r.DB, id)
}
func (r *mutationResolver) UpdateBook(ctx context.Context, input model.UpdateBook) (*model.UpdatedBook, error) {
	return controller.UpdateBook(ctx, r.DB, input)
}
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {	return controller.CreateNewUser(ctx, r.DB, input)
}

func (r *mutationResolver) LoginUser(ctx context.Context, email string, password string) (*model.LoginDetails, error) {
	return controller.LoginUserController(ctx, r.DB, email, password)
}

func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	return controller.GetBook(ctx, id)
}

func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	return controller.GetBooks(ctx, r.DB)
}

func (r *queryResolver) GetBooksWithLimitOffset(ctx context.Context, limit int, offset int) ([]*model.Book, error) {
	return controller.GetBooksUsingLimitOffset(ctx, r.DB, limit, offset)
}

func (r *queryResolver) BookWithChapter(ctx context.Context, id string) (*model.BookDetails, error) {
	return controller.BookWithChapter(ctx, r.DB, id)
}

func (r *queryResolver) BooksWithChapters(ctx context.Context, ids []string) ([]*model.BookDetails, error) {
	return controller.BooksWithChapters(ctx, ids)
}

func (r *queryResolver) GetBooksWithIds(ctx context.Context, ids []string) ([]*model.Book, error) {
	return controller.GetBooksWithIds(ctx, ids)
}

func (r *queryResolver) GetCurrentTime(ctx context.Context) (*model.Subscription, error) {
	return r.GetCurrentTime(ctx)
}

func (r *queryResolver) TotalBooks(ctx context.Context, id string) (*model.TotalBooks, error) {
	panic(fmt.Errorf("not implemented: TotalBooks - totalBooks"))
}

func (r *queryResolver) GetUser(ctx context.Context, id int) (*model.User, error) {
	return controller.GetUserWithUserId(ctx, r.DB, id)
}

func (r *queryResolver) GetUsers(ctx context.Context) ([]*model.User, error) {
	return controller.GetUsers(ctx, r.DB)
}

func (r *subscriptionResolver) CurrentTime(ctx context.Context) (<-chan *model.Time, error) {
	ch := make(chan *model.Time)
	go func() {
		defer close(ch)
		for {
			time.Sleep(1 * time.Second)
			currentTime := time.Now()
			t := &model.Time{
				UnixTime:  int(currentTime.Unix()),
				TimeStamp: currentTime.Format(time.RFC3339),
			}
			
			select {
			case <-ctx.Done():
				fmt.Println("subscription closed.")
				return
			case ch <- t:
			}
		}
	}()
	return ch, nil
}

func (r *totalBooksResolver) BookConnection(ctx context.Context, obj *model.TotalBooks, first *int, after *string) (*model.BookConnection, error) {
	panic(fmt.Errorf("not implemented: BookConnection - bookConnection"))
}

func (r *Resolver) BookDetails() BookDetailsResolver { return &bookDetailsResolver{r} }

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

func (r *Resolver) TotalBooks() TotalBooksResolver { return &totalBooksResolver{r} }

type bookDetailsResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type totalBooksResolver struct{ *Resolver }
