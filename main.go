package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/sana/proto"
	"google.golang.org/grpc"
)

type sanaService struct{}

func (s *sanaService) StreamNetworks(stream proto.CxMateService_StreamNetworksServer) error {
	workDir, err := ioutil.TempDir("/tmp", "sana")
	if err != nil {
		log.Fatalf("failed to create temp directory for SANA networks: %v", err)
	}
	originDir, err := os.Getwd()
	if err != nil {
		log.Printf("could not get current working directory: %v", err)
		return err
	}
	if err := os.Chdir(workDir); err != nil {
		log.Printf("could not change working directory: %v", err)
		return err
	}
	tmpSrcFile, err := ioutil.TempFile(workDir, "source_net_two")
	if err != nil {
		log.Printf("could not create temp source network file: %v", err)
		return err
	}
	tmpTgtFile, err := ioutil.TempFile(workDir, "source_net_one")
	if err != nil {
		log.Printf("could not create temp target network file %v", err)
		return err
	}
	params := []*proto.Parameter{}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println("input stream closed")
			break
		}
		if err != nil {
			log.Printf("err in receive stream: %v", err)
			return err
		}
		switch e := in.GetElement().(type) {
		case *proto.NetworkElement_Parameter:
			params = append(params, e.Parameter)
		case *proto.NetworkElement_Edge:
			srcNode := e.Edge.GetSourceId()
			tgtNode := e.Edge.GetTargetId()
			edgeLine := fmt.Sprintf("%d %d\n", srcNode, tgtNode)
			if in.Label == "network one" {
				tmpSrcFile.WriteString(edgeLine)
			} else {
				tmpTgtFile.WriteString(edgeLine)
			}
		}
	}
	tmpOutFile, err := ioutil.TempFile(workDir, "out_net")
	if err != nil {
		log.Println("could not create temp out network file")
		return err
	}
	cmd := exec.Command("sana", "-fg1", tmpSrcFile.Name(), "-fg2", tmpTgtFile.Name(), "-t", "1", "-o", tmpOutFile.Name())
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("running sana")
	if err := cmd.Run(); err != nil {
		log.Printf("running sana failed: %v", err)
		return err
	}
	log.Println("sana finished")
	var edgeID int64
	outFile, err := os.Open(fmt.Sprint(tmpOutFile.Name(), ".align"))
	if err != nil {
		log.Printf("opening the temp out network file failed: %v", err)
		return err
	}
	fScanner := bufio.NewScanner(outFile)
	for fScanner.Scan() {
		var srcNode, tgtNode int64
		edgeLine := fScanner.Text()
		fmt.Sscan(edgeLine, &srcNode, &tgtNode)
		stream.Send(&proto.NetworkElement{
			Label: "alignment",
			Element: &proto.NetworkElement_Edge{
				Edge: &proto.Edge{
					Id:       edgeID,
					SourceId: srcNode,
					TargetId: tgtNode,
				},
			},
		})
		edgeID++
	}
	regex := `^\a+: [\d.]+$`
	if err := os.Chdir(originDir); err != nil {
		log.Printf("could not change working directory: %v", err)
		return err
	}
	if err := os.RemoveAll(workDir); err != nil {
		log.Printf("could not delete SANA work directory: %v", err)
		return err
	}
	log.Println("requested completed")
	return nil
}

func main() {
	log.Println("starting sana service wrapper")
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("failed to listen for connections: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterCxMateServiceServer(s, &sanaService{})
	log.Println("starting sana grpc server")
	s.Serve(l)
}
