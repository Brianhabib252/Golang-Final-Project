# Golang Final Project My Food Gram
this api have fitur to upload a foto, upload your biodata, make a comment for a photo, and fitur to track your food and the nutrition that you have eat and have a security for authentication and authorization  

## Table of Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Authentication](#authentication)
  - [Endpoints](#endpoints)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

- Go 1.15 or higher
- MySQL database
- Postman or any other REST API client for testing

### Installation

1. Clone the repository:

```bash
git clone <repository_url>
```

2. Navigate to the project directory:

```bash
cd Golang-Final-Project
```

3. Install dependencies:

```bash
go mod tidy
```

4. Set up environment variables by creating a `.env` file and filling it with the required configurations:

```bash
DB_USERNAME=your mysql username
DB_PASSWORD=your mysql password
DB_HOST=your mysql host
DB_PORT=your mysql port
DB_NAME=your mysql database name
SECRET_KEY=secret key for the jwt authentication
...

# Add other necessary environment variables
```

5. Run the application:

```bash
go run main.go
```

## Usage

### Authentication

The authentication using a jwt and bycrypt package so the user password gonna be hash and save in the database. jwt token will automatically generated while login and it will be set as a cookie.
the authorization using a jwt and id params that will use in every endpoint that need authorization, so it will be check if the id params and the user id get from the jwt cookies are match.

### Endpoints

List all the endpoints available in your API along with their descriptions and example requests/responses.

- **Sign Up** endpoint for register a new user
  - Method: `Post`
  - Path: `http://localhost:8080/signup`
  - Description: endpoint for register a new user
  - Request:
    - Body:
      ```json
      {
          	"age": int,
          	"email": string,unique,email,
          	"password": string,min 6 character,
          	"username": string,unique
      }
      ```
  - Response:
    - Status: 200 OK
    - Body:
      ```json
      {
          "key": "value"
      }
      ```

- **Endpoint 2**: Description of endpoint 2
  - Method: `POST`
  - Path: `/endpoint2`
  - Description: Description of what this endpoint does
  - Request:
    - Body:
      ```json
      {
          "key": "value"
      }
      ```
  - Response:
    - Status: 201 Created

## Examples

Provide some examples or code snippets demonstrating how to use your API.

## Contributing

Explain how people can contribute to your project if you are open to contributions.

## License

Specify the license under which your project is distributed.

---

Feel free to customize this template according to your specific API requirements and structure.
