## Burp

Burp is a CRUD app managing beers. Front-end is written with flutter.

See it in action here : https://burp.audrenbdb.fr/

## Details

Root project is where domain models and usecases are defined.

`cmd` package contains the main function that injects dependencies and starts the server.

Other packages are grouped around their dependencies.

By default, in-memory repository is used, even though a working implementation of PSQL repository is included in this repo.

## Tests

For usecases test doubles, I chose mix of spies and stubs.

For psql tests, I chose ory/dockertest that spins up a database container with the actual schema.

For http rest handlers I chose to do e2e tests against a fake repository.


## Dependencies

I generally stick to stdlib but I feel like it lacks convenient ways to handle parameterized routes. Therefore, I picked up chi router.

- [google/uuid](https://github.com/google/uuid)
- [go-chi/chi](https://github.com/go-chi/chi)
- [jackc/pgx](https://github.com/jackc/pgx)