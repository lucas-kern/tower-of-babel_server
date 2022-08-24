# tower-of-babel_server

## Development

### Code Organization
To start developing a `Go` project the `go mod init` command needs to be run with a domain of choice. This project is using `github.com/lucas-kern/tower-of-babel_server` as the GitHub repository so has been initialized with that domain.

In order to import a package into another package for use it will need to be in the format of

`import github.com/lucas-kern/tower-of-babel_server/{PATH}/{TO}{FILE}`

With the {PATH}/{TO}{FILE} being the directory structure of the imported package.

To add other packages from GitHub that are not in this repository use `go get {PACKAGE-DOMAIN}`

If needed can run `go mod tidy` to update current packages and remove unused ones.

Read [`How to Write Go Code`](https://go.dev/doc/code) for a simple understanding of how to organize the code

### Starting the project

To start the project run `go run app/main.go`

Stop the server by hitting `ctrl + c`