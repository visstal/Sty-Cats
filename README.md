# Spy Cat Agency

## Quick Start

```bash
# 1. Navigate to the project directory
cd Sty-Cats

# 2. Start everything
docker-compose up --build
```

## Access Points

- **Frontend Application**: http://localhost:4300/
- **API Documentation**: http://localhost:3001/swagger/index.html
- **Backend API**: http://localhost:3001

## Features & Implementation Notes

### Frontend
For your comfort, i vibecoded simple frontend app - http://localhost:4300/  
But you can check API docs right into http://localhost:3001/swagger//index.html#/

### External API Integration
https://api.thecatapi.com/v1/breeds is stored in local memory cache, with 1hr TTL

### Architecture
I used DDD-like arclitecture, and implemented Rich Domain Model as much as time let me

### Logging
Logs - they are structured, for each request-response, and DB query  
Are written directly into docker container logs.

## Database Implementation

### Transactions
I used transactions for multi-table updates, like completing mission

### ORM Choice
I used GORM, because its the fastest way to develop, although there are drawbacks  
If i had more time - i would have chosen SQLC

### Concurrency Handling
There is no handling of possible DB anomalies, again, because of limited time  
In real app it worth to implement optimistic locks, or even Pessimistic, because  
SCA does not look like highload app

## What's Missing

Due to time constraints, the following features are not implemented:

- Proper handling if https://api.thecatapi.com/v1/breeds is unavaiable
- Auth - as its not required
- Tests - same
- Rate liniting
