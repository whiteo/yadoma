FROM gradle:9.2.1-jdk25 AS builder

WORKDIR /app

COPY proto /app/proto

COPY ui /app/ui

RUN apt-get update && \
    apt-get install -y curl && \
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    rm -rf /var/lib/apt/lists/*

COPY webapp /app/webapp

WORKDIR /app/webapp
RUN gradle build -x test --no-daemon

FROM eclipse-temurin:25-jre-alpine

WORKDIR /app

RUN addgroup -S yadoma && adduser -S yadoma -G yadoma

COPY --from=builder /app/webapp/build/libs/yadoma.jar /app/yadoma.jar

RUN chown -R yadoma:yadoma /app

USER yadoma

EXPOSE 8080

ENV JAVA_OPTS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0 -XX:+UseG1GC"

HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/actuator/health || exit 1

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar /app/yadoma.jar"]
