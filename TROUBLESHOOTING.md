# Troubleshooting Guide

## Container Startup Issues

### Error: "The executable `docker` could not be found"

This error message is misleading and usually indicates one of the following issues:

1. **Database Connection Problem**
2. **Missing Dependencies**
3. **Permission Issues**
4. **Configuration Problems**

## Debugging Steps

### 1. Test Database Connection

Run the database connection test:

```bash
# Build and run the test
docker build -t gorbac-test .
docker run --rm gorbac-test go run test_db.go
```

### 2. Check Container Logs

```bash
# View container logs
docker logs <container_name>

# Or if using docker-compose
docker-compose logs api-rbac
```

### 3. Test with Docker Compose

```bash
# Start the full stack
docker-compose up -d

# Check if database is ready
docker-compose logs database-mariadb

# Test the API
curl http://localhost:8229/roles
```

### 4. Manual Container Testing

```bash
# Run container interactively
docker run -it --rm gorbac-test /bin/bash

# Inside container, test:
./test_db.go
./main
```

## Common Solutions

### 1. Database Not Available

If the database is not running or accessible:

```bash
# Start database first
docker-compose up -d database-mariadb

# Wait for database to be ready
docker-compose logs database-mariadb

# Then start the API
docker-compose up -d api-rbac
```

### 2. Environment Configuration

Ensure the correct environment is set:

```bash
# For local development
export GO_ENV=local

# For production
export GO_ENV=production
```

### 3. Network Issues

If containers can't communicate:

```bash
# Check network
docker network ls
docker network inspect gorbac_desarrollo

# Ensure containers are on the same network
docker-compose up -d
```

### 4. Permission Issues

If there are permission problems:

```bash
# Rebuild with proper permissions
docker build --no-cache -t gorbac-test .
```

## Health Checks

The container includes health checks. Monitor them:

```bash
# Check container health
docker ps
docker inspect <container_name> | grep Health -A 10
```

## Logs Directory

If you need persistent logs, the override file mounts a logs directory:

```bash
# Create logs directory
mkdir -p logs

# View logs
tail -f logs/app.log
```

## Railway Deployment

For Railway deployment, ensure:

1. Environment variables are set correctly
2. Database connection string is valid
3. Port 8229 is exposed
4. Health check endpoint is accessible

## Contact

If issues persist, check:
- Container logs
- Database connectivity
- Environment configuration
- Network connectivity between services 