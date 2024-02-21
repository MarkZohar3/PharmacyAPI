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

### Development environment:

Only docker-compose-dev.yml file is needed. Command:
```bash
docker-compose -f docker-compose-dev.yml up
```
### Local setup

Clone the project:
```bash
$ git clone https://gitlab.com/mark.zohar3/pharmacyapi.git
```
Open root directory (where docker-compose.yml file resides) and run command: 
```bash
$ docker -f docker-compose-local.yml up
```

## Usage
 
 Navigate to http://localhost:8080/swagger/index.html, where you can run CRUD operations with the help of swagger.

## Local Kubernetes Setup
### Prerequisites
- kubectl
- minikube

### Steps
Start minikube 

```bash
$ minikube start --memory=4096 --driver=hyperv
```

Get IP, for locating service
```bash
$ minikube ip 
```
You should receive an ip address, to which we will refer as {MINIKUBE-XYZ-IP}.

Run kubernetes scripts in k8s folder (this command assumes you are placed in root directory of the project):
```bash
$ kubectl create -f k8s
```

Some of the optional commands, to see configuration:
```bash
$ kubectl get svc
$ kubectl get pods -o wide
$ kubectl get deploys
```

In browser, open swagger to see endpoints (due to CORS, currently db and backend can't communicate with each other via browser):
```bash
$ {SOME-XYZ-IP}:30008/swagger/index.html
```

cURL examples:

Get all pharmacies
```bash
$ curl --location --request GET "http://{SOME-XYZ-IP}:30008/api/get_pharmacies"
```

Create a new farmacy:
```bash
curl --location --request POST "http://{SOME-XYZ-IP}:30008/api/create_pharmacy" --header "Content-Type: application/json" --data-raw "{\"owner\": \"Astra\",\"name\": \"Lekarna Celje\", \"address\": \"Askerceva 14, Celje, Slovenia\"}"
```


