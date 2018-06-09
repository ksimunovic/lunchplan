#docker-compose down -v && docker-compose up --build mongo-seed && docker-compose up
#docker-compose down -v && docker-compose up
FROM mongo

COPY . .
# Define mountable directories.
#VOLUME ["/data/db"]

# Define working directory.
#WORKDIR /data

ENV MONGO_INITDB_DATABASE=admin

# Define default command.
CMD ["mongod"]

# Expose ports.
#   - 27017: process
#   - 28017: http
EXPOSE 27017
EXPOSE 28017