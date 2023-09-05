# Reverse Image Search with Inception V3

Welcome to the Reverse Image Search project! This project enables you to perform reverse image searches using the Inception V3 model for image encoding. It consists of a web-based client deployed in the browser using TensorFlow.js, a Golang-based backend utilizing gRPC, and Malvus vector DB for efficient storage and cosine similarity search based on inner product. Nginx is also employed to serve images, and Docker Compose is used for easy project configuration.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Features

- Reverse image search using Inception V3 model for encoding.
- Web-based client deployed in the browser using TensorFlow.js.
- Golang-based backend with gRPC for efficient communication.
- Utilizes Malvus vector DB for image storage and cosine similarity search with inner product.
- Nginx for serving images efficiently.
- Docker Compose for simplified project setup and configuration.

## Getting Started

To get started with the Reverse Image Search project, follow the steps below.

### Prerequisites

Before you begin, ensure you have the following installed:

- Docker: [Docker Installation](https://docs.docker.com/get-docker/)
- Docker Compose: [Docker Compose Installation](https://docs.docker.com/compose/install/)
- Go Programming Language: [Golang Installation](https://golang.org/doc/install)


### Installation

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/yeaung276/ImageSearch.git
   ```
2. Nevagate to the project directory.
    ```bash
    cd ImageSearch
    ```
3. Build and start the project using docker-compose:
    ```bash
    docker-compose up
    ```
This will set up all the necessary components, including the backend, Malvus vector DB, and Nginx, and make them ready for use.

## Usage
Once the project is up and running, you can access the reverse image search application through your web browser. The web client will allow you to upload an image, and the system will search for similar images based on the features extracted by the Inception V3 model. The Golang-based backend and Vector DB handles the encoding, search, and retrieval of images efficiently.
###### http client: `localhost`
###### db interface: `localhost:8000`
![client](https://github.com/yeaung276/ImageSearch/blob/91bab8728fd696cb854ea5ed4246edb716be3fc6/docs/screenshot-web.png)
## Architecture

The Reverse Image Search project comprises the following key components:

- **Web Client (TensorFlow.js)**: This is the client-side application that runs in the user's web browser. It enables users to upload images and get encodings for reverse image searches.

- **GRPC server**: The backend is the server-side application responsible for image storage and connecting to vector DB. It leverages gRPC for efficient communication with the vector database and returns search results to the client.

- **Malvus Vector DB**: Malvus Vector DB is a specialized database designed for storing and efficiently querying high-dimensional vectors. It plays a critical role in the project, particularly for similarity searches.

- **Nginx**: Nginx serves as a web server that ensures efficient image delivery and manages incoming web traffic.

- **Docker**: This project use docker for managing dependencies components.

This architectural design enables efficient image searches and provides a robust foundation for the Reverse Image Search project.
## License

This project is licensed under the **MIT License**. Feel free to use, modify, and distribute it according to the terms of the license.
