# GORBAC - RABC con golang

**SPBAC** es un sistema de control de acceso basado en roles y paquetes de servicios, pensado para gestionar usuarios, aplicaciones, servicios y la asignación de permisos de manera flexible y escalable.

## Características principales

- **Gestión de Usuarios:** Alta, baja y modificación de usuarios.
- **Gestión de Roles y Permisos:** Asignación de roles a usuarios y permisos a roles.
- **Gestión de Productos:** Alta, baja y modificación de productos.
- **Autenticación tradicional y con Google:** Permite login local y mediante Google OAuth.
- **Control de acceso a endpoints:** Los permisos se gestionan a nivel de endpoint, según el rol del usuario.
- **API RESTful:** Backend desarrollado en Go, expone endpoints para todas las operaciones principales.

## Configuración de la base de datos

La aplicación utiliza un archivo `config.json` ubicado en la carpeta `config` para la configuración de la base de datos. Modifica estos parámetros según tu entorno:

```json
{
  "db_driver": "mysql",
  "db_host": "host",
  "db_port": "3306",
  "db_name": "gorbac",
  "db_user": "root",
  "db_password": "root"
}
```

En la carpeta `db` encontrarás el script SQL con la estructura de tablas y relaciones necesarias para el sistema.

## Instalación y ejecución

1. Clona el repositorio.
2. Configura la base de datos y el archivo `config.json`.
3. Ejecuta las migraciones o el script SQL de la carpeta `db`.
4. Compila y ejecuta el backend en Go.
5. (Opcional) Configura y ejecuta el frontend en React para la interfaz de usuario.

## Autenticación con Google

Para habilitar el login con Google:
- Configura el proveedor en Firebase y Google Cloud Console.
- Usa el Client ID correspondiente en el backend para validar los tokens.

## Licencia

MIT

---

**Desarrollado por Sebastián Bustelo**
