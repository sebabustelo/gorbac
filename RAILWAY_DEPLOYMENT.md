# Railway Deployment Guide

## Configuración para Railway

### 1. Variables de Entorno Requeridas

En Railway, necesitas configurar las siguientes variables de entorno:

#### Variables de Base de Datos:
```
DATABASE_HOST=mysql.railway.internal
DATABASE_PORT=3306
DATABASE_NAME=railway
DATABASE_USER=root
DATABASE_PASSWORD=tu_password_aqui
```

#### Variables de Entorno:
```
GO_ENV=railway
PORT=8229
```

### 2. Configuración en Railway Dashboard

1. Ve a tu proyecto en Railway
2. Selecciona tu servicio
3. Ve a la pestaña "Variables"
4. Agrega las variables de entorno mencionadas arriba

### 3. Configuración de Base de Datos

Si usas Railway MySQL:

1. Crea un servicio MySQL en Railway
2. Railway automáticamente proporcionará las variables de entorno:
   - `MYSQL_HOST`
   - `MYSQL_PORT` 
   - `MYSQL_DATABASE`
   - `MYSQL_USER`
   - `MYSQL_PASSWORD`

3. Conecta tu servicio de aplicación con el servicio MySQL

### 4. Deployment

#### Opción 1: GitHub Integration
1. Conecta tu repositorio de GitHub
2. Railway detectará automáticamente el Dockerfile
3. Configura las variables de entorno
4. Deploy automático

#### Opción 2: Manual Upload
1. Haz build local: `docker build -t gorbac .`
2. Sube la imagen a Railway
3. Configura las variables de entorno

### 5. Troubleshooting

#### Error: "The executable `docker` could not be found"

Este error en Railway generalmente indica:

1. **Problema de configuración de base de datos**
   - Verifica que las variables de entorno estén configuradas
   - Asegúrate de que la base de datos esté conectada

2. **Problema de puerto**
   - Railway asigna puertos dinámicamente
   - La aplicación debe usar `$PORT`

3. **Problema de permisos**
   - Los archivos deben tener permisos correctos

#### Verificar Logs

```bash
# En Railway Dashboard
# Ve a tu servicio > Logs
# Busca errores específicos
```

#### Health Check

La aplicación incluye un health check en `/roles`. Railway verificará:
- `GET /roles` debe responder con 200 OK

### 6. Configuración de Dominio

1. Ve a tu servicio en Railway
2. En la pestaña "Settings"
3. Configura tu dominio personalizado

### 7. Monitoreo

Railway proporciona:
- Logs en tiempo real
- Métricas de uso
- Health checks automáticos
- Alertas de fallos

### 8. Escalado

Para escalar tu aplicación:
1. Ve a tu servicio
2. Ajusta el número de réplicas
3. Railway manejará el balanceo de carga automáticamente

### 9. Backup y Recuperación

Railway maneja automáticamente:
- Backups de base de datos
- Rollbacks de deployments
- Recuperación de fallos

## Comandos Útiles

```bash
# Ver logs en tiempo real
railway logs

# Ver variables de entorno
railway variables

# Conectar a la base de datos
railway connect

# Deploy manual
railway up
```

## Soporte

Si tienes problemas:
1. Revisa los logs en Railway Dashboard
2. Verifica las variables de entorno
3. Asegúrate de que la base de datos esté conectada
4. Contacta al soporte de Railway si es necesario 