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
| Register User | ✅ |  
| List Users | ✅ |  |
| Product Management | ✅ |  
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
  - Prevented using GORM parameterized queries (no raw SQL concatenation).
- **Authentication & Authorization**
  - JWT-based authentication using the `Authorization: Bearer` header.
  - Role-based access control (Admin and Cashier).
- **XSS & Clickjacking Protection**
  - Secure HTTP headers are applied globally.
- **CSRF**
  - Not applicable, as authentication does not rely on cookies but uses JWT headers.
- **Rate Limiting**
  - Global rate limiting is enabled to reduce brute-force and abusive requests.
- **Password Security**
  - Passwords are hashed using bcrypt and never returned by the API.
- **CORS**
  - Explicitly configured allowed origins for frontend access.

---

## ⚙️ Environment Variables

Create a `.env` file in the project root:

```env
APP_PORT=8080
APP_ENV=development

JWT_SECRET=your-jwt-secret

DB_HOST=your-supabase-host
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres.xxxxx
DB_PASSWORD=your-db-password
DB_SSLMODE=require
```

---

## 📬 Postman Documentation

https://web.postman.co/documentation/17852367-82112a18-bbb1-43c3-917c-70e12745f85a/publish?workspaceId=3d6c77c2-a920-4a68-8ae0-d53e5d27efea&authFlowId=a22ae019-4dca-49d6-a7cc-2cc864c2911c
