package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"

	pb "texttospeach/internal/adapters/framework/left/grpc/pb/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// serve grpc
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()
	logrus.Infof("Listening to port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrus.Fatalf("could not listen to port %d: %v", *port, err)
	}
	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	if err := s.Serve(lis); err != nil {
		logrus.Fatal("could not serve: %v", err)
	}

}

type server struct {
	pb.UnimplementedTextToSpeechServer
}

func (server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("Could not open file: %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("Could not close %s: %v", f.Name(), err)
	}

	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s", data)
	}
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file: %v", err)
	}

	return &pb.Speech{Audio: data}, nil
}
