docker build -t a2go-client . \
  && docker run --rm -it -p 8080:8080 a2go-client
