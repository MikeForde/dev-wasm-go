# Stage 1: Build the Go application
FROM registry.access.redhat.com/ubi8/go-toolset:1.18 AS builder

WORKDIR /opt/app-root/src

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the final image
FROM registry.access.redhat.com/ubi8/ubi-minimal

WORKDIR /opt/app-root/src

# Copy the built application and static files from the builder stage
COPY --from=builder /opt/app-root/src/main .
COPY --from=builder /opt/app-root/src/static ./static

EXPOSE 8080

CMD ["./main"]
