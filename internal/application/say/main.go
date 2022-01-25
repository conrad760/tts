package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	pb "texttospeach/internal/adapters/framework/left/grpc/pb/proto"

	"google.golang.org/grpc"
)

func main() {
	backend := flag.String("b", "localhost:8080", "addree of the say backend")
	output := flag.String("o", "output.wav", "wav file where the output file will be written")
	flag.Parse()
	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", *backend, err)
	}
	defer conn.Close()
	client := pb.NewTextToSpeechClient(conn)
	text := &pb.Text{Text: flag.Arg(0)}
	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speech\"\n", os.Args[0])
		os.Exit(1)
	}

	res, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("I am not going to say %s: %v", text.Text, err)
	}
	if err := ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatalf("could not write to file %s: %v", *output, err)
	}

}
