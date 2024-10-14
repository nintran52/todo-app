# **Todo App API Server 📋**

This is the APIs server for the **Todo App**, built using **Golang** and **Gin framework**. The app provides a simple way to manage tasks, including creating, updating, retrieving, and deleting items.

---

## **Features**

- Create a new item
- Retrieve all items
- Retrieve an item by ID
- Update an existing item
- Delete an item by ID
- Well-documented API using **Swagger**

---

## **Prerequisites**

Ensure you have the following installed:

- **Go** (version 1.18 or higher)
- **Docker** (for running with containers)
- **Postman or cURL** (for testing the APIs)
- **Gin Framework** installed

---

## **Installation**

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/todo-app.git
   cd todo-app
   ```

2. Install the required dependencies:

   ```bash
   go mod download
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

### **Run by Docker**

```bash
docker-compose up --build
```

---

## **API Documentation**

Swagger is integrated for easy API exploration.  
After running the server, navigate to `http://localhost:8080/swagger/index.html` to view the Swagger UI.

---

## **API Endpoints**

### **1. Create an Item**

- **Endpoint:** `POST /items`
- **Request Body:**
  ```json
  {
    "title": "Sample Task",
    "description": "This is a sample task"
  }
  ```
- **Response:**
  ```json
  {
    "data": "item_id"
  }
  ```

### **2. Get All Items**

- **Endpoint:** `GET /items`
- **Response:**
  ```json
  {
    "data": [
      {
        "id": "uuid",
        "title": "Sample Task",
        "description": "This is a sample task"
      }
    ]
  }
  ```

### **3. Get Item by ID**

- **Endpoint:** `GET /items/{id}`
- **Response:**
  ```json
  {
    "data": {
      "id": "uuid",
      "title": "Sample Task",
      "description": "This is a sample task"
    }
  }
  ```

### **4. Update an Item**

- **Endpoint:** `PUT /items/{id}`
- **Request Body:**
  ```json
  {
    "title": "Updated Task",
    "description": "This is the updated description"
  }
  ```
- **Response:**
  ```json
  {
    "success": true
  }
  ```

### **5. Delete an Item**

- **Endpoint:** `DELETE /items/{id}`
- **Response:**
  ```json
  {
    "success": true
  }
  ```

---

## **Error Handling**

- **400 Bad Request:** Returned when the request is invalid (e.g., incorrect ID format).
- **404 Not Found:** Returned when the item with the given ID is not found.
- **500 Internal Server Error:** Returned for unexpected server issues.

---

## **Project Structure**

```
/todo-app
│
├── /domain                # Domain models (e.g., Item, ItemUpdate)
├── /internal              # API handlers (Create, Read, Update, Delete)
├── /pkg                   # Client models and error handling
├── /item                  # Business logic and item operations
├── /users                 # Business logic and user operations
├── main.go                # Entry point of the application
├── go.mod                 # Dependencies file
└── README.md              # Project documentation
```

---

## **Run with Docker (Optional)**

1. Build the Docker image:

   ```bash
   docker build -t todo-app .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 todo-app
   ```

---

## **Testing the API**

You can test the API using **Postman** or **cURL**:

Example using `cURL` to create an item:

```bash
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name": "New Task", "description": "Task details"}'
```

---

## **Contributing**

Feel free to contribute by opening a pull request or raising issues. Contributions are welcome!

---

## **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## **Contact**

If you have any questions or need support, reach out at [your-email@example.com](mailto:your-email@example.com).
