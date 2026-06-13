# 📝 TO-DO LIST API

[English](#english) |

---

## English

A clean and robust RESTful API for a To-Do List application, built with **Go**, **PostgreSQL**, and **Docker**. This is my first Go project designed to demonstrate hands-on experience with backend development, containerization, and database integration.

### 🚀 Features
* **Full CRUD Lifecycle:** Create, Read, Update, and Delete tasks.
* **Database Persistence:** Integrated with PostgreSQL via `pgx/v5` driver.
* **Environment Configuration:** Secure credential management using `.env` files and `godotenv`.
* **Containerization:** Ready-to-go database infrastructure managed via Docker Compose.

### 🛠️ Tech Stack
* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Tools:** Docker, Docker Compose, Postman (for API testing)
* **Libraries:** `github.com/jackc/pgx/v5`, `github.com/joho/godotenv`

### 📦 Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/GudOwner/TO-DO-list-by-Gud.git](https://github.com/GudOwner/TO-DO-list-by-Gud.git)
   cd TO-DO-list-by-Gud

2. **Configure environment variables**
Create a .env file in the root directory based on
3. **Spin up the database (Docker)**
    ```bash 
    docker compose up -d 

4. **Run the GO server**
    ```bash 
    go run main.go 




### 🛣️ API Endpoints

| Method | Endpoint | Description | Sample Request Body (JSON) |
| :--- | :--- | :--- | :--- |
| **GET** | `/task` | Get all tasks | *None* |
| **POST** | `/task` | Create a new task | `{"task": "Buy protein", "status": false}` |
| **PUT** | `/task` | Update an existing task | `{"id": 1, "task": "Buy protein", "status": true}` |
| **DELETE** | `/task` | Delete a task | `{"id": 1}` |
