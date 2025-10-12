package dev.whiteo.yadoma.config;

import container.v1.ContainerServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Configuration class for setting up gRPC clients for container service.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Configuration
public class GrpcConfig {

    /**
     * Creates a gRPC channel for communication with the container service.
     *
     * @return ManagedChannel configured for localhost:50051 with plaintext encryption.
     */
    @Bean
    public ManagedChannel grpcChannel() {
        return ManagedChannelBuilder.forAddress("localhost", 50051)
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
}