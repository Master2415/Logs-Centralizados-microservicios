
FROM postgres:latest

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=12345

COPY initUsers.sql /docker-entrypoint-initdb.d/
COPY initLogs.sql /docker-entrypoint-initdb.d/

# docker run --name server_database -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=12345 -d server_database