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