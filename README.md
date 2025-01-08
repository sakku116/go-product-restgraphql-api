# Go Product REST & GraphQL API

This repository contains a backend API built with Go, supporting both RESTful and GraphQL endpoints. The application provides CRUD operations for users and products, secured with JWT authentication, and uses MongoDB as the database.

## Features

- **User Management**
  - Create, Read, Update, and Delete (CRUD) operations for users.

- **Product Management**
  - CRUD operations for products.

- **Authentication**
  - JWT-based authentication for secure access to APIs.

- **RESTful API**
  - API documentation and testing via Swagger UI available at `/swagger/index.html`.

- **GraphQL API**
  - GraphQL endpoint for authentication at `/auth/graphql` (playground at `/auth/graphql/playground`).
  - GraphQL endpoint for users at `/users/graphql` (playground at `/users/graphql/playground`).
  - GraphQL endpoint for products at `/products/graphql` (playground at `/products/graphql/playground`).

- **Database**
  - MongoDB for data storage.

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sakku116/go-product-restgraphql-api.git
   cd go-product-restgraphql-api
   ```

2. **Run the application using Docker Compose**: Ensure you have Docker and Docker Compose installed, then run:

    ```bash
    docker-compose up -d
    ```
## APIs

- RESTful API Documentation: http://localhost:8001/swagger/index.html
- GraphQL:
    - Auth:
        - http://localhost:8001/auth/graphql
        - http://localhost:8001/auth/graphql/playground
    - Users:
        - http://localhost:8001/users/graphql
        - http://localhost:8001/users/graphql/playground
    - Products:
        - http://localhost:8001/products/graphql
        - http://localhost:8001/products/graphql/playground