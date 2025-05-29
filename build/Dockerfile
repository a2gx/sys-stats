# Этап сборки
FROM golang:1.24 as build

ENV BIN_FILE /opt/sys-stats/daemon
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем бинарник для linux, статически, без CGO
ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "$LDFLAGS" \
    -o ${BIN_FILE} ./cmd/daemon

# Финальный образ — минимальный
FROM alpine:3.19

LABEL SERVICE="sys-stats"
LABEL MAINTAINERS="alxbuylov@gmail.com"

# Удаляем ненужные пакеты (procps не нужен)
ENV BIN_FILE /opt/sys-stats/daemon
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sys-stats/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} run -n 3 -m 7 --config ${CONFIG_FILE}
