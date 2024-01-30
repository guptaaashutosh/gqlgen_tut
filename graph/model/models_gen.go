// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Book struct {
	ID              string `json:"id"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationYear int    `json:"publication_year"`
}

type BookConnection struct {
	Edges    []*BookEdges `json:"edges"`
	PageInfo *PageInfo    `json:"pageInfo"`
}

type BookDetails struct {
	ID              string     `json:"id"`
	Author          string     `json:"author"`
	Title           string     `json:"title"`
	PublicationYear int        `json:"publication_year"`
	Chapter         []*Chapter `json:"chapter,omitempty"`
}

type BookEdges struct {
	Cursor string `json:"cursor"`
	Node   *Book  `json:"node,omitempty"`
}

type Chapter struct {
	Cid      int `json:"cid"`
	Pages    int `json:"pages"`
	Duration int `json:"duration"`
}

type DeletedBook struct {
	ID string `json:"id"`
}

type LoginDetails struct {
	IsLoggedIn *bool  `json:"isLoggedIn,omitempty"`
	Token      string `json:"token"`
}

type Mutation struct {
}

type PageInfo struct {
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
	HasNextPage *bool  `json:"hasNextPage,omitempty"`
}

type Query struct {
}

type Subscription struct {
}

type Time struct {
	UnixTime  int    `json:"unixTime"`
	TimeStamp string `json:"timeStamp"`
}

type TotalBooks struct {
	ID              string          `json:"id"`
	Author          string          `json:"author"`
	Title           string          `json:"title"`
	PublicationYear int             `json:"publication_year"`
	BookConnection  *BookConnection `json:"bookConnection,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewBook struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publicationYear"`
}

type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateBook struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publicationYear"`
}

type UpdatedBook struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publicationYear"`
}

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

var AllRole = []Role{
	RoleAdmin,
	RoleUser,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}