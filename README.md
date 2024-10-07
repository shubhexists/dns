# 🌐 DNS Server and Platform

A lightning-fast, feature-rich DNS server and platform implemented in Go, boasting advanced capabilities including NameSpace Resolution.

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14.0%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Project Status](https://img.shields.io/badge/Status-Active-brightgreen)

## 🚀 Features

- 📡 Complete DNS server functionality
- 🔍 NameSpace Resolution support
- 🏎️ Built with Go for high performance and concurrency
- 💾 Utilizes PostgreSQL for robust data storage
- ⚡ Integrates DiceDB (a Redis alternative) for caching and fast data retrieval
- 🌐 Exposes UDP port for handling DNS connections
- 🖥️ Provides HTTP interface for database connections and management

## 🚫 Intentionally Omitted Features

In the interest of maintaining a streamlined and modern DNS server, we've intentionally excluded certain obsolete or rarely used features:

### Obsolete QTypes

The following obsolete QTypes are not implemented:

- QTYPE_MD (3): Mail Destination
- QTYPE_MF (4): Mail Forwarder
- QTYPE_MB (7): Mailbox Domain Name
- QTYPE_MG (8): Mail Group Member
- QTYPE_MR (9): Mail Rename Domain Name

These QTypes are no longer in use in modern DNS systems, so we've opted to exclude them for simplicity and efficiency.

Altough this is subjected to change on the later stages of the project.

### Class Limitation

This DNS server only supports the Internet (IN) class. Other classes are not implemented as they are rarely used in contemporary DNS deployments.

These design decisions help keep our codebase clean, efficient, and focused on modern DNS use cases.

## 🛠️ Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git

## 📦 Installation and Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/shubhexists/dns
   cd dns
   ```

2. Set up the environment file:

   ```bash
   cp env.sample .env
   ```

   Edit the `.env` file to configure your environment variables as needed.

3. Start the required services using Docker Compose:

   ```bash
   docker compose --env-file .env up -d
   ```

   This command will start PostgreSQL, DiceDB, and any other required services defined in your `docker-compose.yml` file.

4. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## 🚀 Usage

To start the DNS server:

```bash
go run cmd/main.go
```

This command will run your main application, connecting to the services you've set up with Docker Compose. Once started:

- The server will listen for DNS queries on the configured UDP port.
- An HTTP server will be available for database connections and management tasks.

Make sure your firewall allows traffic on the configured UDP port for DNS queries and the HTTP port for management.

## ⚙️ Configuration

The main configuration is handled through the `.env` file. Make sure to review and adjust the variables in this file according to your needs. Key configuration points include:

- UDP port for DNS connections
- HTTP port for database management
- Database connection details
- DiceDB connection details

## 📚 API Documentation

Yet to be Documented, Please stay tuned

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **shubhexists** - [shubhexists](https://github.com/shubhexists)
- **well-mannered-goat** - [well-mannered-goat](https://github.com/well-mannered-goat)

## Thanks

Do star ⭐ the project to give us a boost. Would mean a lot! 