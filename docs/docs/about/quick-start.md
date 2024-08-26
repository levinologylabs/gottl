---
---

# Getting Started

## Pre-requisites

First off, you need to have the following installed on your machine:

- [Taskfile](https://taskfile.dev/#/installation) - Task runner
- [Docker](https://docs.docker.com/get-docker/) - Containerization platform
- [Go 1.22+](https://golang.org/dl/) - Programming language
- [Scaffold](https://hay-kot.github.io/scaffold/introduction/quick-start.html) - Scaffolding CRUD operations
- [Swaggo](https://github.com/swaggo/swag?tab=readme-ov-file#getting-started) - Generating Swagger documentation
- [Sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) - Generating Go code from SQL queries
- [Dagger.io](https://dagger.io/) - Running Pipelines
- [Golangci-lint](https://golangci-lint.run/usage/install/) - Linting Go code

## Starting Your Project

### Step 1: Clone the Repository

```bash
git clone https://github.com/levinologylabs/gottl
```

### Step 2: Ensure You Have the Required Tools

Ensure that you have working environment before you start to make changes. You can check by running the following command:

```bash
task run
```
This should:

1. Start the required containers
2. Run the JSON Api server

### Step 3: Modifying Go Module

You can modify the `go.mod` file to include your module name. You can do this by running the following command:

```bash
// TODO
```