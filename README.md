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
git clone https://github.com/Brianhabib252/Golang-Final-Project.git
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

- **Sign Up**
  - Method: `Post`
  - Path: `/signup`
  - Description: endpoint for register a new user
  - Request:
    - Body:
      ```json
      {
          	"age": "int",
          	"email": "string,unique,email",
          	"password": "string,min 6 character",
          	"username": "string,unique"
      }
      ```
  - Response:
    - Status: 201 Created
    - Body:
      ```json
      {
          	"age": "int",
          	"email": "string",
          	"id": "uint",
          	"username": "string"
      }
      ```
- **Sign In**
  - Method: `Post`
  - Path: `/signin`
  - Description: endpoint for sign in and get token by cookies
  - Request:
    - Body:
      ```json
      {
          	"email": "string,unique,email",
          	"password": "string,min 6 character",
      }
      ```
  - Response:
    - Status: 200 OK
    - Body:
      ```json
      {
          	"token": "string",
          	"id": "uint",
      }
      ```
- **Create Social Media**
  - Method: `POST`
  - Path: `/sosmed`
  - Description: post your social media and biodata to the database my sql
  - Request:
    - Body:
      ```json
      {
          "name":"string",
          "social_media_url":"string",
          "height":"int"(cm),
          "weight":"int"(kg),
          "gender":"string"(MALE or FEMALE)
      }
      ```
  - Response:
    - Status: 201 Created
    - Body:
      ```json
      {
          "name":"string",
          "social_media_id":"uint",
          "social_media_url":"string",
          "height":"int"(cm),
          "weight":"int"(kg),
          "gender":"string"(MALE or FEMALE)
          "user": {
              "email": "string",
              "username": "string"
          },
      }
      ```
- **Get All Social Media**
  - Method: `GET`
  - Path: `/sosmed`
  - Description: get all data of social media and biodata from the database my sql
  - Request: none
  - Response:
    - Status: 200 OK
    - Body:
      ```json
      {
          "name":"string",
          "social_media_id":"uint",
          "social_media_url":"string",
          "height":"int"(cm),
          "weight":"int"(kg),
          "gender":"string"(MALE or FEMALE)
          "user": {
              "email": "string",
              "username": "string"
          },
      }
      ```

- **Get By ID Social Media**
  - Method: `GET`
  - Path: `/sosmed/:social_media_id`
  - Description: get data of social media and biodata by social media id from the database my sql
  - Request:
    - params : social_media_id
  - Response:
    - Status: 201 Created
    - Body:
      ```json
      {
          "name":"string",
          "social_media_id":"uint",
          "social_media_url":"string",
          "height":"int"(cm),
          "weight":"int"(kg),
          "gender":"string"(MALE or FEMALE)
          "user": {
              "email": "string",
              "username": "string"
          },
      }
      ```
- **Update Social Media**
  - Method: `PUT`
  - Path: `/sosmed/id/social_media_id`
  - Description: update your social media and biodata to the database my sql
  - Request:
    - params : id(dser id), social_media_id
    - Body:
      ```json
      {
          "name":"string",
          "social_media_url":"string",
          "height":"int"(cm),
          "weight":"int"(kg),
      }
      ```
  - Response:
    - Status: 200 OK
    - Body:
      ```json
      {
          "name":"string",
          "social_media_id":"uint",
          "social_media_url":"string",
          "height":"int"(cm),
          "weight":"int"(kg),
          "gender":"string"(MALE or FEMALE)
          "user": {
              "email": "string",
              "username": "string"
          },
      }
      ```
- **Delete Social Media**
  - Method: `POST`
  - Path: `/sosmed/id/social_media_id`
  - Description: delete your social media and biodata from the database my sql
  - Request:
    - params : id(dser id), social_media_id
  - Response:
    - Status: 200 OK
    - Body:
      ```json
      {
          "message": "Social media deleted successfully"
      }
      ```


## Examples

Provide some examples or code snippets demonstrating how to use your API.

