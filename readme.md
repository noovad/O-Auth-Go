## 📚 **OAuth Project with Gin, GORM, PostgreSQL**

---

# **OAuth Project**

This is my learning project about **OAuth2** authentication using **Google OAuth2** with **Gin Framework** , **GORM** as the ORM, and **PostgreSQL** as the database. The project demonstrates login using Google accounts and utilizes the **Google OAuth2 package for Go**.

[Google OAuth2 Go Package](https://pkg.go.dev/golang.org/x/oauth2/google#pkg-overview

---

## 📦 **Tech Stack**

- **Go (Golang)**
- **Gin** - Web framework for Go.
- **OAuth2** - Authentication via Google.
- **GORM** - ORM library for Golang.
- **PostgreSQL** - Relational database.
- **OAuth2** - Authentication via Google.
- **Air** - Hot reload for development.
- **Wire** - Dependency injection for managing dependencies.

---

## ⚙️ **Prerequisites**

Make sure you have installed:

- **Go**: [Download Go](https://go.dev/dl/)
- **PostgreSQL**: [Download PostgreSQL](https://www.postgresql.org/download/)
- **Air** for hot reload:
  ```bash
  go install github.com/cosmtrek/air@latest
  ```
- **Wire** for dependency injection:
  ```bash
  go install github.com/google/wire/cmd/wire@latest
  ```
- **Git**: [Install Git](https://git-scm.com/downloads)

---

## 🚀 **Project Installation**

1. **Clone the Repository**
2. **Environment Configuration**
   Copy the example.env file as .env and update it with your own configuration.
3. **Set up Database**
   Make sure PostgreSQL is running and the database is created. Adjust the database configuration in .env accordingly.
4. **Install Dependencies**
   ```bash
   make install
   ```

---

## 🌐 **Google OAuth2 Configuration**

Follow the instructions to set up your Google OAuth2 credentials:
[Google OAuth2 Guide](https://developers.google.com/identity/protocols/oauth2)

Update the following variables in your **`.env`** file:

```env
GOOGLE_CLIENT_ID='YOUR_GOOGLE_CLIENT_ID'
GOOGLE_CLIENT_SECRET='YOUR_GOOGLE_CLIENT_SECRET'
```

---

## 🔧 **Dependency Injection with Wire**

The project uses **Google Wire** for dependency injection to simplify and manage the dependencies between components.

### How to use Wire:

- **Modify dependencies** in `injector.go` inside the **api** folder.
- To regenerate the dependency graph, use the command:

```bash
make gen
```

This command runs `wire gen ./api` to generate the `wire_gen.go` file.

If you add new services or repositories that require dependency injection, update the injector and regenerate the Wire code.

---

## 📂 **Project Structure**

```plaintext
LEARN_O_AUTH/
│── api/
│   ├── controller/
│   │   └── auth_controller.go      # Controller for OAuth authentication
│   ├── repository/
│   │   └── user_repository.go      # Repository for user data operations
│   ├── service/
│   │   ├── auth_service.go         # Business logic for authentication
│   │   └── user_service.go         # Business logic for user management
│   ├── injector.go                 # Dependency injection using Wire
│   └── wire_gen.go                 # Generated file from Wire
│
│── cmd/
│   └── main.go                     # Application entry point
│
│── config/
│   ├── database.go                 # PostgreSQL database configuration
│   ├── oauth.go                    # Google OAuth2 configuration
│   └── validator.go                # Data validation configuration
│
│── data/
│   ├── google_oauth.go             # Data model for Google OAuth response
│   └── users.go                    # Data structure for user information
│
│── helper/
│   ├── error.go                    # Custom error handling
│   ├── generate_state.go           # State generator for OAuth
│   ├── jwt.go                      # JWT (JSON Web Token) utilities
│   └── response.go                 # Helper for JSON response formatting
│
│── model/
│   ├── migrate.go                  # Database migration file
│   └── user_model.go               # ORM model for user data (GORM)
│
│── presentation/
│   ├── home_page.go                # HTML page for the home view
│   └── login_page.go               # HTML page for the login view
│
│── router/
│   ├── oauth_route.go              # OAuth authentication routes
│   └── setup_router.go             # Main router setup for the application
│
│── tmp/                            # Temporary directory for hot reload (Air)
│── .air.toml                       # Air configuration for hot reload
│── .env                            # Environment variables
│── .gitignore                      # Git ignore file
│── example.env                     # Example configuration for environment variables
│── go.mod                          # Go module definition
│── go.sum                          # Dependency checksum file
│── Makefile                        # Makefile for automation tasks
│── README.md                       # Project documentation
```

---

## 💻 **Running the Project**

Use the available **Makefile** commands:

- **Install dependencies:**

  ```bash
  make install
  ```
- **Generate dependency injection with Wire:**

  ```bash
  make gen
  ```
- **Run with Air (Hot Reload):**

  ```bash
  make run
  ```

---