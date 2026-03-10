# IT Inventory Web Application

## Overview
This project is a simple web-based IT inventory management system written in Go.  
It allows users to register, log in, and manage inventory items through a web interface.

The application provides basic inventory management features including authentication, role-based access control, and CRUD operations for inventory items.

---

## Features

### User Management
- User registration
- User login with session-based authentication
- Role-based access (`admin` and `user`)

### Inventory Management
- Create inventory items
- View inventory
- Update inventory items
- Delete inventory items
- Search inventory

### Security
- Password hashing using bcrypt
- Cookie-based sessions
- Basic access control for protected routes

---

## Project Structure

```
.
├── go.mod
├── go.sum
├── main.go
├── models.go
├── handlers.go
├── routes.go
├── middleware.go
├── inventory.json
├── users.json
├── handlers_test.go
├── static/
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   ├── dashboard.html
│   ├── inventory.html
│   └── script.js
```

### File Description

| File | Description |
|-----|-------------|
| `main.go` | Application entry point |
| `routes.go` | Defines all HTTP routes |
| `handlers.go` | Contains API logic and request handlers |
| `middleware.go` | Authentication and authorization middleware |
| `models.go` | Data structures for users and inventory items |
| `inventory.json` | Storage file for inventory data |
| `users.json` | Storage file for user accounts |
| `handlers_test.go` | Basic test cases for handlers |
| `static/` | Frontend HTML and JavaScript files |

---

## Technologies Used

- **Go (Golang)**
- **gorilla/mux** for routing
- **bcrypt** for password hashing
- **HTML / JavaScript** for frontend
- **JSON files** for persistent storage

---

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd IT-Inventory
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the application

```bash
go run main.go
```

The server will start on:

```
http://localhost:8080
```

---

## Default Data Storage

The application stores its data in local JSON files:

- `users.json`
- `inventory.json`

These files act as a simple database.

---

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|------|-----------|-------------|
| POST | `/register` | Register new user |
| POST | `/login` | Login user |
| POST | `/logout` | Logout |

### Inventory

| Method | Endpoint | Description |
|------|-----------|-------------|
| GET | `/inventory` | Get all inventory items |
| POST | `/inventory` | Create new item |
| PUT | `/inventory/{id}` | Update item |
| DELETE | `/inventory/{id}` | Delete item |
| GET | `/search` | Search inventory |

---

## Running Tests

```bash
go test ./...
```

---

## Known Issues / Limitations

- JSON files are used instead of a database
- ID generation may lead to duplicates
- No pagination for inventory lists
- Search endpoint is not protected by authentication
- Concurrency safety is limited when multiple users write simultaneously

---

## Possible Improvements

- Replace JSON storage with a database (PostgreSQL / SQLite)
- Implement JWT authentication
- Add pagination and filtering
- Improve role-based access control
- Add API documentation (Swagger)
- Implement better error handling
- Add frontend framework (React / Vue)

---

## License

This project is intended for educational purposes.