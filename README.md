# ASCII Canvas

A simple web server to perform drawing operations on a text canvas.

This is intended for educational purposes.

## Build Requirements
The project was built using `go 1.17`.

In case of code modification, the following projects are used for code generation and as such, must be installed. 

They are executed automatically on `go generate`.

| Name    |  Link                              | Generates
|---------|------------------------------------|---------------
| swag    |  https://github.com/swaggo/swag    | Swagger documentation for http endpoints
| mockery |  https://github.com/vektra/mockery | Mock implementations for testing

## Building

It is recommended to use `make` to build the project

| Command      | Description |
|--------------|-------------|
| `make build` | Build the binary into the `build` subfolder
| `make clean` | Clean up previous builds

## Running

Build the binary first and then execute:

```
cd build
./ascanvas
```

You need to have the `build` folder as writable and port `1337` opened.

You may customize the database and listen address by creating a file `ascanvas.json` in the same directory as the binary (`build` folder).

A sample of this file is included as `ascanvas.json.dist`

## Using

The server can be accessed via web interface:

| URL                            | Description |
|--------------------------------|-----------------
| http://127.0.0.1:1337/swagger  | View API endpoints and perform requests using Swagger UI           
| http://127.0.0.1:1337/         | View listing of canvas items and access **live update UI**    

## License

This project is provided under the MIT license. A copy of the license found in this repository.
