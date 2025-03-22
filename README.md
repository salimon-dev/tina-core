# Nexus

salimon nexus service is the gate that connects users to entities and manages profiles, authorization, credits and balancing

## Environment variables

a sample of all required environment variables are in `.env.sample` make sure you have them set as envrionment variable while service is running.
if you provide a `.env` same as `.env.example`, while using docker compose to bootstrap the service, it will take values from there.

```

PGSQL_HOST= pgsql database port
PGSQL_PORT= pgql port
PGSQL_DBNAME= nexus database name in pgsql
PGSQL_USERNAME= pgsql username
PGSQL_PASSWORD= pgsql user password

SMTP_ENDPOINT= webmail smtp tcp endpoint
SMTP_USERNAME= webmail username
SMTP_PASSWORD= webmail password
SMTP_FROM= sender email address of nexus
SMTP_PORT= webmail port

JWT_SECRET= random secret token for jwt generation

ENV= dev or production

```

## Building project locally

to build and run the service locally you need latest golang installed on your local machine. then run:

```shell
go build -o bootstrap .
```

and then:

```shell
./bootstrap
```

## Running with docker

make sure you have docker installed and running on host machine and you have all required envrionment variables `.env` then run:

```shell
docker compose up
```
