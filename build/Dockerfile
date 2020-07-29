ARG base
FROM $base
ADD . /home/dev/ 
RUN go build -o main cmd/spinner.go
ENTRYPOINT [ "./main" ]