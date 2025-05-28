## About
Senior Software Engineer Task

This repo contains a lightweight headless web application that includes basic authentication and data retrieval.

### Design choices

- We are using Postgres to store out data for its robustness and small footprint.
- We are also using it to implement the rate limiting logic to keep the number of dependencies low instead of using an in memory solution perhaps.
- We used [sqlc](https://sqlc.dev/) and [goose](https://github.com/pressly/goose) to manage out database schema.
- To hash the passwords before storing them we are using the *Argon2* algorithm with a 64MB memory cost for the strong security and the memory-hardness to brute-force attacks that it offers.
- As an authentication mechanism we are issuing JWT tokens signed with a static key for simplicity and skipping the usage of public and private keys.

### Future improvements

- Gracefully shutdown the server
- Use a public/private key to sign the JWTs

### Local Development

To run the application locally, you can use Docker compose:

```bash
docker compose up -d
```

Then you can start making requests to the server:

```bash
curl localhost:8080/api/health
```

### Coverage report summary 
```
$> go test -coverprofile=coverage.out ./...
        github.com/cv711/odin-takehome/server           coverage: 0.0% of statements
ok      github.com/cv711/odin-takehome/server/api       0.551s  coverage: 63.6% of statements
        github.com/cv711/odin-takehome/server/db                coverage: 0.0% of statements
ok      github.com/cv711/odin-takehome/server/internal  0.377s  coverage: 83.3% of statements
```

for more info check [coverage.html](./coverage.html)