# Builder Images
# ---------------------------------------------------
FROM telkomindonesia/alpine:go-1.17 AS go-builder

# Set Working Directory
WORKDIR /app

# Copy Go Source Code File
COPY . ./

# Install Go Dependencies & Compile Go File
RUN go build -o main main.go

# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base

# Set Working Directory
WORKDIR /app

# Copy Anything The Application Needs
COPY --from=go-builder /app/main ./

COPY app.env .

# Expose Application Port
EXPOSE 8080

# Run The Application
CMD ["/app/main"]
