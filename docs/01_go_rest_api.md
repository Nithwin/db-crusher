# Chapter 1: Teaching Go REST APIs

Welcome to Chapter 1! In this chapter, we will learn what a REST API is and how to build one from absolute scratch using the Go programming language.

## What is a REST API?
Imagine you are sitting in a restaurant. You are the **Client**. The kitchen is the **Server** (where the food/data is made). You cannot just walk into the kitchen and cook your own food. You need a waiter. 

An **API** (Application Programming Interface) is the waiter. You give the API an "order" (a Request), the API takes it to the server, and the API brings your "food" (a Response) back to you. **REST** is just a set of rules on how that waiter should behave (using HTTP methods like GET, POST, DELETE).

## Step 1: Setting Up the File
To start, we need to create our main entry point. In your terminal, type:
```bash
touch main.go
```
Open `main.go`. Every executable Go program *must* start with `package main`. 

We also need to import tools. Go has a massive "standard library" (built-in tools). We will import `fmt` (to print text to the screen) and `net/http` (the waiter!).

```go
package main

import (
    "fmt"
    "net/http"
)
```

## Step 2: Creating a Handler (The Kitchen Staff)
When an order comes in, someone has to process it. In Go, this is called a **Handler Function**. 

Let's write a simple function that responds to a "Health Check" (a way to ask the server, "Are you alive?").

```go
// Every HTTP handler must take these two exact arguments!
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    // 1. We use 'w' to talk BACK to the user. 
    // We tell them we are sending plain text.
    w.Header().Set("Content-Type", "text/plain")
    
    // 2. We write the actual message into the response.
    w.Write([]byte("Server is alive and well!"))
}
```

## Step 3: Routing (The Hostess)
Now we have a kitchen staff member (`HealthHandler`), but how does an order get to them? We need a Router. 

We write our `main` function. This is the heart of the application.

```go
func main() {
    // We tell the Go Router: 
    // "If a GET request comes to /health, give it to the HealthHandler"
    http.HandleFunc("GET /health", HealthHandler)

    // Print a message so we know it started
    fmt.Println("Server is running on http://localhost:8080")

    // Start the server on port 8080. This function runs forever!
    err := http.ListenAndServe(":8080", nil)
    
    if err != nil {
        fmt.Println("Server crashed:", err)
    }
}
```

## Step 4: Running the Code!
You've just written a complete web server! Let's run it.
Open your terminal and type:
```bash
go run main.go
```
You will see `Server is running on http://localhost:8080`. 
Now, open your web browser and type `http://localhost:8080/health`. You will see your message: **"Server is alive and well!"**

In Chapter 2, we will replace this simple text response with real data from a PostgreSQL database!
