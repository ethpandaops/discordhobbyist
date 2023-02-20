FROM gcr.io/distroless/static-debian11:latest
COPY discordhobbyist* /discordhobbyist
ENTRYPOINT ["/discordhobbyist"]
