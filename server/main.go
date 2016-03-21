package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mailgun/mailgun-go"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	pb "mailgun-sender/protos"
)

var conf Config

type Email struct {
	ID      int64  `gorm:"primary_key"`
	Address string `gorm:"type:varchar(100)"`
	Status  string `gorm:"type:varchar(100)"`
}

type server struct {
	DB *gorm.DB
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func senderMail(db *gorm.DB, id int64, email string, msg string) {
	mg := mailgun.NewMailgun(conf.Domain, conf.Key, conf.PubKey)

	m := mg.NewMessage(
		conf.Sender,
		"Test Message",
		msg,
		email,
	)

	st, _, err := mg.Send(m)

	if err != nil {
		log.Fatal(err)
	}

	updaterecord := Email{}
	db.First(&updaterecord, id).Update("Status", st)
}

func (s *server) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	if !validateEmail(req.Email) {
		return &pb.SendResponse{Id: int64(0)}, nil
	}

	newrecord := Email{Address: req.Email}
	s.DB.Create(&newrecord)
	go senderMail(s.DB, newrecord.ID, req.Email, req.Message)
	return &pb.SendResponse{Id: newrecord.ID}, nil
}

func (s *server) Status(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	var email Email
	var status string

	s.DB.First(&email, req.Id)

	if email.ID != int64(0) {
		status = email.Status
	} else {
		status = "Error: ID does not exist"
	}
	return &pb.StatusResponse{Status: status}, nil
}

type Config struct {
	Port   string `yaml:"port"`
	Key    string `yaml:"mailgun_key"`
	PubKey string `yaml:"mailgun_pub_key"`
	Domain string `yaml:"mailgun_domain"`
	Sender string `yaml:"mailgun_sender"`
	TypeDB string `yaml:"type_db"`
	DBConn string `yaml:"db_connect"`
}

func (c *Config) Init() {
	file, err := os.Open("settings.yaml")
	if err != nil {
		log.Fatalf("failed to open settings file: %v", err)
	}
	defer file.Close()

	stat, _ := file.Stat()

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		log.Fatalf("failed to open settings file: %v", err)
	}

	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func main() {

	conf.Init()

	db, err := gorm.Open(conf.TypeDB, conf.DBConn)

	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	defer db.Close()

	if !db.HasTable("emails") {
		db.CreateTable(&Email{})
	}

	lis, err := net.Listen("tcp", conf.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterSenderServer(s, &server{DB: db})

	fmt.Println("Server started on port", conf.Port)
	s.Serve(lis)
}
