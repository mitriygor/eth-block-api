FROM mongo:4.2.16-bionic
COPY ./mongo-init/ ./docker-entrypoint-initdb.d/
