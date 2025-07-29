# EDNET FIBER API

A comprehensive EdTech API built with the Go Fiber framework and Ent ORM. This project provides a robust backend solution for an educational platform, handling user management, course administration, content delivery, and more.

![alt text](https://github.com/kayprogrammer/ednet-fiber-api/blob/main/display/fiber.png?raw=true)

## Features

*   **User Authentication & Authorization:** Secure user registration, login, and role-based access control (Admin, Instructor, Student).
*   **Course Management:** Create, read, update, and delete courses, categories, lessons, quizzes, questions, and answers.
*   **Enrollment & Progress Tracking:** Manage student enrollments in courses and track their lesson progress.
*   **Payments:** Handle payment records for course purchases.
*   **Reviews:** Allow users to leave reviews for courses.
*   **Seeding:** Initial data seeding for quick setup and testing.

## Technologies Used

*   **Framework:** Go Fiber (a fast, expressive, and minimalist web framework for Go).
*   **ORM:** Ent (an entity framework for Go, providing a simple yet powerful API for modeling and querying data).
*   **Database:** PostgreSQL (a powerful, open-source object-relational database system).
*   **Dependency Management:** Go Modules.
*   **Hot Reloading:** Air (for Go applications during development).
*   **Containerization:** Docker & Docker Compose.

## Getting Started

Follow these instructions to set up and run the project locally.

### Prerequisites

*   Go (version 1.18 or higher)
*   Docker & Docker Compose
*   PostgreSQL client (optional, for direct database interaction)

### 1. Clone the Repository

```bash
git clone https://github.com/kayprogrammer/ednet-fiber-api.git
cd ednet-fiber-api
```

### 2. Environment Variables

Create a `.env` file in the root directory by copying `.env.example` and fill in the necessary details.

```bash
cp .env.example .env
```

Edit the `.env` file with your database credentials and other configurations.

### 3. Run with Docker Compose (Recommended for Production-like Environment)

This will set up the PostgreSQL database and run the Go application in Docker containers.

```bash
docker-compose up --build
```

The API will be accessible at `http://localhost:8000` (or the port specified in your `.env` file).

### 4. Run Locally with Air (for Development with Hot Reloading)

First, ensure you have `air` installed:

```bash
go install github.com/cosmtrek/air@latest
```

Then, start the PostgreSQL database using Docker Compose (without the Go app):

```bash
docker-compose up db
```

Once the database container is running, you can run the Go application with hot-reloading:

```bash
air
```

The API will be accessible at `http://localhost:8000` (or the port specified in your `.env` file). Any code changes will automatically trigger a rebuild and restart of the application.

### 5. Database Migrations

Ent ORM handles migrations. To generate and run migrations, you typically use Ent's codegen and migrate commands. These are usually integrated into the `go generate ./ent` command and then applied when the application starts.

To manually generate Ent assets and migrations:

```bash
go generate ./ent
```

The application will apply pending migrations on startup if configured to do so.

## API Documentation

API documentation (Swagger/OpenAPI) is available at `/swagger/index.html` when the application is running.

## Live URL - (EDNET API Docs)[https://ednet-api.fly.io]

![alt text](https://github.com/kayprogrammer/ednet-fiber-api/blob/main/display/disp1.png?raw=true)
![alt text](https://github.com/kayprogrammer/ednet-fiber-api/blob/main/display/disp2.png?raw=true)
![alt text](https://github.com/kayprogrammer/ednet-fiber-api/blob/main/display/disp3.png?raw=true)
![alt text](https://github.com/kayprogrammer/ednet-fiber-api/blob/main/display/disp4.png?raw=true)

---