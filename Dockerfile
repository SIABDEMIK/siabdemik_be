# Menggunakan base image resmi Go
FROM golang:1.20-alpine

# Set environment variable
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Salin file go.mod dan go.sum lalu download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode aplikasi ke dalam container
COPY . .

# Build aplikasi
RUN go build -o server ./cmd/myapp

# Ekspos port yang akan digunakan
EXPOSE 8080

# Perintah untuk menjalankan aplikasi
CMD ["./server"]
