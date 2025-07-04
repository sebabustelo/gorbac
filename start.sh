#!/bin/bash

# Startup script for the RBAC API

echo "Starting RBAC API..."

# Check if the main executable exists
if [ ! -f "./main" ]; then
    echo "ERROR: main executable not found!"
    exit 1
fi

# Check if config directory exists
if [ ! -d "./config" ]; then
    echo "ERROR: config directory not found!"
    exit 1
fi

# Check if RSA keys exist
if [ ! -f "./private.rsa" ] || [ ! -f "./public.rsa.pub" ]; then
    echo "WARNING: RSA keys not found!"
fi

# Detect Railway environment
if [ ! -z "$RAILWAY_ENVIRONMENT" ]; then
    export GO_ENV="railway"
    echo "Railway environment detected, setting GO_ENV to 'railway'"
elif [ ! -z "$PORT" ]; then
    export GO_ENV="railway"
    echo "Railway PORT detected, setting GO_ENV to 'railway'"
elif [ -z "$GO_ENV" ]; then
    export GO_ENV="local"
    echo "Setting GO_ENV to 'local'"
fi

echo "Environment: $GO_ENV"
echo "Railway Environment: $RAILWAY_ENVIRONMENT"
echo "Port: $PORT"

# Check database environment variables for Railway
if [ "$GO_ENV" = "railway" ]; then
    echo "Checking Railway database configuration..."
    
    # Check if database variables are set
    if [ -z "$DATABASE_HOST" ] && [ -z "$MYSQL_HOST" ]; then
        echo "WARNING: DATABASE_HOST not set, using default config"
    fi
    
    if [ -z "$DATABASE_PORT" ] && [ -z "$MYSQL_PORT" ]; then
        echo "WARNING: DATABASE_PORT not set, using default config"
    fi
    
    if [ -z "$DATABASE_NAME" ] && [ -z "$MYSQL_DATABASE" ]; then
        echo "WARNING: DATABASE_NAME not set, using default config"
    fi
    
    if [ -z "$DATABASE_USER" ] && [ -z "$MYSQL_USER" ]; then
        echo "WARNING: DATABASE_USER not set, using default config"
    fi
    
    if [ -z "$DATABASE_PASSWORD" ] && [ -z "$MYSQL_PASSWORD" ]; then
        echo "WARNING: DATABASE_PASSWORD not set, using default config"
    fi
fi

echo "Starting application..."

# Use Railway's PORT if available, otherwise use default 8229
if [ ! -z "$PORT" ]; then
    echo "Using Railway PORT: $PORT"
    # Modify the application to use the PORT environment variable
    export APP_PORT="$PORT"
else
    echo "Using default port: 8229"
    export APP_PORT="8229"
fi

# Run the application
exec ./main 