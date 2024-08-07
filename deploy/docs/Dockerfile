# Use the official Swagger UI base image
FROM swaggerapi/swagger-ui

# Copy your Swagger spec file into the container
COPY ./docs/open-api.yaml /usr/share/nginx/html/swagger.yaml

ENV API_URL=swagger.yaml
ENV BASE_URL=/swagger

EXPOSE 8080