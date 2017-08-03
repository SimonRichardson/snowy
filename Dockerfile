FROM golang:1.8

EXPOSE 7650

COPY ./ /go/src/github.com/trussle/snowy
WORKDIR /go/src/github.com/trussle/snowy

ARG mode
ENV MODE=${mode}

RUN make
RUN chmod +x ./dist/documents

CMD ["sh", "-c", "./dist/documents ${MODE}"]