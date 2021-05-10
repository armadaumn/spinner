ARG base
FROM $base
ADD . /home/dev/ 
RUN go build -o main cmd/spinner.go
ENTRYPOINT [ "./main" ]

#FROM alpine:3.12
#RUN apk --no-cache add ca-certificates
#WORKDIR /home/dev/
#COPY --from=0 /home/dev/main .
#ENTRYPOINT [ "./main" ]