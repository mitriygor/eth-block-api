FROM mongo:4.2.16-bionic
COPY ./ ./docker-entrypoint-initdb.d/
