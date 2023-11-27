# Movie API

API to manage films and user's favourites list using Golang.

The API connects to an MySQL database using GORM (ORM library for Golang) and the project can be deployed using docker compose.

## Features

- [x] View, create, delete and edit movies
- [x] User authentication (login and signup) using JWT tokens
- [x] Add movies to favourite and manage user's favourite lists

## Deployment

The API uses enviroment variables for configuration and docker compose for easy deployment.

1. Copy `.env.example` file and rename to `.env`. Modify the configuration as needed (port, database credentials, etc)

2. Run docker compose file

```
docker compose up -d
```
3. Test the API !

To shutdown resources and volumes:

```
docker compose down --rmi all -v
```

## F.A.Q

#### Why use `httprouter` for HTTP API routing instead of another library or the standard library?

Although the standard library provides enough funcionalities to develop a basic API with routing, middlewares, etc, I missed method based routing and I think I could use an HTTP library for this.

For this API, I decided to use `httprouter` as it is a very lightweight package and it is one of the fastest HTTP routing packages. As this API is not really complex at the moment, I think it is the best choice for its simplicity and performance.


## Future improvements

- [ ] Use HTTPS in the API (TLS)


