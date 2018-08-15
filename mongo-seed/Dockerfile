#docker-compose down -v && docker-compose up --build mongo-seed && docker-compose up
#docker-compose down -v && docker-compose up
FROM mongo
COPY . .
ENV MONGO_INITDB_DATABASE=admin
CMD ["mongod"]
EXPOSE 27017
EXPOSE 28017