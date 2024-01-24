# Pharmacy API

This is a sample project in Golang, with basic CRUD operations.


## Overview

This is a sample project in Golang, with basic CRUD operations. It is integrated with PostGreSQL database and features swagger.

## Features

- Create, Delete, Get and Update Pharmacy
- Containerisation
- Swagger

## Prerequisites

- Docker Desktop 

## Getting Started

Clone the project:
```bash
$ git clone https://gitlab.com/mark.zohar3/pharmacyapi.git
```
Open root directory (where docker-compose.yml file resides) and run command: 
```bash
$ docker compose build
```

## Usage

From root directory of project, run command:

```bash
$ docker compose up
```
 
 Navigate to http://localhost:8080/swagger/index.html, where you can run CRUD operations with the help of swagger.