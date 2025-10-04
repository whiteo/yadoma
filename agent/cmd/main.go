// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package main starts the Yadoma Docker agent. It connects to the Docker Engine,
// initializes gRPC services for container, image, network, volume, and system domains,
// and serves a gRPC API over TCP.
package main

import (
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	docker "github.com/whiteo/yadoma/internal/dockers"
	grpcserver "github.com/whiteo/yadoma/internal/grpc"
	"github.com/whiteo/yadoma/internal/services/container"
	"github.com/whiteo/yadoma/internal/services/image"
	"github.com/whiteo/yadoma/internal/services/network"
	"github.com/whiteo/yadoma/internal/services/system"
	"github.com/whiteo/yadoma/internal/services/volume"
	_ "github.com/whiteo/yadoma/pkg/loggers"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

func main() {
	var (
		socket = flag.String("dockers-socket",
			"/var/run/docker.sock",
			"Path to dockers engine socket",
		)
		tcpPort = flag.String("agent-tcp-port",
			":50001",
			"Run gRPC over TCP",
		)
	)

	flag.CommandLine.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = io.WriteString(w, "Yadoma - Yet Another DOcker MAnager\n")
		_, _ = io.WriteString(w, "Lightweight agent for Docker containers management over gRPC\n")
		_, _ = io.WriteString(w, "Optimized for internal networks and high performance\n\n")
		_, _ = io.WriteString(w, "Usage:\n")
		_, _ = io.WriteString(w, "  yadoma [flags]\n\n")
		_, _ = io.WriteString(w, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	log.Info().Msg("Starting Yadoma Docker Agent")

	c, err := initializeConnectToDockerEngine(*socket)
	if err != nil {
		log.Error().
			Err(err).
			Str("socket", *socket).
			Msg("Cannot connect to Docker engine. Make sure to run as root")
		return
	}
	defer func() {
		_ = c.Close()
		log.Info().Msg("Docker client connection closed")
	}()
	log.Info().Msg("Successfully connected to Docker engine")

	layer := docker.NewLayer(c)
	log.Info().Msg("Docker layer initialized")

	containerService := container.NewContainerService(layer)
	imageService := image.NewImageService(layer)
	networkService := network.NewNetworkService(layer)
	volumeService := volume.NewVolumeService(layer)
	systemService := system.NewSystemService(layer)
	log.Info().Msg("All gRPC services initialized")

	gRPC, err := startGRPCServer(*tcpPort, containerService, imageService, networkService, volumeService, systemService)
	if err != nil {
		log.Error().Err(err).Msg("failed to start gRPC server")
		return
	}
	defer func() {
		gRPC.GracefulStop()
		log.Info().Msg("gRPC server stopped gracefully")
	}()
	log.Info().Msgf("gRPC server is now serving on port %s", *tcpPort)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Info().Msg("Received stop signal, shutting down")
}

func initializeConnectToDockerEngine(socket string) (*client.Client, error) {
	c, err := client.NewClientWithOpts(client.WithHost("unix://"+socket), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func startGRPCServer(addr string,
	containerService *container.Service,
	imageService *image.Service,
	networkService *network.Service,
	volumeService *volume.Service,
	systemService *system.Service,
) (*grpc.Server, error) {
	lst, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()

	registerGrpcServers(grpcServer,
		containerService,
		imageService,
		networkService,
		volumeService,
		systemService,
	)

	go func() {
		if fErr := grpcServer.Serve(lst); fErr != nil {
			log.Error().Err(fErr).Msg("failed to initialize GRPC Server")
		}
	}()
	return grpcServer, nil
}

func registerGrpcServers(rpc *grpc.Server, services ...grpcserver.Registrator) {
	grpcserver.RegisterAll(rpc, services...)
}
