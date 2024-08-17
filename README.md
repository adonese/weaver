# Exinity payment gateway

This program simulates a payment gateway solution. In payment terms we call those Escrow Integration, or Service Providers Integrations. This system processes deposits and withdrawals requests, while providing a callback to update a transaction status. One simple scenario for this is an Exinity user who wants to "withdraw" some of their funds from their exinity wallet, to their bank account. In this case, we can call the withdraw api and that will initiate the transaction, and we provide a webhook (this callback) to the bank such that when the transaction is completed in their end, they will trigger this webhook (callback), for us to update the transaction status. 

## How to run the program
The program uses generated file(s) from weaver namely `weaver_gen.go` 
- Make sure that go is installed on your machine
- Clone the repo, `git clone https://github.com/adonese/exinity`
- cd into the project directory
- download depenencies `go get`
- run `go build -o payment_system`
- run `./payment_system`

### optional steps
- install weaver cli tool, `go install github.com/ServiceWeaver/weaver/cmd/weaver@latest`, to generate weaver files
- run `weaver generate`


## What is weaver?

Service Weaver is a new framework for writing microservices (modular monolith) in Go. It is very interesting in that, it has the benefits of the both worlds:
- for monolith, you still write your application the same way you would write a typical monolith application: a single binary that is `go build`-able
- that same single binary runs as separate microservices, managed by the service weaver runtime.
So essenatially, one gets the benefits of ease of develpment of a monolith, and the benefits of a microservice architecture.

`weaver.toml` is the configuration file for the service weaver runtime, it is a very compact toml with only `binary` as required field.

## Tests 

The program is composed of an integration test that tests the full flow of a transaction from initiation to completion. You can run the tests via `go test -v ./...`

## Design decisions

- we decided to keep the backing storage of the transactions in memory, for the sake of simplicity and shared for both gateways
- we have http handlers for withdraw and deposit, and a callback to update the transaction status. We also have a transaction status api
- http handler will call 