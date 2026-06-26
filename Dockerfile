FROM scratch

WORKDIR /app

COPY ./build/hou-yi ./

EXPOSE 8080

ENTRYPOINT ["/app/hou-yi"]