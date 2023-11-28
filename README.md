# Movie API

API to manage films and user's favourites list using Golang.

The API connects to an MySQL database using GORM (ORM library for Golang) and the project can be deployed using docker compose.

## Features

- View, create, delete and edit movies
- User authentication (login and signup) using JWT tokens
- Add movies to favourite and manage user's favourite lists
- Configurable using .env file
- Easy deployment using docker compose
- Documentation using Swagger 

## Deployment

The API uses enviroment variables for configuration and docker compose for easy deployment.

1. Copy `.env.example` file and rename to `.env`. Modify the configuration as needed (port, database credentials, etc)

2. Run docker compose file

```
docker compose up -d
```
3. Test the API. Use `docs/swagger.yaml` documentation for help.

*Database will automatically be populated with sample data and users when migrating database schema first time.*

**NOTE:** To check server logs while using the API:
```
docker logs movies-api
```

**NOTE:** To shutdown resources and volumes:

```
docker compose down --rmi all -v
```

## F.A.Q

#### Why use `httprouter` for HTTP API routing instead of another library or the standard library?

Although the standard library provides enough funcionalities to develop a basic API with routing, middlewares, etc, I missed method based routing and I think I could use an HTTP library for this.

For this API, I decided to use `httprouter` as it is a very lightweight package and it is one of the fastest HTTP routing packages. As this API is not really complex at the moment, I think it is the best choice for its simplicity and performance.


## Future improvements

- [ ] Use HTTPS in the API (TLS).
- [ ] Add `/top` endpoint to list most favourited movies.
- [ ] Create `directors`, `actors` and `genre` database tables to allow more complex relations and queries.

