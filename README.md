# Tasker

Tasker is a small project aimed at providing a practical platform for practicing various aspects of web development. It's a Go-based application that follows the hexagonal design pattern, utilizes adapters to interact with external APIs, connects to MySQL and Redis databases, employs go-chi as a router library, utilizes stretchr/testify for testing, and relies on DATA-DOG/go-sqlmock for database testing. This project includes a Docker Compose file for easy local setup.

## Table of Contents

- [Introduction](#introduction)
- [Usage](#usage)
- [Features](#features)
- [Getting Started](#getting-started)
- [Endpoints](#endpoints)
- [License](#license)

## Introduction

Tasker is an application designed to assist you in creating and managing tasks with complex sequences of steps. It is particularly useful for scenarios where you need to automate a series of actions or processes, with each step being customizable and flexible.

Here's what Tasker can do:

- **Task Creation**: Define the tasks you want to automate. Each task consists of a series of different steps, and you can specify what each step should do, including making API calls, data manipulation, or other actions.

- **Data Flow**: You have full control over how data flows between steps. Each step can use data provided by the previous step or include its parameters, allowing you to build flexible and interconnected workflows.

- **Error Handling**: For robust execution, Tasker supports error handling at each step. You can specify how the application should respond if an error occurs during a step's execution.

- **Task Scheduling**: If you want tasks to run automatically, you can schedule them using cron syntax. Define the schedule for each task, and Tasker will ensure they execute at the specified times.

## Usage

Here's how you can use Tasker:

1. **Create Tasks**: Define the tasks you want to automate. Each task consists of a series of different steps, and you can specify what each step should do, including making API calls, data manipulation, or other actions.

2. **Manage Data Flow**: You have full control over how data flows between steps. Each step can use data provided by the previous step or include its parameters, allowing you to build flexible and interconnected workflows.

3. **Handle Errors**: For robust execution, Tasker supports error handling at each step. You can specify how the application should respond if an error occurs during a step's execution.

4. **Schedule Execution**: If you want tasks to run automatically, you can schedule them using cron syntax. Define the schedule for each task, and Tasker will ensure they execute at the specified times.

## Features

Tasker showcases various aspects of web development, including:

- **Hexagonal Design Pattern**: Tasker follows a hexagonal architecture, promoting separation of concerns and making the application more maintainable and testable.

- **Adapter Integration**: It integrates with external APIs through adapters, allowing for easy extensibility and flexibility when interacting with different services.

- **Databases**: Tasker connects to MySQL and Redis databases, demonstrating how to use multiple data stores within a single application.

- **Router Library**: It uses go-chi as the router library, showcasing how to set up routing and handling HTTP requests in a Go application.

- **Testing**: Tasker utilizes stretchr/testify for writing unit and integration tests. Additionally, DATA-DOG/go-sqlmock is employed for testing interactions with the database.

## Getting Started

To get started with Tasker, follow these steps:

1. Clone the repository: `git clone https://github.com/your-username/tasker.git`
2. Navigate to the project directory: `cd tasker`
3. Start the application using Docker Compose: `docker-compose up -d`

Ensure you have Docker and Docker Compose installed on your system before running the above commands.

## Endpoints

Tasker provides the following endpoints for you to explore and interact with:

- **POST /jobs/execute-scheduled-tasks**: Execute scheduled tasks.

- **POST /schedule/**: Create a new schedule.

- **POST /task/**: Create a new task.

- **GET /task/{taskID}**: Retrieve a specific task by its ID.

- **POST /task/{taskID}/execute/{scheduleID}**: Execute a specific task associated with a schedule.


## License

Tasker is licensed under the [MIT License](LICENSE). You are free to use, modify, and distribute this project as per the terms of the license.

---

*Please note that Tasker is still a work in progress, and some features may be incomplete or under development.*
