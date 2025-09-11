# docker build -t scheduler:v001 -f ./build/scheduler/Dockerfile .

# Собираем в гошке
FROM golang:1.23.12-alpine3.22 AS build

ENV BIN_FILE /opt/calendar/scheduler-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.22

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="scheduler"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/scheduler-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/scheduler_config.yaml
COPY ./configs/scheduler_config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE}
