FROM golang:1.19.6-alpine3.17

RUN apk add git
RUN mkdir /diva-challenge
WORKDIR /diva-challenge
COPY go.mod .
COPY go.sum .

ARG ACCESS_TOKEN
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN export GOPROXY="https://proxy.golang.org"
RUN go mod download
ADD . go/src/github.com/PlanningDiva/diva-challenge
COPY . .
EXPOSE 8080
RUN go install github.com/PlanningDiva/diva-challenge
ENTRYPOINT ["/go/bin/diva-challenge"]