FROM golang:1.22.1 AS build

WORKDIR /app

COPY . .

RUN go build -o prb_care_scheduler cmd/prb_care_scheduler/main.go

FROM alpine:3.20.2

RUN apk add --no-cache tzdata

ENV TZ=Asia/Jakarta

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=build /app/fcm.json /app/fcm.json

COPY --from=build /app/prb_care_scheduler /app/prb_care_scheduler

CMD ["/app/prb_care_scheduler"]