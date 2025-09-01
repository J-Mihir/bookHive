BookHive: A Go Library Management API
This repository contains the backend for BookHive, a complete library management API crafted with Go. It's designed to be both powerful and secure, offering a full suite of features for managing users, cataloging books, and handling complex transactions like borrowing and reservations. The entire system is protected by a modern JWT authentication layer with role-based access for admins and members.

✨ Features
User Management: Full CRUD operations for library members and staff.

JWT Authentication: Secure user registration and login using JSON Web Tokens.

Role-Based Access Control: Differentiates between staff (admin) and student (member) roles, protecting sensitive endpoints.

Book & Inventory Management: Full CRUD for books, including tracking total copies and availability.

Category Management: Organize books by genre or category.

Transaction System:

Borrow and return books.

Automatic due date tracking.

Fine calculation for overdue books.

Borrowing Rules: Enforces business logic, such as a limit on the number of books a member can borrow.

Reservation System: Allows members to reserve books that are currently borrowed.

🛠️ Tech Stack
Language: Go

Router: Gorilla Mux

ORM: GORM

Authentication: golang-jwt/jwt

Password Hashing: bcrypt

Database: MySQL (or any GORM-supported SQL database)

🚀 Getting Started
Prerequisites
Go (version 1.18 or higher)

A running SQL database instance (e.g., MySQL)

Postman for API testing

Installation & Setup
Clone the repository:

git clone <your-repository-url>
cd BookHive

Install dependencies:

go mod tidy

Configure your database:

Open the pkg/config/app.go file.

Update the database connection string with your credentials:

d, err := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

Set Environment Variables:

The application uses a secret key to sign JWTs. Set this in the terminal session where you will run the server.

export JWT_SECRET_KEY="your-long-and-super-secret-string"

Run the server:

go run cmd/main/main.go

The server should now be running on http://localhost:9010.

🧪 API Endpoints & Testing with Postman
Below is a complete guide to testing all available endpoints.

Authentication
1. Register a New User
Method: POST

URL: /register

Body (JSON):

{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "membership_id": "MEMBER001",
    "role": "student"
}

2. User Login
Method: POST

URL: /login

Body (JSON):

{
    "email": "test@example.com",
    "password": "password123"
}

Response:

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Note: Copy this token. You will need it for all protected requests.

Books
Protected routes require a Bearer Token in the Authorization header.

Get All Books (Public)
Method: GET

URL: /books

Get Book by ID (Public)
Method: GET

URL: /books/{bookId}

Create a Book (Admin Only)
Method: POST

URL: /books

Auth: Bearer Token

Body (JSON):

{
    "name": "The Go Programming Language",
    "author": "Alan Donovan",
    "isbn": "9780134190440",
    "copies": 5,
    "category_id": 1
}

Transactions
Borrow a Book
Method: POST

URL: /transactions/borrow

Body (JSON):

{
    "user_id": 1,
    "book_id": 1
}

Return a Book
Method: PUT

URL: /transactions/{transactionId}/return

Reservations
Create a Reservation
Method: POST

URL: /reservations

Auth: Bearer Token

Body (JSON):

{
    "user_id": 2,
    "book_id": 1
}
