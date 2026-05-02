# Student Management API

A robust RESTful API built with Go (v1.25.0) for managing students, courses, and their enrollments. It features Role-Based Access Control (RBAC) with regular users and admins, JWT-based authentication, and a PostgreSQL database.

## Features

- **Authentication & Authorization:** JWT-based user registration, login, token refresh, and admin routes.
- **Student Management:** Register, update, delete, and list students.
- **Course Management:** Create, update, delete, and list courses.
- **Enrollment Management:** Enroll students in courses, unenroll, and list course/student associations. 
- **Roles:** Specific actions (like creating courses or deleting students) are restricted to standard `admin` roles, while regular users have restricted access to public or personal data.
- **Database Migrations:** SQL migration management via `goose`.
- **Type-Safe SQL:** Queries and models are auto-generated using `sqlc`.

## Tech Stack

- **Language:** Go (1.25.0), using standard library routing (`net/http` ServeMux).
- **Database:** PostgreSQL 16
- **Database Driver:** `jackc/pgx/v5`
- **SQL Generator:** [sqlc](https://sqlc.dev/)
- **Migrations:** [goose](https://github.com/pressly/goose)
- **Live Reload:** [air](https://github.com/cosmtrek/air)
- **Authentication:** `golang-jwt/jwt/v5` and `golang.org/x/crypto/bcrypt`
- **Infrastructure:** Docker & Docker Compose (for the DB layer)

## Prerequisites

Ensure you have the following installed to run and develop the project:

- [Go](https://golang.org/) (v1.25.0+)
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (usually pre-installed on Unix systems)
- [goose](https://github.com/pressly/goose) (for migrations)
- [sqlc](https://sqlc.dev/) (for generating Go code from SQL)
- [air](https://github.com/cosmtrek/air) (for live reloading during development)

## Getting Started

### 1. Configure Environment Variables

Create a `.env` file in the root directory by copying the following template. Update the placeholder values if necessary.

```env
VERSION=1.0.0
HTTP_PORT=8080
SERVICE_NAME="Student Management System"
JWT_SECRET_KEY="very_strong_secret_key123"

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=student_db
DB_SSL_MODE=false
```

### 2. Start the Database

Boot up the PostgreSQL database via Docker Compose:

```bash
make run
```
*(Note: `make run` also starts `air`. You may want to run `docker compose up -d` separately if you prefer running via `go run main.go` first.)*

### 3. Run Database Migrations

Apply the database schema changes using goose:

```bash
make migrate-up
```

### 4. Generate SQL Code (Optional)

If you modify any SQL queries in `sql/queries/` or schema files in `sql/migrations/`, regenerate the Go models via sqlc:

```bash
make sqlc
```

### 5. Seed Admin User

You can create a default admin user via the seed script to start accessing admin-only endpoints:

```bash
make seed-admin
```

### 6. Run the Application locally

Start the continuous development server (Air) and the database container:
```bash
make run
```

To stop the Docker containers and clear up resources:
```bash
make stop
```

## API Endpoints

### Auth 
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Authenticate a user and receive tokens
- `POST /auth/refresh` - Refresh an active JWT
- `GET /auth/me` - Get current user profile (Protected)
- `POST /auth/logout` - Logout a user (Protected)

### Students
- `POST /students` - Register a student profile (Protected)
- `GET /students/me` - Get logged-in student's profile (Protected)
- `GET /students` - List all students (Protected - Admin only)
- `GET /students/{id}` - Get a student by ID (Protected - Admin only)
- `PATCH /students/{id}` - Update a student (Protected - Admin only)
- `DELETE /students/{id}` - Delete a student (Protected - Admin only)

### Courses
- `GET /courses` - List all courses (Protected)
- `GET /courses/{id}` - Get a specific course (Protected)
- `POST /courses` - Create a course (Protected - Admin only)
- `PATCH /courses/{id}` - Update a course (Protected - Admin only)
- `DELETE /courses/{id}` - Delete a course (Protected - Admin only)

### Enrollments
- `POST /enrollments` - Enroll a student in a course (Protected - Admin only)
- `DELETE /enrollments/{studentID}/{courseID}` - Unenroll a student (Protected)
- `GET /students/{id}/courses` - Check courses for a student (Protected)
- `GET /courses/{id}/students` - List students in a course (Protected - Admin only)
