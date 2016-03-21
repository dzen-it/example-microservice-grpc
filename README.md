# example-microservice-grpc

**Microservice on golang.** 

**Used libraries:**

 - [grpc](https://github.com/grpc/grpc-go) 
 - [gorm](https://github.com/jinzhu/gorm)
 - [yaml](https://github.com/go-yaml/yaml)
 - [mailgun-go](https://github.com/mailgun/mailgun-go)

The server receives the messages and records the ID in the database. Next send a message through the API [Mailgun](www.mailgun.com), records the status in the database.
Using the method `status` and request  `ID`,  response received status messages from the database

`client-test` is designed to test server.

For configuration use the `settings.yaml`

For run `client-test`, use the argument of `help`
