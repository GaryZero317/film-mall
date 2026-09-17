# Film Mall

English | [简体中文](README.zh-CN.md)

An e-commerce platform for photographic film, built as my undergraduate final-year project between **December 2024 and March 2025**.

Film Mall combines an online store with photography community features. It includes Go backend services, a Vue administration interface, and a WeChat Mini Program for customers.

## Features

- **Users and addresses:** registration, login, profile management, and delivery addresses.
- **Product catalogue:** product categories, product details, images, and search.
- **Shopping cart:** add and remove items, update quantities, and select items for checkout.
- **Orders:** order creation, status management, cancellation, and unpaid order expiry handling.
- **Photography community:** share photographs, browse posts, and add likes and comments.
- **Customer support:** chat, FAQs, and service requests.
- **Administration:** manage products, orders, users, community content, and sales statistics.
- **Film services:** film-related orders and photo management.

## Technology Stack

| Area | Technologies |
| --- | --- |
| Backend | Go, go-zero, gRPC |
| Database | MySQL, GORM |
| Cache and expiry tracking | Redis |
| Service discovery | etcd |
| Administration interface | Vue 3, Vite, Element Plus, Pinia, ECharts |
| Customer interface | WeChat Mini Program |
| Authentication | JWT |
| Real-time communication | WebSocket |
| API testing | Postman collection included in the repository |

## Backend Implementation Highlights

### Product caching

Product lookup checks the cache before querying MySQL. The database provides a fallback when cached product data is missing or invalid.

Implementation: [product_gorm.go](service/product/model/product_gorm.go).

### Concurrent inventory task processing

Inventory deduction and restoration requests are placed in a buffered Go channel and processed by a pool of worker goroutines. Results are returned through per-request channels, and requests are rejected when the queue is full.

Implementation: [stock_processor.go](service/product/model/stock_processor.go).

### Unpaid order expiry

Redis sorted sets track order expiry times. A background task periodically retrieves expired entries and initiates cancellation for unpaid orders.

Implementation: [ordertimeout.go](service/order/rpc/internal/task/ordertimeout.go).

These sections describe the mechanisms implemented in the project, rather than measured throughput or production-scale performance guarantees.

## Repository Structure

```text
film-mall/
├── admin-frontend/          # Vue administration interface
├── Wechat/                 # WeChat Mini Program project
├── service/                # Backend services
│   ├── user/               # Users, authentication, and customer support
│   ├── product/            # Products and inventory
│   ├── order/              # Orders and expiry handling
│   ├── cart/               # Shopping cart
│   ├── pay/                # Payment-related APIs and RPC service
│   ├── address/            # Delivery addresses
│   ├── community/          # Photography posts, likes, and comments
│   ├── film/               # Film services and photos
│   └── statistics/         # Sales and activity statistics
├── common/                 # Shared backend utilities
├── docs/                   # Supporting documentation
├── mall.sql                # Database SQL script
├── gomall.postman_collection.json
├── go.mod
└── go.sum
```

## Local Setup

### Prerequisites

- **Go:** `go.mod` declares Go 1.22.7 and the Go 1.23.3 toolchain.
- **MySQL 8.0**, **Redis**, and **etcd** running locally or on accessible hosts.
- **Node.js and npm:** use a Node.js version compatible with Vite 6, as declared by the frontend dependencies.
- **WeChat Developer Tools** to run the customer Mini Program.

Service configuration files contain environment-specific connection settings. Update them for your own environment before starting the application. The repository contains multiple services that must be started separately.

### 1. Clone the repository

```bash
git clone https://github.com/GaryZero317/film-mall.git
cd film-mall
go mod download
```

### 2. Prepare the database

Create a dedicated development database named `mall`, or choose another name and update the service database connections accordingly.

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS mall CHARACTER SET utf8mb4;"
mysql -u root -p mall < mall.sql
```

The SQL script recreates tables, so import it only into a database intended for this project.

### 3. Configure the backend

Review the YAML files under `service/<service>/api/etc/` and `service/<service>/rpc/etc/` where present. Configure:

- MySQL connection details and database name.
- Redis addresses and authentication settings.
- etcd endpoints and RPC discovery keys.
- API and RPC listening ports.
- JWT settings, keeping the signing configuration consistent between related services.
- Any upload directories and service-specific paths used in your environment.

Keep private credentials out of public commits. In `service/product/rpc/etc/product.yaml`, consolidate the repeated top-level `Etcd` block into a single block before starting the product RPC service.

### 4. Start the backend services

Start MySQL, Redis, and etcd first. Run each backend service in its own terminal, from its service directory, so relative paths resolve correctly.

For example, start the product RPC service:

```bash
cd service/product/rpc
go run product.go -f etc/product.yaml
```

In a separate terminal, starting from the repository root, start the product API:

```bash
cd service/product/api
go run product.go -f etc/product.yaml
```

Use the following entry points for the other services. Commands in the last column are run from the directory in the middle column.

| Service | Directory | Command |
| --- | --- | --- |
| User RPC | `service/user/rpc` | `go run user.go -f etc/user.yaml` |
| Order RPC | `service/order/rpc` | `go run order.go -f etc/order.yaml` |
| Payment RPC | `service/pay/rpc` | `go run pay.go -f etc/pay.yaml` |
| User API | `service/user/api` | `go run user.go -f etc/user.yaml` |
| Order API | `service/order/api` | `go run order.go -f etc/order.yaml` |
| Payment API | `service/pay/api` | `go run pay.go -f etc/pay.yaml` |
| Cart API | `service/cart/api` | `go run cart.go -f etc/cart.yaml` |
| Address API | `service/address/api` | `go run address.go -f etc/address-api.yaml` |
| Statistics API | `service/statistics/api` | `go run statistics.go -f etc/statistics.yaml` |
| Film API | `service/film/api` | `go run film.go -f etc/film.yaml` |
| Community API | `service/community/api` | `go run community.go -f etc/community.yaml` |

Start the RPC services before exercising APIs that depend on them. The default API ports are 8000–8008; the RPC ports are 9000–9003. Check the YAML files for the active settings.

### 5. Start the administration interface

From the repository root:

```bash
cd admin-frontend
npm install
npm run dev
```

The Vite configuration uses port **5173**. Check `admin-frontend/vite.config.js` and the frontend API configuration to ensure requests point to your running backend services.

### 6. Run the WeChat Mini Program

Open the project under `Wechat/` in WeChat Developer Tools. Configure the appropriate AppID and update the Mini Program's API addresses for your environment. For testing on a physical device, use a backend address reachable from that device.

## API Testing

Import [gomall.postman_collection.json](gomall.postman_collection.json) into Postman. Review the request URLs, authentication tokens, and sample identifiers before running requests against your local services.

A useful manual check is to register or log in, browse products, add an item to the cart, create an order, and inspect the corresponding order and inventory records.

## Project Background

This project was developed over several months as part of my undergraduate Software Engineering degree. It provided practical experience building an application across customer interfaces, backend services, database access, caching, and background task processing.

**Development period:** December 2024 – March 2025.
