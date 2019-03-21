# Bank of Oracle API

## To Build:
- [] Create a `.env` file in your project root. Add your database environment variables like so:
```
ATP_DEV_USERNAME=
ATP_DEV_PASSWORD=
ATP_DEV_NAME=

ATP_TEST_USERNAME=
ATP_TEST_PASSWORD=
ATP_TEST_NAME=

ATP_PROD_USERNAME=
ATP_PROD_PASSWORD=
ATP_PROD_NAME=
```
- [] Create a `wallet/` directory in your project root. Unzip your Oracle Autonomous Transaction Processing credential wallet into it.
- [] Run `$ docker build .`
- [] Run `$ docker run <your container id>`

## To Test:
- Run `$ go test -v`
