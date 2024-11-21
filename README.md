# Microservices Architecture with Go and Kubernetes

## Project Overview

A modern microservices-based application built with Go, featuring:

- API Gateway for routing and service orchestration
- Product Service for product management
- User Service for authentication and user management
- Payment Service for transaction processing
- Notification Service for alerts and communications
- Delivery Service for shipment tracking

## Architecture

```mermaid
graph TD
   Client --> Gateway[API Gateway]
   Gateway --> PS[Product Service]
   Gateway --> US[User Service]
   Gateway --> PAY[Payment Service]
   Gateway --> NS[Notification Service]
   Gateway --> DS[Delivery Service]

   PS --> DB[(PostgreSQL)]
```
