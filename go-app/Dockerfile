# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copia o código
COPY main.go .

# Inicializa o modulo e instala dependencias dentro do container
RUN go mod init go-app && go mod tidy

# Compila o binário estaticamente
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage (imagem mínima baseada em alpine)
FROM alpine:latest  

WORKDIR /root/

# Copia o binário construído
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
