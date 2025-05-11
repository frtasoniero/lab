package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/frtasoniero/server-streaming-server/src/pb/department"
	"google.golang.org/grpc"
)

type server struct {
	department.DepartmentServiceServer
}

func (s *server) ListPerson(req *department.ListPersonRequest, srv department.DepartmentService_ListPersonServer) error {
	file, err := os.Open("./data.csv")
	if err != nil {
		return fmt.Errorf("error on open file. error: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		time.Sleep(1 * time.Second)

		data := strings.Split(scanner.Text(), ";")
		id, _ := strconv.Atoi(data[0])
		name := data[1]
		email := data[2]
		income, _ := strconv.Atoi(data[3])
		departmentId, _ := strconv.Atoi(data[4])

		if int32(departmentId) == req.GetDepartmentId() {
			if err := srv.Send(&department.ListPersonResponse{
				Id:           int32(id),
				Name:         name,
				Email:        email,
				Income:       int32(income),
				DepartmentId: int32(departmentId),
			}); err != nil {
				return fmt.Errorf("error on send. error: %+v", err)
			}
		}
	}

	return nil
}

func main() {
	fmt.Println("Starting grpc server...")
	listener, err := net.Listen("tcp", ":5005")
	if err != nil {
		log.Fatalln("error on get listener. error: ", err)
	}

	s := grpc.NewServer()
	department.RegisterDepartmentServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalln("error on serve. error: ", err)
	}
}
