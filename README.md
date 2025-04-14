# Data Provider Service

## Docker

### Build
```bash
docker build -t data-provider-service -f ./.docker/Dockerfile .
```

### Run
```bash
docker run -p 8080:8080 data-provider-service
```

## Docker Compose

### Build
```bash
docker-compose -f ./.docker/docker-compose.yaml build
```

### Run
```bash
docker-compose -f ./.docker/docker-compose.yaml up -d
```