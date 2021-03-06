package util

import (
	"dist/common/log"
	"dist/pb"
	"net"

	"google.golang.org/grpc"
)

type CommunicationServer struct {
	pb.UnimplementedCommunicationServer
}

func SetupServer(f *string, address string, srv pb.CommunicationServer) {
	log.Log(f, "<SetupServer(f:\"%s\", address:\"%s\")>", *f, address)

	log.Log(f, "<SetupServer> Listening on \"%s\"", address)
	lis, err := net.Listen("tcp", address)
	log.FailOnError(f, err, "Couldn't start listening on \"%s\"", address)

	log.Log(f, "<SetupServer> Starting and registering new server...")
	server := grpc.NewServer()
	pb.RegisterCommunicationServer(server, srv)
	log.Log(f, "<SetupServer> Done registration")

	go func() {
		log.Log(f, "<SetupServer> Serving on address")
		err = server.Serve(lis)
		log.FailOnError(f, err, "Couldn't serve on address")
	}()

	log.Log(f, "<SetupServer> Done")
}

func SetupClient(f *string, target string) (*grpc.ClientConn, *pb.CommunicationClient) {
	log.Log(f, "<SetupClient(f:\"%s\", target:\"%s\")>", *f, target)

	log.Log(f, "<SetupClient> Dialing address with options: grpc.WithInsecure(), grpc.WithBlock() ...")
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.FailOnError(f, err, "Couldn't connect with grpc.Dial with address: %s", target)
	log.Log(f, "<SetupClient> Done dial")

	log.Log(f, "<SetupClient> Setting up Client...")
	client := pb.NewCommunicationClient(conn)
	log.Log(f, "<SetupClient> Done client")

	return conn, &client
}
