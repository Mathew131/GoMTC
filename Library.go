package main

import (
    "fmt"
    "math/rand"
    "strconv"
    "time"
)

type Book struct {
    ID    string
    Title string
    Author string
}

type Searcher interface {
    SearchByID(id string) (Book, bool)
}

type Library interface {
    SearchByName(name string) (Book, bool)
}

type LibraryS struct {
    books []Book
    booksMap map[string]Book
}

func NewLibrary() *LibraryS {
    return &LibraryS{
        books: make([]Book, 0),
        booksMap: make(map[string]Book),
    }
}

func generateID() string {
    rand.Seed(time.Now().UnixNano())
    return strconv.Itoa(rand.Intn(1000000))
}

func (lib *LibraryS) AddBook(title, author string) {
    id := generateID()
    book := Book{
        ID:    id,
        Title: title,
        Author: author,
    }
    lib.books = append(lib.books, book)
    lib.booksMap[id] = book
}

func (lib *LibraryS) SearchByID(id string) (Book, bool) {
    book, found := lib.booksMap[id]
    return book, found
}

func (lib *LibraryS) SearchByName(name string) (Book, bool) {
    for _, book := range lib.books {
        if book.Title == name {
            return book, true
        }
    }
    return Book{}, false
}

func main() {

    library := NewLibrary()

    library.AddBook("Война и мир", "Лев Толстой")
    library.AddBook("Преступление и наказание", "Фёдор Достоевский")
    library.AddBook("Мастер и Маргарита", "Михаил Булгаков")

    book, found := library.SearchByName("Мастер и Маргарита")
    if found {
        fmt.Printf("Найдена книга: %s от автора %s (ID: %s)\n", book.Title, book.Author, book.ID)
    } else {
        fmt.Println("Книга не найдена")
    }

    bookID := book.ID
    foundBook, found := library.SearchByID(bookID)
    if found {
        fmt.Printf("Найдена книга по ID: %s — %s от автора %s\n", foundBook.ID, foundBook.Title, foundBook.Author)
    } else {
        fmt.Println("Книга по ID не найдена")
    }
}
