FROM debian:latest
COPY discordhobbyist* /discordhobbyist
ENTRYPOINT ["/discordhobbyist"]
