# movimientos_contables_mid

API MID para la gestión de movimientos contables dentro del sistema financiero de la universidad.

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones

* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)

### Variables de Entorno

Resumidas en [la plantilla](template.env)

```sh
# Copiar
cp template.env NOMBRE_COPIA.env

# Ajustar variables según sea necesario
# por ejemplo, con nano (puede usarse otro editor)
nano NOMBRE_COPIA.env

# Cargar variables
source NOMBRE_COPIA.env
# Hacer esto cada que cambie NOMBRE_COPIA.env
# o cada que se abra nuevamente el terminal
```

### Ejecución del Proyecto

```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/movimientos_contables_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/movimientos_contables_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
MOVIMIENTOS_CONTABLES_MID_HTTP_PORT=8080 MOVIMIENTOS_CONTABLES_MID_SOME_VARIABLE=some_value bee run
```

### Ejecución Dockerfile

```shell
# Implementado para despliegue del Sistema de integración continua CI.
```

### Ejecución docker-compose

```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/movimientos_contables_mid

#2. Moverse a la carpeta del repositorio
cd movimientos_contables_mid

#3. Crear un fichero con el nombre **custom.env**
touch custom.env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

### Ejecución Pruebas

Pruebas unitarias

```shell
# En Proceso
```

## Estado CI

| Develop | Release 0.1.0 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/movimientos_contables_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/movimientos_contables_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/movimientos_contables_mid/status.svg?ref=refs/heads/release/0.1.0)](https://hubci.portaloas.udistrital.edu.co/udistrital/movimientos_contables_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/movimientos_contables_mid/status.svg?ref=refs/heads/master)](https://hubci.portaloas.udistrital.edu.co/udistrital/movimientos_contables_mid) |

## Licencia

This file is part of movimientos_contables_mid.

movimientos_contables_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

movimientos_contables_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with movimientos_contables_mid. If not, see https://www.gnu.org/licenses/.
