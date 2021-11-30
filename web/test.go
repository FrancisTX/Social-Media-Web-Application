package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

// Book ...
type Book struct {
    Title  string
    Author string
}

func main() {
    r := gin.Default()
    r.LoadHTMLFiles("web/src/html/test.html")

    books := make([]Book, 0)
    books = append(books, Book{
        Title:  "Title 1",
        Author: "Author 1",
    })
    books = append(books, Book{
        Title:  "Title 2",
        Author: "Author 2",
    })

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "test.html", gin.H{
            "books": books,
        })
    })
    log.Fatal(r.Run())
}