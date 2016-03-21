package main

import (
	"log"
	"flag"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
    pb "mailgun-sender/protos"
)

const (
	address     = "localhost:50051"
	defaultNum = 0
)

var(
    command = flag.String("c", "None", "Command: \n\tsend - for sending new email\n\tstatus - for check the status ")
    id = flag.Int64("id", int64(0), "ID of sending email")
    email = flag.String("email", "None", "Email address")
    message = flag.String("msg", "None", "Message of sending Email")
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSenderClient(conn)

	flag.Parse()

	switch *command{
	    case "send":
	        r, err := c.Send(context.Background(), &pb.SendRequest{Email: *email, Message: *message})
        	if err != nil {
        		log.Fatalf("could not greet: %v", err)
        	}
        	log.Println("ID: ", r.Id)
        case "status":
            r, err := c.Status(context.Background(), &pb.StatusRequest{Id: *id})
        	if err != nil {
        		log.Fatalf("could not greet: %v", err)
        	}
        	log.Println("Status: ", r.Status)
	}
}