package dev.whiteo.yadoma.config;

import container.v1.ContainerServiceGrpc;
import org.springframework.beans.factory.annotation.Value;
import system.v1.SystemServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Configuration class for setting up gRPC clients for container and system services.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Configuration
public class GrpcConfig {

    @Value("${grpc.host}")
    private String grpcHost;

    /**
     * Creates a gRPC channel for communication with the container service.
     *
     * @return ManagedChannel configured for grpcHost:50051 with plaintext encryption.
     */
    @Bean
    public ManagedChannel grpcChannel() {
        return ManagedChannelBuilder.forAddress(grpcHost, 50051)
                .usePlaintext()
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
}