FROM iron/go

EXPOSE 8080

WORKDIR /app
ADD documents /app/

ENTRYPOINT ["./documents"]
