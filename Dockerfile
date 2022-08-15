# Builder Image
# ---------------------------------------------------
FROM telkomindonesia/alpine:go-<Go_Version> AS go-builder

# Set Working Directory
WORKDIR /usr/src/app

# Copy Go Source Code File
COPY . ./

# Install Go Dependencies & Compile Go File
# Set CGO_ENABLED=1 when there are some embedded C code
# Set GOOS=<platform> when needed to build for other platform
#   Ex. platform value can be ["linux", "darwin", "windows"]
RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -a -o app *.go \
    && cp app /tmp/app


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base

# Set Working Directory
WORKDIR /usr/src/app

# Copy Anything The Application Needs
COPY --from=go-builder /tmp/app ./

# Expose Application Port
EXPOSE 8080

# Run The Application
CMD ["./app"]
