# BookHive: A Go Library Management API

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)

A complete library management API crafted with Go. Designed to be both powerful and secure, offering a full suite of features for managing users, cataloging books, and handling complex transactions like borrowing and reservations. The entire system is protected by a modern JWT authentication layer with role-based access for admins and members.

✨ Features

- 🔐 **User Management**: Full CRUD operations for library members and staff
- 🎟️ **JWT Authentication**: Secure user registration and login using JSON Web Tokens
- 👥 **Role-Based Access Control**: Differentiates between `staff` (admin) and `student` (member) roles, protecting sensitive endpoints
- 📚 **Book & Inventory Management**: Full CRUD for books, including tracking total copies and availability
- 🏷️ **Category Management**: Organize books by genre or category
- 🔄 **Transaction System**:
  - Borrow and return books
  - Automatic due date tracking
  - Fine calculation for overdue books
- 📖 **Borrowing Rules**: Enforces business logic, such as a limit on the number of books a member can borrow
- 📅 **Reservation System**: Allows members to reserve books that are currently borrowed

🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go** | Primary language |
| **Gorilla Mux** | HTTP router |
| **GORM** | ORM for database operations |
| **golang-jwt/jwt** | JWT authentication |
| **bcrypt** | Password hashing |
| **MySQL** | Database (supports any GORM-compatible SQL database) |

🚀 Getting Started

### Prerequisites

- Go (version 1.18 or higher)
- MySQL or any GORM-supported SQL database
- Postman (for API testing)

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <your-repository-url>
   cd BookHive
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure your database**
   
   Open `pkg/config/app.go` and update the database connection string:
   ```go
   d, err := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
   ```

4. **Set environment variables**
   ```bash
   export JWT_SECRET_KEY="your-long-and-super-secret-string"
   ```

5. **Run the server**
   ```bash
   go run cmd/main/main.go
   ```

The server will be running on `http://localhost:9010`

🧪 API Endpoints & Testing

### Authentication

#### Register a New User
```http
POST /register
Content-Type: application/json

{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "membership_id": "MEMBER001",
    "role": "student"
}
```

#### User Login
```http
POST /login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "password123"
}
```

**Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

> **Note:** Copy this token for use in protected requests as `Authorization: Bearer <token>`

### Books

> 🔒 Protected routes require Bearer Token in Authorization header

#### Get All Books (Public)
```http
GET /books
```

#### Get Book by ID (Public)
```http
GET /books/{bookId}
```

#### Create a Book (Admin Only)
```http
POST /books
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "The Go Programming Language",
    "author": "Alan Donovan",
    "isbn": "9780134190440",
    "copies": 5,
    "category_id": 1
}
```

### Transactions

#### Borrow a Book
```http
POST /transactions/borrow
Content-Type: application/json

{
    "user_id": 1,
    "book_id": 1
}
```

#### Return a Book
```http
PUT /transactions/{transactionId}/return
```

### Reservations

#### Create a Reservation
```http
POST /reservations
Authorization: Bearer <token>
Content-Type: application/json

{
    "user_id": 2,
    "book_id": 1
}
```


*Built with ❤️ using Go • Open for contributions*
