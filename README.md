# Arritech Hometask Challenge 

Golang + Vue.js Hometask challenge

## Getting Started
This application is written in Golang + Vue.js let you manage users along with MongoDB

You can easily download the source from github <br>
`[Github Repo](https://github.com/secret-point/simple-user-management)`

### Prerequisites

What things you need to install the software and how to install them:

- Go (Programming language)
- Node.js (JavaScript runtime)

Ensure you have Go installed by running `go version` from your command line. The project requires a Go version that supports modules (Go 1.11+).

Ensure you have Node.js installed by running `node -v` from your command line. The project requires Node.js version 18 or higher.

### Installing

A step-by-step series of examples that tell you how to get a development environment running.

#### Backend Setup

From the root directory

```bash
cd backend
```

Install backend packages:

```bash
go mod tidy
```

Run the backend server:

```bash
go run main.go
```

Frontend Setup
From the root directory

```bash
cd frontend
```

Install frontend modules:

```bash
npm i
```

Run the project in development mode:

```bash
npm run dev
```

## Technologies

### Frontend

- Vue v3.4.21
- Typescript v5.4.0
- lodash: v4.17.21
- axios v1.6.8
- yup: v1.4.0
- jest: v29.7.0

### Backend

- go v1.21.3
- mongodb
