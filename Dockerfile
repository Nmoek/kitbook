FROM ubuntu:20.04
COPY kitbook /app/kitbook
WORKDIR /app
CMD ["/app/kitbook"]
