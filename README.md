# What is Chirpy

Chirpy is a  social media backend API written in Go

## Installation/Usage

Clone the following repository "https://github.com/Wayne-Francis/chirpy"

Run go build -o out && ./out

Set up a .env in the root of your project and add it to your .gitignore file

for the .env file you will require the follwing params

DB_URL="username://postgres:username@localhost:xxx/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="example string"
POLKA_KEY=example number

## Features

- **User Authentication**: Secure password hashing using Argon2 and JWT-based session management.
- **Chirp Management**: Full CRUD operations for "chirps" with profanity filtering.
- **Admin Tools**: Protected metrics and health check endpoints.
- **Webhooks**: Integration for handling upgraded memberships via external services.

Language: Go (Golang)
Database: PostgreSQL
Tools: Goose (for migrations), SQLC (for type-safe SQL)