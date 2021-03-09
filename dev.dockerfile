FROM golang@sha256:84a409b4c174965a51e393064e46f6eb32adb80daa6097851268458136fd68b6

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o app" --command=./app
