![gqlgen](https://user-images.githubusercontent.com/980499/133180111-d064b38c-6eb9-444b-a60f-7005a6e68222.png)


# gqlgen [![Integration](https://github.com/99designs/gqlgen/actions/workflows/integration.yml/badge.svg)](https://github.com/99designs/gqlgen/actions) [![Coverage Status](https://coveralls.io/repos/github/99designs/gqlgen/badge.svg?branch=master)](https://coveralls.io/github/99designs/gqlgen?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/99designs/gqlgen)](https://goreportcard.com/report/github.com/99designs/gqlgen) [![Go Reference](https://pkg.go.dev/badge/github.com/99designs/gqlgen.svg)](https://pkg.go.dev/github.com/99designs/gqlgen) [![Read the Docs](https://badgen.net/badge/docs/available/green)](http://gqlgen.com/)

## What is gqlgen?

[gqlgen](https://github.com/99designs/gqlgen) is a Go library for building GraphQL servers without any fuss.<br/>

- **gqlgen is based on a Schema first approach** — You get to Define your API using the GraphQL [Schema Definition Language](http://graphql.org/learn/schema/).
- **gqlgen prioritizes Type safety** — You should never see `map[string]interface{}` here.
- **gqlgen enables Codegen** — We generate the boring bits, so you can focus on building your app quickly.

Still not convinced enough to use **gqlgen**? Compare **gqlgen** with other Go graphql [implementations](https://gqlgen.com/feature-comparison/)

## Quick start (gqlgen demo implementation)
1. [Initialise a new go module](https://golang.org/doc/tutorial/create-module)

       mkdir example
       cd example
       go mod init github.com/guptaaashutosh/gqlgen-prac   

2. Initialise gqlgen config and generate models

       go run github.com/99designs/gqlgen init

       go mod tidy

4. Create Schema in schema.graphql file as below example (schema reference : https://graphql.org/learn/schema/)

```graphql
      type Book {
        id: String!
        author: String!
        title: String!
        publication_year: Int!
       }  
```

5. Start the graphql server

       go run server.go

More help to get started:
 - [Getting started tutorial](https://gqlgen.com/getting-started/) - a comprehensive guide to help you get started
 - [Real-world examples](https://github.com/99designs/gqlgen/tree/master/_examples) show how to create GraphQL applications
 - [Reference docs](https://pkg.go.dev/github.com/99designs/gqlgen) for the APIs

When you have nested or recursive schema like this:

```graphql
type BookDetails {
  id: String!
  author: String!
  title: String!
  publication_year: Int!
  chapter: [Chapter]  # @goField(forceResolver: true) is built-in indirectives by gqlgen that allow us to force generate a resolver
}
```

You need to tell gqlgen that it should only fetch chapter if the BookDetails requested it. There are two ways to do this;

- #### Using Custom Models

Write a custom model that omits the chapter field:

```go
type Chapter {
  cid: Int!
  pages: Int!
  duration: Int!
}
```

And reference the model in `gqlgen.yml`:

```yaml
# gqlgen.yml
models:
  User:
    model: github.com/you/pkg/model.Chapter # go import path to the User struct above
```

- #### Using Explicit Resolvers

If you want to Keep using the generated model, mark the field as requiring a resolver explicitly in `gqlgen.yml` like this:

```yaml
# gqlgen.yml
  BookDetails:
    fields:
      chapter:
        resolver: true # force a resolver to be generated
```

After doing either of the above and running generate we will need to provide a resolver for chapter:

```go
func (r *bookDetailsResolver) Chapter(ctx context.Context, obj *model.BookDetails) ([]*model.Chapter, error) {
	 // write here to fetch chapters of given bookId
  return chapters,  nil
}
```

For authentication and Authorization using directive 

```graphql
directive @isAuthenticated on FIELD_DEFINITION

directive @hasRole(role: Role!) on FIELD_DEFINITION

enum Role {
    ADMIN
    USER
}
```

## Other Resources

- https://gqlgen.com/
- https://github.com/99designs/gqlgengqlgen-a-graphql-server-generator-for-go/)

copy reference : [link](https://github.com/99designs/gqlgen?tab=readme-ov-file)
