## Burp - clean architecture app

Burp is a CRUD app managing beers. Front-end is written with flutter.

See it in action here : https://burp.audrenbdb.fr/

## Details

Root project is where models are defined and usecases containing business logic are exposed. There is little to no business logic in that example project.

Cmd folder contains starts our server and injects dependencies.

By default, in-memory repository is used, even though a working implementation of PSQL repository is included in this repo.

## Tests

To test our external dependencies, a few mocks are generated from https://github.com/golang/mock with `go generate`.

To test our usecases, in-memory implementation of repository is used.