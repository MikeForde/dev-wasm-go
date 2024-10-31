# Stage 1: Build the Go application and WASM module
FROM registry.access.redhat.com/ubi8/go-toolset:1.18 AS builder

WORKDIR /opt/app-root/src

# Copy the source code
COPY . .

# Build the Go application binary
RUN go build -o main .

# Compile the WASM module using the standard Go compiler
RUN GOOS=js GOARCH=wasm go build -o static/main.wasm ./frontend

# Copy wasm_exec.js for WebAssembly runtime support
RUN cp $(go env GOROOT)/misc/wasm/wasm_exec.js static/

# Stage 2: Create the final image
FROM registry.access.redhat.com/ubi8/ubi-minimal

WORKDIR /opt/app-root/src

# Copy the built application and static files from the builder stage
COPY --from=builder /opt/app-root/src/main .
COPY --from=builder /opt/app-root/src/static ./static

EXPOSE 8080

CMD ["./main"]
