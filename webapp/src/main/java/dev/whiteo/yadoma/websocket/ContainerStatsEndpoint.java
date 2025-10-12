package dev.whiteo.yadoma.websocket;

import container.v1.Container;
import container.v1.ContainerServiceGrpc;
import dev.whiteo.yadoma.config.SpringWebSocketConfigurator;
import dev.whiteo.yadoma.security.TokenInteract;
import dev.whiteo.yadoma.service.UserService;
import io.grpc.stub.StreamObserver;
import jakarta.websocket.CloseReason;
import jakarta.websocket.OnClose;
import jakarta.websocket.OnOpen;
import jakarta.websocket.Session;
import jakarta.websocket.server.PathParam;
import jakarta.websocket.server.ServerEndpoint;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * WebSocket endpoint for container statistics.
 * Handles WebSocket connections for container statistics and manages statistic streaming.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Component
@RequiredArgsConstructor
@SuppressWarnings("unused")
@ServerEndpoint(value = "/ws/containers/{id}/stats", configurator = SpringWebSocketConfigurator.class)
public class ContainerStatsEndpoint {

    private final ContainerServiceGrpc.ContainerServiceStub containerStub;
    private final TokenInteract tokenInteract;
    private final UserService userService;

    private ExecutorService executor;

    /**
     * Handles WebSocket connection for container statistics.
     *
     * @param session the WebSocket session
     * @param id the ID of the container to stream statistics from
     */
    @OnOpen
    @SuppressWarnings("unused")
    public void onOpen(Session session, @PathParam("id") String id) {
        try {
            List<String> tokens = session.getRequestParameterMap().get("token");
            String token = (tokens != null && !tokens.isEmpty()) ? tokens.getFirst() : null;

            if (token == null || !tokenInteract.validateToken(token)) {
                session.close(new CloseReason(CloseReason.CloseCodes.VIOLATED_POLICY, "Unauthorized"));
                return;
            }

            String userId = tokenInteract.getUserId(token);
            if (!userService.isContainerIdContains(id, userId)) {
                session.close(new CloseReason(CloseReason.CloseCodes.VIOLATED_POLICY, "Access denied"));
                return;
            }

            executor = Executors.newSingleThreadExecutor();
            executor.submit(() -> streamStats(session, id));

        } catch (Exception _) {
            closeSession(session);
        }
    }

    /**
     * Closes the WebSocket session and shutdown the executor service.
     *
     * @param session the WebSocket session to close
     */
    @OnClose
    @SuppressWarnings("unused")
    public void onClose(Session session) {
        closeSession(session);
    }

    /**
     * Streams container statistics to the WebSocket session.
     *
     * @param session the WebSocket session
     * @param containerId the ID of the container to stream statistics from
     */
    private void streamStats(Session session, String containerId) {
        Container.GetContainerStatsRequest request = Container.GetContainerStatsRequest.newBuilder()
                .setStream(true)
                .setId(containerId)
                .build();

        containerStub.getContainerStats(request, new StreamObserver<>() {
            @Override
            public void onNext(Container.GetContainerStatsResponse value) {
                try {
                    String json = """
                            {
                              "cpu": %d,
                              "memUsage": %d,
                              "memLimit": %d,
                              "netInput": %d,
                              "netOutput": %d
                            }
                            """.formatted(value.getCpuUsage(),
                            value.getMemUsage(), value.getMemLimit(), value.getNetInput(), value.getNetOutput());

                    session.getBasicRemote().sendText(json.trim());
                } catch (IOException _) {
                    closeSession(session);
                }
            }

            @Override
            public void onError(Throwable t) {
                try {
                    session.getBasicRemote().sendText("[error] " + t.getMessage());
                } catch (IOException _) {
                }
                closeSession(session);
            }

            @Override
            public void onCompleted() {
                closeSession(session);
            }
        });
    }

    /**
     * Closes the WebSocket session and shutdown the executor service.
     *
     * @param session the WebSocket session to close
     */
    private void closeSession(Session session) {
        try {
            session.close();
        } catch (IOException _) {
        }
        if (executor != null && !executor.isShutdown()) {
            executor.shutdownNow();
        }
    }
}