Go HTTP Proxy Server
====================

A simple HTTP proxy server written in Go with the following features:

-   Parsing a configuration file in YAML format: The proxy server reads a configuration file in YAML format to determine its behavior. The configuration file includes fields such as the port to bind to, a map of headers to add or modify in the forwarded request, and the target host and port to forward requests to.

-   Validation of the configuration file using the `validator` package: The proxy server uses the `validator` package to validate the fields in the configuration file. If any errors are found, the proxy server will print a list of the errors and exit.

-   An option to generate a default configuration file using the `-g` flag: The proxy server can generate a default configuration file named `config.yaml` by using the `-g` flag when starting the server.

-   An option to run in verbose mode using the `-v` flag: The proxy server can be run in verbose mode by using the `-v` flag. In verbose mode, the server will print additional information about its behavior to the console.

-   An option to specify a custom configuration file using the `-c` flag: The proxy server can use a custom configuration file by specifying the path to the file using the `-c` flag. The default configuration file is `config.yaml`.

-   Reading the contents of an incoming HTTP request and forwarding it to a target host and port specified in the configuration file: When the proxy server receives an incoming HTTP request, it reads the contents of the request and forwards it to the target host and port specified in the configuration file.

-   Modifying the headers of the forwarded request based on the `header_map` field in the configuration file: The proxy server can modify the headers of the forwarded request based on the `header_map` field in the configuration file. This field specifies a mapping of header names and values to add or modify in the forwarded request.

-   Returning the response from the target host to the client: The proxy server receives the response from the target host and returns it to the client that made the original request.

Components
----------

The code consists of the following main components:

-   The `Config` struct: This struct represents the configuration for the proxy server. It includes fields for the port to bind to, a map of headers to add or modify in the forwarded request, and the target host and port to forward requests to. The `Config` struct also includes validation tags for the fields, which are used by the `validator` package to perform validation.

-   The `Parse` function: This function reads and parses the configuration file, performs validation using the `validator` package, and returns a pointer to a `Config` struct. If any errors are found during validation, the `Parse` function will print a list of the errors and exit the program.

-   The `main` function: This is the entry point of the program. It handles command line flags, reads the configuration file, and starts an HTTP server to handle incoming requests.

Usage
-----

To start the proxy server, run the following command:

Copy code

```bash
go run main.go
```

Use the following flags to customize the behavior of the proxy server:

-   `-g`: Generate a default configuration file named `config.yaml`
-   `-v`: Run in verbose mode
-   `-c`: Specify a custom configuration file. The default is `config.yaml`.

For example, to generate a default configuration file and run the proxy server in verbose mode, use the following command:

Copy code

```bash
go run main.go -g -v
```

To specify a custom configuration file, use the `-c` flag followed by the path to the configuration file. For example:

Copy code

```bash
go run main.go -c /path/to/custom/config.yaml
```

