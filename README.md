# EPP client

The EPP client is meant to connect to an
EPP server to be able to make requests and read responses.
There is currently only one interface to the EPP client
but more can be added to make it easy to use the client in different ways.

## Prompt interface

Prompt interface is meant to run on a commandline. The user will be given
options for available commands and their available data.
There is also an option to send custom xml to the epp server.
And the user can validate and print their commands before sending them.

### Run

Connecting to an EPP server on you local machine that listens to
port 700:\
`go run ./cmd/prompt/*.go --host 127.0.0.1 --port 700 --cert path.to.cert --key path.to.key`

### Options

| Parameter                   | Description                             | Default value       |
|-----------------------------|-----------------------------------------|---------------------|
| `port` or `p`               | the port to send requests to            | 7000                |
| `host` or `h`               | the host to send requests to            | 127.0.0.0           |
| `cert` or `c`               | path to the cert to use for tls         | some-cert-path.cert |
| `key` or `k`                | path to the key to use for tls          | some-key-path.key   |
| `keep-alive` or `a`         | keep connection to the epp server alive | false               |
| `validate-responses` or `v` | validate responses from epp server      | true                |

### Validation

If validation of responses is active the result will
be printed under the response itself. Either if any errors
were found the output from libxml2 will be printed or "ok" if no
errors were found.
