# gin-simple-token-middleware 🔑

A purposely simple middleware for securing access in the [Gin framework](https://github.com/gin-gonic/gin) with a statically-issued token. Tokens are provided once at creation and are not meant to be changed. They may be provided via either GET (`?access_token=<token>`) or the [HTTP ‘Authorization’ header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization) by using the ‘Token’ type (`Authorization: Token <token>`).

## Usage

You may secure particular routes or complete groups — depending on where you will fit the middleware handler. The following example illustrates how to use the middleware for securing a group on a default Gin instance:

```go
router := gin.Default()
group := router.Group("/", tokenmiddleware.NewHandler("some_token"))
group.POST("/graphql", graphQLHandler)
```

Access then is granted if requests either provide the respective token via GET (`/graphql?access_token=some_token`) or the HTTP header (`Authentication: Token some_token`). If not, access is denied with an HTTP 401.
