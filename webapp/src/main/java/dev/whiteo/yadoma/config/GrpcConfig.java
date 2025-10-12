package dev.whiteo.yadoma.config;

import container.v1.ContainerServiceGrpc;
import image.v1.ImageServiceGrpc;
import org.springframework.beans.factory.annotation.Value;
import system.v1.SystemServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.concurrent.TimeUnit;

/**
 * Configuration class for setting up gRPC clients for container, image, and system services.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Configuration
public class GrpcConfig {

    @Value("${grpc.host}")
    private String grpcHost;

    @Value("${grpc.port:50001}")
    private int grpcPort;

    /**
     * Creates a gRPC channel for communication with the container service.
     *
     * @return ManagedChannel configured for grpcHost:grpcPort with plaintext encryption and timeouts.
     */
    @Bean
    public ManagedChannel grpcChannel() {
        return ManagedChannelBuilder.forAddress(grpcHost, grpcPort)
                .usePlaintext()
                .keepAliveTime(2, TimeUnit.MINUTES)
                .keepAliveTimeout(20, TimeUnit.SECONDS)
                .keepAliveWithoutCalls(false)
                .idleTimeout(10, TimeUnit.MINUTES)
                .maxInboundMessageSize(16 * 1024 * 1024)
                .enableRetry()
                .maxRetryAttempts(3)
                .build();
    }

    /**
     * Creates a gRPC stub for asynchronous communication with the container service.
     *
     * @param channel ManagedChannel for gRPC communication.
     * @return ContainerServiceGrpc.ContainerServiceStub for asynchronous operations.
     */
    @Bean
    public ContainerServiceGrpc.ContainerServiceStub containerServiceStub(ManagedChannel channel) {
        return ContainerServiceGrpc.newStub(channel);
    }

    /**
     * Creates a gRPC stub for synchronous communication with the container service.
     *
     * @param channel ManagedChannel for gRPC communication.
     * @return ContainerServiceGrpc.ContainerServiceBlockingStub for synchronous operations.
     */
    @Bean
    public ContainerServiceGrpc.ContainerServiceBlockingStub containerServiceBlockingStub(ManagedChannel channel) {
        return ContainerServiceGrpc.newBlockingStub(channel);
    }

    /**
     * Creates a gRPC stub for asynchronous communication with the system service.
     *
     * @param channel ManagedChannel for gRPC communication.
     * @return SystemServiceGrpc.SystemServiceStub for asynchronous operations.
     */
    @Bean
    public SystemServiceGrpc.SystemServiceStub systemServiceStub(ManagedChannel channel) {
        return SystemServiceGrpc.newStub(channel);
    }

    /**
     * Creates a gRPC stub for synchronous communication with the system service.
     *
     * @param channel ManagedChannel for gRPC communication.
     * @return SystemServiceGrpc.SystemServiceBlockingStub for synchronous operations.
     */
    @Bean
    public SystemServiceGrpc.SystemServiceBlockingStub systemServiceBlockingStub(ManagedChannel channel) {
        return SystemServiceGrpc.newBlockingStub(channel);
    }

    /**
     * Creates a gRPC stub for asynchronous communication with the image service.
     *
     * @param channel ManagedChannel for gRPC communication.
     * @return ImageServiceGrpc.ImageServiceStub for asynchronous operations.
     */
    @Bean
    public ImageServiceGrpc.ImageServiceStub imageServiceStub(ManagedChannel channel) {
        return ImageServiceGrpc.newStub(channel);
    }

    /**
     * Creates a gRPC stub for synchronous communication with the image service.
     *
     * @param channel ManagedChannel for gRPC communication.
     * @return ImageServiceGrpc.ImageServiceBlockingStub for synchronous operations.
     */
    @Bean
    public ImageServiceGrpc.ImageServiceBlockingStub imageServiceBlockingStub(ManagedChannel channel) {
        return ImageServiceGrpc.newBlockingStub(channel);
    }
}