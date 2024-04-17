# testDashboardTask
# Project Testing Documentation

## Overview

This project includes automated tests for validating the functionality of the "Run report" button and monitoring the statistics of transitions through referral links using automated web tests.

## Technologies

- **Go**: The project is written in Go.
- **Testify**: Utilized for assertions and test suite structuring.
- **ChromeDP**: Used to drive browser-based tests, simulating user interactions with the web interface.

## Getting Started

### Prerequisites

- Go (version 1.15 or later recommended)
- Chrome or Chromium browser installed
- Appropriate ChromeDriver for your Chrome version

### Installing

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/ansirenko/testDashboardTask.git
cd testDashboardTask
go mod download
go mod tidy
go test -v
```

## Running Tests in Docker

Running the tests within a Docker container can help ensure a consistent testing environment. This section covers the steps to build the Docker image and run tests using Docker.

### Building the Docker Image

First, you need to build the Docker image which includes all the necessary dependencies installed. Make sure you are in the project root directory and run the following command to build the Docker image:

```bash
docker build -t testdashboardtask .
```

### Running Tests in the Docker Container

Once the Docker image is built, you can run the tests in the Docker container. Use the following command to run the tests:

```bash
docker run --rm testdashboardtask
```
