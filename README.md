# go-task-challenge

## Overview

go-task-challenge CLI is a command-line tool designed to manage customer registrations and KYC (Know Your Customer) verifications. The application is structured following Domain-Driven Design (DDD) principles to ensure that the business logic is well-organized, modular, and easy to maintain.

## Problem Statement

In any customer-facing application, ensuring the legitimacy of customer data is crucial. This involves two critical tasks:

1. **Customer Registration**: The system must allow customers to register with their basic information and set an initial KYC status.
2. **Customer KYC Verification**: The system should verify the customer's identity through multiple external services and update their KYC status accordingly.

The challenge lies in structuring the application in a way that is modular, scalable, and easy to maintain. Additionally, verifying customer data against multiple external services can be resource-intensive, so efficient concurrency patterns are necessary.

## Solution

To address the problem, go-task-challenge CLI implements the following:

1. **Domain-Driven Design (DDD)**: The application is structured into clearly defined layers (Domain, Application, Infrastructure), each handling a specific aspect of the business logic.
2. **Concurrency Patterns**: Specifically, the Fan-in and Fan-out concurrency patterns are used to handle KYC verification against multiple external services efficiently.

## Tasks

### Task 1: Register Customer

**Objective**: Implement a command that registers a customer and sets their initial KYC status to "pending."

**Steps**:
1. **Domain Layer**: Define the `Customer` entity, which represents a customer in the system. This entity includes attributes such as `FirstName`, `LastName`, `Email`, `Phone`, `Address`, and `KYCStatus`.
2. **Application Layer**: Implement the `CustomerService` class, which handles the registration process.
3. **Infrastructure Layer**: Implement an in-memory repository (`InMemoryCustomerRepository`) to store customer data.
4. **Command Layer**: Create the `registerCmd` command, which accepts user input (first name, last name, email, phone, address) and uses the `CustomerService` to register the customer.

**Improvements**:
- **Modularity**: By separating the customer registration process into different layers, the system becomes more maintainable and easier to extend. However, this implementation could be improved by integrating a persistent database solution instead of an in-memory repository, which would allow the system to handle more extensive data sets and retain information across sessions.
- **Error Handling**: The current implementation includes basic error handling, but this could be enhanced by implementing more detailed logging and error tracking mechanisms. This would make the system more robust in production environments where identifying and resolving issues quickly is crucial.
- **Scalability**: While the current setup works well for a small number of customers, as the system scales, we might need to optimize how customer data is stored and accessed, potentially integrating caching mechanisms to improve performance.

### Task 2: Verify Customer KYC

**Objective**: Implement a command that verifies a customer's KYC status using multiple external services concurrently.

**Steps**:
1. **Domain Layer**: Define the `KYCService` interface that abstracts the KYC validation process.
2. **Application Layer**: Extend the `CustomerService` class to include a `VerifyCustomerKYC` method. This method uses the Fan-in and Fan-out concurrency patterns to validate the customer's KYC status across multiple external services.
3. **Infrastructure Layer**: Implement the `KYCAdapter`, which acts as an adapter for the external KYC validation services. The service randomly approves or rejects customers to simulate a real-world scenario.
4. **Command Layer**: Create the `verifyCmd` command, which accepts the customer's email as input and uses the `CustomerService` to verify the customer's KYC status.

**Concurrency Pattern**:
- **Fan-Out**: The `VerifyCustomerKYC` method fans out the verification requests to multiple external services concurrently. Each service operates in its own goroutine, allowing the system to handle multiple verification processes simultaneously.
- **Fan-In**: The results from all external services are fanned back in to aggregate the final KYC status. This ensures that the verification process is completed efficiently and that the system can determine the final KYC status based on multiple sources.

**Improvements**:
- **Efficiency**: By using the Fan-in and Fan-out concurrency patterns, the system efficiently handles multiple external service calls, reducing overall verification time. However, the implementation can be improved by implementing a more sophisticated aggregation strategy. For instance, weighting the reliability of different KYC services and using that to influence the final decision.
- **Error Handling**: The current implementation does not fully handle scenarios where one or more external services fail. Future improvements could include retry mechanisms or fallbacks to ensure the system remains resilient even if some services are temporarily unavailable.
- **Scalability**: As the number of external services increases, the current approach may face bottlenecks. To improve scalability, we could implement a more dynamic concurrency management system that adjusts the number of goroutines based on system load and external service response times.

### Task 3: Implement a Simplified Redis-Like Cache System using Fan-in and Fan-out Concurrency Patterns in Go

#### Problem Statement:
Design and implement a simplified Redis-like in-memory key-value store in Go. The system should be capable of handling concurrent read and write requests efficiently. 

To achieve this, use the Fan-In and Fan-Out concurrency design patterns. Specifically, implement the following features:

1. **SET command**: Store a key-value pair in the cache.
2. **GET command**: Retrieve the value for a given key from the cache.
3. **DEL command**: Delete a key-value pair from the cache.
4. **TTL**: Implement time-to-live (TTL) functionality for each key, where the key-value pair automatically expires after the TTL duration.

#### Concurrency Requirements:
1. **Fan-Out**: When multiple clients send requests to the cache system, each request should be processed concurrently.
2. **Fan-In**: Use a worker pool to handle these requests efficiently, gathering results back in a centralized manner.

#### Additional Constraints:
- The cache should be thread-safe.
- Use channels to implement the Fan-In and Fan-Out patterns.
- Implement a mechanism to automatically clean up expired keys using TTL.

#### Key Components:

1. **Cache Data Structure**:
    - A `map[string]interface{}` to store the key-value pairs.
    - A `map[string]time.Time` to track TTL expiration for each key.
  
2. **Worker Pool**:
    - A fixed-size pool of goroutines to handle incoming requests.
  
3. **Request Channels**:
    - A channel for incoming requests (Fan-Out).
    - A channel for worker results (Fan-In).

4. **Request Handler**:
    - Handles `SET`, `GET`, and `DEL` commands.
    - Manages TTL expiration.

5. **Cleaner Goroutine**:
    - Periodically checks and removes expired keys.


## How to Run the Application

1. **Clone the Repository**:
   ```
   git clone https://github.com/macadrich/go-task-challenge.git
   cd go-task-challenge
   ```

2. **Run the Application**:
   ```
   go run main.go
   ```

3. **Register a Customer**:
   ```
   Enter command: register --first-name John --last-name Doe --email john.doe@example.com --phone 1234567890 --address "123 Main St"
   ```

4. **Verify Customer KYC**:
   ```
   Enter command: verify --email john.doe@example.com
   ```

5. **Redis-Cache: Set Key-Value with TTL of 60 seconds**:
   ```
   Enter command: set mykey myvalue -t 60

6. **Redis-Cache: Get Value Redis-Cache**:
   ```
   Enter command: get mykey

7. **Exit the Application**:
   ```
   Enter command: exit
   ```