package dev.whiteo.yadoma.config;

import container.v1.ContainerServiceGrpc;
import image.v1.ImageServiceGrpc;
import io.grpc.ManagedChannel;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.test.util.ReflectionTestUtils;
import system.v1.SystemServiceGrpc;

import static org.junit.jupiter.api.Assertions.assertNotNull;

class GrpcConfigTest {

    private GrpcConfig grpcConfig;

    @BeforeEach
    void setUp() {
        grpcConfig = new GrpcConfig();
        ReflectionTestUtils.setField(grpcConfig, "grpcHost", "localhost");
        ReflectionTestUtils.setField(grpcConfig, "grpcPort", 50001);
    }

    @Test
    void grpcChannel_shouldCreateChannel() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        assertNotNull(channel);
        channel.shutdown();
    }

    @Test
    void containerServiceStub_shouldCreateStub() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        ContainerServiceGrpc.ContainerServiceStub stub = grpcConfig.containerServiceStub(channel);

        assertNotNull(stub);
        channel.shutdown();
    }

    @Test
    void containerServiceBlockingStub_shouldCreateStub() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        ContainerServiceGrpc.ContainerServiceBlockingStub stub = grpcConfig.containerServiceBlockingStub(channel);

        assertNotNull(stub);
        channel.shutdown();
    }

    @Test
    void systemServiceStub_shouldCreateStub() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        SystemServiceGrpc.SystemServiceStub stub = grpcConfig.systemServiceStub(channel);

        assertNotNull(stub);
        channel.shutdown();
    }

    @Test
    void systemServiceBlockingStub_shouldCreateStub() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        SystemServiceGrpc.SystemServiceBlockingStub stub = grpcConfig.systemServiceBlockingStub(channel);

        assertNotNull(stub);
        channel.shutdown();
    }

    @Test
    void imageServiceStub_shouldCreateStub() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        ImageServiceGrpc.ImageServiceStub stub = grpcConfig.imageServiceStub(channel);

        assertNotNull(stub);
        channel.shutdown();
    }

    @Test
    void imageServiceBlockingStub_shouldCreateStub() {
        ManagedChannel channel = grpcConfig.grpcChannel();

        ImageServiceGrpc.ImageServiceBlockingStub stub = grpcConfig.imageServiceBlockingStub(channel);

        assertNotNull(stub);
        channel.shutdown();
    }
}
