# docker build -t 1aeef68f-integration-tests:develop -f ./build/test/Dockerfile .


FROM golang:1.23.12-alpine3.22

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . ${CODE_DIR}

CMD ["go", "test", "-tags=integration", "./..."]
