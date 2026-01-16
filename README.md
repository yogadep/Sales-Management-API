# Sales Management API

A simple **Point of Sale (POS) / Sales Management REST API** built using **Golang** and the **Echo Framework**.  
This project was developed as part of a Golang Developer technical assessment.

---

## 🚀 Features

- Authentication & Authorization (JWT)
- Role-based Access Control (Admin & Cashier)
- User Management (Admin only)
- Product Management (CRUD)
- Sales Transaction
  - One transaction can contain multiple products
  - Stock is reduced atomically
- Sales Reports
  - JSON
  - PDF export
  - Excel export
- PostgreSQL database (Supabase)
- Security best practices applied

---

## 👥 Roles & Permissions

| Feature | Admin | Cashier |
|--------|-------|---------|
| Login | ✅ | ✅ |
| Register User | ✅ | |
| List Users | ✅ | |
| Product Management | ✅ | |
| Create Sales Transaction | ✅ | ✅ |
| View Sales | ✅ | ✅ |
| Sales Report (JSON / PDF / Excel) | ✅ | ✅ |

---

## 🧰 Tech Stack

- **Language**: Go
- **Framework**: Echo
- **ORM**: GORM
- **Database**: PostgreSQL (Supabase)
- **Authentication**: JWT
- **Password Hashing**: bcrypt
- **Report Export**:
  - PDF: gofpdf
  - Excel: excelize

---

## 🔐 Security Implementation

This application applies multiple security measures:

- **SQL Injection**  
  Prevented using GORM parameterized queries (no raw SQL concatenation).
- **Authentication & Authorization**  
  JWT-based authentication using the `Authorization: Bearer` header with role-based access control.
- **XSS & Clickjacking Protection**  
  Secure HTTP headers are applied globally.
- **CSRF**  
  Not applicable, as authentication does not rely on cookies but uses JWT headers.
- **Rate Limiting**  
  Global rate limiting is enabled to reduce brute-force and abusive requests.
- **Password Security**  
  Passwords are hashed using bcrypt and never returned by the API.
- **CORS**  
  Explicitly configured allowed origins for frontend access.

---

## ⚙️ Environment Variables

Create a `.env` file in the project root:

APP_PORT=8080  
APP_ENV=development  

JWT_SECRET=your-jwt-secret  

DB_HOST=your-supabase-host  
DB_PORT=5432  
DB_NAME=postgres  
DB_USER=postgres.xxxxx  
DB_PASSWORD=your-db-password  
DB_SSLMODE=require  

---

## ▶️ How to Run

1. Make sure **Go** is installed (Go 1.21+ recommended)
2. Create a `.env` file (see Environment Variables section)
3. Install dependencies:
4. Run the application:

go run ./cmd/api

The server will start at:  
http://localhost:8080

---

## 📮 API Documentation (Postman)

The API documentation is publicly available via Postman:

https://web.postman.co/documentation/17852367-82112a18-bbb1-43c3-917c-70e12745f85a/publish

The documentation provides:
- List of all available endpoints
- Request & response examples
- Authorization details (JWT Bearer Token)

All requests use environment variables:
- `BASE_URL`
- `TOKEN`

---

## 📝 Notes

- All sales transactions are handled atomically using database transactions.
- Product stock is locked during sales creation to prevent race conditions.
- Password hashes are never exposed via API responses.
- This project focuses on backend API design and security best practices for assessment purposes.

---

## 📄 License

This project was created for technical assessment purposes.

