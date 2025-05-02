# 🔐 Go Auth Demo

[![Go Version](https://img.shields.io/badge/go-1.24+-brightgreen)](https://golang.org/dl/)
[![Dockerized](https://img.shields.io/badge/docker-ready-blue)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)

A simple and secure user authentication service written in **Go**, using **password hashing**, **MongoDB** for storage, and **Docker** for easy setup.

> 🔗 GitHub: [github.com/frtasoniero/pwh-auth](https://github.com/frtasoniero/pwh-auth)


---

## 🚀 Features

- ✅ Password-based authentication with bcrypt
- ✅ User registration and login endpoints
- ✅ JSON API responses
- ✅ Swagger endpoint documentation
- ✅ MongoDB persistence
- ✅ Docker + Makefile integration
- ✅ `.env`-based configuration

---

## 🗂 Project Structure

```
auth-demo/
├── cmd/ # Application entrypoint
├── internal/
│ ├── auth/ # Password hashing, user logic
│ └── db/ # MongoDB connection
├── Dockerfile
├── docker-compose.yml
├── .env
├── Makefile
└── README.md
```

---

## ⚙️ Requirements

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Make](https://www.gnu.org/software/make/) (optional but recommended)

---

## 🧪 Running Locally

1. **Start MongoDB locally** (if not using Docker)
2. **Set environment variables** in `.env`:

```env
MONGO_URI=mongodb://localhost:27017
PORT=5005
```

---

## 🐳 Running with Docker

- Make sure your ```.env``` exists in the project root
- Then build and start the containers:

```
make docker-up
```

- To stop:

```
make docker-down
```

---

## 🧪 API Endpoints

### 📥 POST /register

> #### Registers a new user

```
{
  "username": "john_doe",
  "password": "supersecret"
}
```

### 🔐 POST /login

> #### Logs in an existing user

```
{
  "username": "john_doe",
  "password": "supersecret"
}
```

---

## 🛠 Make Targets

| Command            | Description               |
| ------------------ | ------------------------- |
| `make run`         | Run app locally           |
| `make docker-up`   | Build and run with Docker |
| `make docker-down` | Stop Docker containers    |
| `make build`       | Build Go binary           |
| `make fmt`         | Format code               |
| `make tidy`        | Clean up go.mod/go.sum    |
| `make clean`       | Remove built binary       |

--- 

## 📝 License

This project is licensed under the [MIT License](LICENSE) © [Felipe R. Tasoniero](https://github.com/frtasoniero).