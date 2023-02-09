FROM golang:1.18.2-alpine

# Install images tools required for project
RUN apk add --no-cache git && apk add --update bash && apk add screen

ARG APP_ENV
ARG COMMAND_TYPE
ARG CONSUMER_TOPIC_NAME
ARG CONSUMER_GROUP

ENV COMMAND_TYPE=$COMMAND_TYPE
ENV CONSUMER_TOPIC_NAME=$CONSUMER_TOPIC_NAME
ENV CONSUMER_GROUP=$CONSUMER_GROUP

# Copy all solution files
WORKDIR /order-service

COPY go.mod ./
COPY go.sum ./

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code
COPY . ./

# Install library dependencies
RUN export GO111MODULE=on

RUN mkdir -pv /src/app/ssl

RUN cd app && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/app/ordersrv .

COPY app/.env.$APP_ENV /src/app/.env
COPY app/ssl/ /src/app/ssl

WORKDIR /src/app
#RUN rm -r /order-service

RUN chmod +x ./ordersrv

EXPOSE 5000
ENTRYPOINT ["/src/app/ordersrv"]
CMD ["ordersrv"]