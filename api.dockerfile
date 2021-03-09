FROM golang@sha256:84a409b4c174965a51e393064e46f6eb32adb80daa6097851268458136fd68b6 AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# set default env
ARG NAME=$NAME
ENV USER=appuser
ENV UID=10001 
ARG VERSION=$VERSION
ARG BRANCH=$BRANCH
ARG BUILD=$BUILD

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# set dependencies
WORKDIR "${GOPATH}/src/${NAME}/"
COPY . .

# Fetch dependencies.
# Using go get.
RUN go mod download

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags="-w -s -extldflags \"-static\" -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Branch=${BRANCH}" -a \
      -o "/go/bin/app" .

FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder "/go/bin/app" "/go/bin/app"

# Use an unprivileged user.
USER appuser:appuser

# Run the binary.
ENTRYPOINT ["/go/bin/app"]
