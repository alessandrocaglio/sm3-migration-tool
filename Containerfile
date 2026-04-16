# Stage 1: Build the frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /build/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go binary
FROM golang:1.26-alpine AS go-builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy built frontend assets into the Go package directory for embedding
COPY --from=frontend-builder /build/frontend/dist/ pkg/api/ui/
RUN CGO_ENABLED=0 GOOS=linux go build -o sm3-migration-tool main.go

# Stage 3: Operational Image
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
LABEL maintainer="Alessandro Caglio"
LABEL name="sm3-migration-tool"
LABEL summary="Red Hat OpenShift Service Mesh 2 to 3 Migration Assistant"

WORKDIR /app
COPY --from=go-builder /build/sm3-migration-tool .
# Ensure standard mock data is available for test cases
COPY testdata/mock-cluster.yaml testdata/mock-cluster.yaml

# Set permissions for OpenShift (arbitrary UID support)
RUN chown -R 1001:0 /app && chmod -R g+w /app
USER 1001

EXPOSE 8080

ENTRYPOINT ["/app/sm3-migration-tool"]
CMD ["serve", "--mode", "live", "--port", "8080"]
