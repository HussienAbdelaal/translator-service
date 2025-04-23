# Translator Service

The Translator Service is a Go-based application that leverages OpenAI's GPT models to translate text from Arabic to English while leaving any given English text unchanged. It is designed to handle large batches of text efficiently, splitting and reconstructing transcriptions as needed.

## Features

- **Translation**: Translates Arabic text to English while preserving any English text as-is.
- **Batch Processing**: Splits large transcriptions into smaller parts for efficient processing and reconstructs them after translation. Also group small transcriptions together into one batch to avoid having a lot of small requests.
- **Database Integration**: Stores and retrieves previously translated scripts.
- **Parallel Execution**: Executes translation tasks in parallel with fail-fast behavior for improved performance.
- **REST API**: Exposes endpoints for retrieving and translating text.

## Architecture

The service is built with a modular architecture, including the following components:

- **Handlers**: Define REST API endpoints using the Gin framework.
- **Services**: Contain the business logic for translation and batch processing.
- **Repositories**: Handle database interactions for storing and retrieving translations.
- **Clients**: Integrate with external services like OpenAI's GPT models.
- **Utilities**: Provide helper functions for tasks like hashing, parallel execution, and text splitting.

## Endpoints

### 1. `GET /translations`
Retrieves all stored translations from the database. This is for development purposes (no pagination implemented).

### 2. `POST /translate`
Translates a batch of text inputs. The request body should be a JSON array of transcriptions with the following structure:
```json
[
  {
    "speaker": "Speaker1",
    "time": "00:00:01",
    "sentence": "مرحبا"
  }
]
```

The response will include the translated text in the same order as the input with the same structure:
```json
[
  {
    "speaker": "Speaker1",
    "time": "00:00:01",
    "sentence": "Hello"
  }
]
```


## Setup

### Prerequisites

- Go 1.24 or later
- OpenAI API key
- PostgreSQL (or any other supported database)

### Installation

#### Install Go
- Install Go from the official website: [golang.org](https://golang.org/dl/) and follow the instructions for your operating system.
- Make sure go is installed by running:
```bash
go version
```

#### Clone the Repository
```bash
git clone git@github.com:HussienAbdelaal/translator-service.git
cd translator-service
```

#### Install Dependencies
```bash
go mod tidy
```

#### Create PostgreSQL Database
1. Using Terraform
    - Install Terraform from the official website: [terraform.io](https://www.terraform.io/downloads.html).
    - Navigate to the `terraform-postgres` directory and run:
      ```bash
      terraform init
      terraform apply
      ```
    - This will create a PostgreSQL database instance on AWS RDS.
    - Don't forget to destroy the resources after use:
      ```bash
      terraform destroy
      ```
2. Using Docker
    - Install Docker from the official website: [docker.com](https://www.docker.com/get-started).
    - Run the following command to start a PostgreSQL container:
      ```bash
      docker run --name postgres -e POSTGRES_USER=<your_db_user> -e POSTGRES_PASSWORD=<your_db_password> -p 5432:5432 -d postgres
      ```
    - This will create a PostgreSQL database instance in a Docker container.

#### Run Migration Script
- start by creating the database inside the PostgreSQL instance:
```sql
postgres=> CREATE DATABASE <your_db_name>;
```
- The migrations are created using [golang-migrate](https://github.com/golang-migrate/migrate) and can be run using the following command:
```bash
migrate -path db/migrations -database postgres://<your_db_user>:<your_db_password>@<db_host>:5432/<your_db_name> up
```
- However, the migration script is already included in the repository, so you can run it directly without using golang-migrate:
```bash
psql -h <db_host> -U <your_db_user> -d <your_db_name> -f db/migrations/000001_create_translation_table.up.sql
```

#### Set Up Environment Variables
Create a `.env` file in the root directory with the following variables:
```plaintext
OPENAI_API_KEY=<your_openai_api_key>
OPENAI_MODEL=gpt-4o-mini
OPENAI_BATCH_SIZE=6000
OPENAI_TEMPERATURE=0.3

DB_USER=<your_db_user>
DB_PASSWORD=<your_db_password>
DB_HOST=<your_db_host>
DB_PORT=5432
DB_NAME=<your_db_name>
```
- If you created the database using Terraform, the `DB_HOST` can be found using the AWS console or running the following command:
```bash
aws rds describe-db-instances --query "DBInstances[?DBInstanceIdentifier=='<your_db_instance_name>'].Endpoint.Address" --output text
```
or 
```
cd terraform-postgres
terraform show | grep -i endpoint
```

### Build and Run

#### Build and Run Locally

1. Build and run the application:
```bash
go build -o translator
./translator
```
2. Alternatively, you can run the application directly without building
```bash
go run main.go
```
The service will be available at `http://localhost:8080`.

#### Docker
A `Dockerfile` is provided for containerizing the application. To build and run the container:

1. Build the Docker image:
   ```bash
   docker build -t translator-service .
   docker run --rm --env-file .env -p 8080:8080 translator-service
   ```

### Testing
Unit tests are provided for key components. To run the tests:
```bash
go test ./...
```
