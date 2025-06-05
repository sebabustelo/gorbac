# SPBAC  [Services Packets Base Access Control]

**Sistema de control de acceso basado en paquetes de servicios**

## Caracteristicas

 * Gestionar Usuarios (ABM)
 * Gestionar Aplicaciones (ABM)
 * Gestionar Servicios (ABM)
 * Gestionar Paquetes (son un conjunto de servicios que pueden ser asignado a uno o muchos usuarios) (ABM)

## Configuración de la base de datos

La aplicación lee la configuración de la base de datos del archivo `config.json` ubicado en la carpeta config, modificar estos
parametros segun corresponda.
En la carpeta db se encuentra el script sql con las tablas y las relaciones de la base de datos.

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
