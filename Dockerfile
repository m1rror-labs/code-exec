FROM backpackapp/build:v0.31.0

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
    mold \
    unzip


# Install TypeScript globally
RUN npm install -g typescript

RUN curl -fsSL https://bun.sh/install | bash

# Copy all files into /app folder
WORKDIR /app
COPY . .

# Copy the Python script into the image
COPY decode_base58.py /app/decode_base58.py

# Install Python and dependencies
RUN apt-get update && apt-get install -y python3 python3-pip && \
    pip3 install base58

# Run the Python script as part of the build step
RUN python3 /app/decode_base58.py

# Download all dependencies and build project
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go

WORKDIR /app/pkg/dependencies/typescript
RUN npm install

WORKDIR /app/pkg/dependencies/rust
RUN cargo fetch
RUN cargo build --bin main

WORKDIR /app/pkg/dependencies/anchor/code-exec
RUN anchor build

RUN cargo install sccache
RUN export RUSTC_WRAPPER=sccache

WORKDIR /app






# Run the application
CMD ["./main"]