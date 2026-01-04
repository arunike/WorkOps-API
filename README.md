# 👔 WorkOps - HR Management System API

The powerful backend powering the WorkOps system. A clear, performant, and RESTful API offering comprehensive employee management capabilities.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white) ![Chi](https://img.shields.io/badge/Chi-00ADD8?style=for-the-badge&logo=go&logoColor=white) ![MySQL](https://img.shields.io/badge/MySQL-00000F?style=for-the-badge&logo=mysql&logoColor=white) ![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

## 📋 Table of Contents
- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
- [Project Structure](#-project-structure)
- [API Endpoints](#-api-endpoints)
- [Authors](#-authors)

## 🌟 Overview
The WorkOps API handles all business logic, data persistence, and authentication for the platform. It provides a secure and scalable foundation for managing associates, departments, tasks, and social interactions.

**Key Capabilities:**
- 🔐 **Robust Authentication**: Secure session management and JWT handling.
- 👥 **Associate Management**: CRUD operations for employee data.
- ⚙️ **Process Automation**: Workflows for tasks, approvals, and time-off requests.
- ❤️ **Social Graph**: Managing "Give Thanks" posts, likes, and comments.
- 📊 **Analytics Data**: Serving aggregated data for frontend dashboards.

## ✨ Features

### 👥 Associate & Hierarchy
- **Associates**: Create, read, update, delete functionality for employee records.
- **Offices & Departments**: Management of organizational structure entities.
- **Hierarchy Awareness**: Logic to understand reporting lines and team structures.

### ⚙️ Operational Workflows
- **Task Management**: Endpoints for creating and processing approval tasks (e.g., salary increases).
- **Time Off**: Managing vacation and leave requests with approval status flow.
- **Document Categories**: Taxonomy management for employee documents.

### ❤️ Social & Recognition
- **Thanks Feed**: CRUD for recognition posts.
- **Interactions**: Toggling likes and managing threaded comments on posts.

### 🛡️ Admin & System
- **Menu Permissions**: Granular control over UI element visibility.
- **System Settings**: Global configuration endpoints (e.g., default password).

## 🛠 Tech Stack

### Core Runtime
- **Go (Golang) 1.24** - Performant, statically typed system language.

### Frameworks & Libraries
- **Chi Router** - Lightweight, idiomatic, and composable router config.
- **MySQL Driver** - Robust database connectivity.
- **Standard Lib** - Heavy use of Go's powerful standard library (`net/http`, `database/sql`).

### Data Persistence
- **MySQL** - Relational database for structured data storage.

## 🚀 Getting Started

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- MySQL instance (local or containerized)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ReactJS-HRCore/api
   ```

2. **Environment Variables**
   The application expects a database connection string. By default (in Docker), this is handled via `DSN`.

3. **Run Locally** (Requires running MySQL)
   ```bash
   go run cmd/api/main.go
   ```
   *The server will start on port `8081`.*

### Docker Deployment
The recommended way to run the API is via Docker Compose from the project root:
```bash
# From project root
docker compose up --build -d api
```

## 📁 Project Structure

```
api/
├── cmd/
│   └── api/                # Application entrypoint
│       ├── main.go         # Server initialization
│       ├── routes.go       # Route definitions
│       ├── handlers.go     # Request handlers (Controllers)
│       └── middleware.go   # HTTP middleware (CORS, Auth)
├── internal/
│   └── data/               # Data access layer (Models)
│       ├── models.go       # Model definitions
│       └── *.go            # Database operations
├── db/                     # Database schemas and migrations
│   └── init.sql            # Initial seeding script
├── go.mod                  # Module definition
└── Dockerfile              # Container definition
```

## 🔌 API Endpoints

### Auth
- `POST /login` - Authenticate user
- `POST /register` - Register a new account

### Associates
- `GET /associates` - List all associates
- `POST /associates` - Create new associate
- `GET /associates/{id}` - Get specific associate
- `PUT /associates/{id}` - Update associate details
- `PUT /associates/{id}/password` - Change password

### Social
- `GET /thanks` - Get recognition feed
- `POST /thanks` - Create recognition post
- `POST /thanks/{id}/like` - Like a post
- `POST /thanks/{id}/comment` - Comment on a post

*...and many more for Tasks, Time Off, Offices, etc.*

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.txt) file for details.

## 👤 Author

**Richie Zhou**

- GitHub: [@arunike](https://github.com/arunike)
