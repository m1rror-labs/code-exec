FROM backpackapp/build:v0.30.1

RUN rustup toolchain uninstall stable && rustup toolchain install stable

RUN apt-get update && apt-get install -y wget curl
RUN wget https://golang.org/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Node.js and npm
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libudev-dev \
    llvm \
    libclang-dev \
    protobuf-compiler \
    libssl-dev \
    libc6 \
    ghdl \
    openssl \
    clang-15 \
    mold


# Install TypeScript globally
RUN npm install -g typescript

# Copy all files into /app folder
WORKDIR /app
COPY . .

# Download all dependencies and build project
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go

WORKDIR /app/pkg/dependencies/runtimes/typescript
RUN npm install

WORKDIR /app/pkg/dependencies/runtimes/rust
RUN cargo fetch
RUN cargo build --bin main

RUN cargo install sccache
RUN export RUSTC_WRAPPER=sccache

WORKDIR /app

# Run the application
CMD ["./main"]