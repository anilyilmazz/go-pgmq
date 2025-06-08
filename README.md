# Go-PGMQ

A simple queue-based processing system written in Go. It creates a message queue using `pgmq` on PostgreSQL, allows sending messages via REST API, and processes messages in the background with Go workers.

> 📖 **Related Article**: [RabbitMQ’suz da Olur: PostgreSQL ile Sade ve Güçlü Bir Kuyruk Sistemi](https://medium.com/@anilyilmaz/rabbitmqsuz-da-olur-postgresql-ile-sade-ve-g%C3%BC%C3%A7l%C3%BC-bir-kuyruk-sistemi-337e8bdb9823)
> 
> This project is a practical implementation of the concepts discussed in the Medium article, demonstrating how to build a PostgreSQL-based worker pool architecture as an alternative to managed queue services.

## 🧩 Features

- PostgreSQL (`pgmq`) based message queue
- Send messages to the queue via REST API
- Worker processes messages using `pgmq_read` and `pgmq_delete`
- Supports multiple workers (using Go routines)
- Easy start with Docker Compose

## 📦 Technologies Used

- **Golang 1.23** - Backend language
- **PostgreSQL + pgmq** - Message queue system
- **REST API** -  Built with Go's standard HTTP package
- **Docker & Docker Compose** - Containerization and orchestration

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/anilyilmazz/go-pgmq.git
cd pgmq-worker-api
```

### 2. Run the project with Docker Compose

```bash
docker-compose up
```

### 3. Test the API

After the services are running, you can send a test message to the queue using this command:

```bash
curl --location 'http://localhost:8080/send' \
--header 'Content-Type: application/json' \
--data '{"type": "email", "payload": "Hello from API"}'
```

After sending the message, you can check your Docker logs to see the message being received and processed by the worker.

## 📝 API Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/send` | Send a message to the queue | `{"type": "string", "payload": "string"}` |

## 🏗️ Project Structure

```
go-pgmq/
├── consumer.go          # Message queue consumer/worker logic
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile          # Docker build instructions
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
├── http_handler.go     # HTTP API handlers
├── main.go             # Main application entry point
```

## 🔧 Configuration

The application uses the following default configuration:

- **API Port**: 8080
- **PostgreSQL**: Running in Docker container
- **Queue Name**: Configurable in the application

## 🧪 Testing

To verify the system is working:

1. Start the services with `docker-compose up`
2. Send a test message using the curl command above
3. Check Docker logs to see message processing:
   
 ```bash
 docker-compose logs -f
 ```

## 🙏 Acknowledgments

- [pgmq](https://github.com/pgmq/pgmq) for the PostgreSQL message queue extension
