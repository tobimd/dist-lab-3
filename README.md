# Laboratorio 03 - Sistemas Distribuidos 2021-2

## Integrantes

|    Integrante    |              Correo             |     Rol     |
|----------------|-------------------------------|-----------|
| Pablo Aravena    | pablo.aravenan@sansano.usm.cl   | 201773044-7   |
| Esteban Carrillo | esteban.carrillo@sansano.usm.cl | 201773032-3 |
| César Paulangelo | cesar.paulangelo@sansano.usm.cl | 201773088-9   |

---

## Como correr el programa

**tl;dr**:

- Cambiar directorio a `~/lab03_grupo46`

- Usar `make dist-<#vm>` (e.g. para dist-181: `make dist-181`)

Para que funcione el programa, se dispuso de un Makefile en cada máquina virtual (dentro del directorio del proyecto `~/lab03_grupo46`), donde sólo basta utilizar el comando `make dist-<#vm>` usando `<#vm>` como el número asignado en la máquina (por ejemplo, si la máquina es dist-181, entonces este se usaría como `make dist-181`)

Por detrás, en cada máquina se dispuso un `target` por defecto en el Makefile que corresponde a la máquina actual (e.g. el target `dist-181` estará presente en todas las máquinas). Este tipo de targets (`dist-<#vm>` son para correr automáticamente todos las entidades relacionadas a alguna máquina. Por otro lado, si se desea correr una entidad en específico, se puede usar `make <nombre_entidad> [id_entidad=0]` como por ejemplo: `make fulcrum` (entidad fulcrum con id 0) o `make leia` que independiente de que se incluya o no un id, se ignorará ya que la entidad solo es una.

Para terminar todos los procesos creados con el Make, se utiliza `make stop` y para borrar los archivos generados (los logs y registros) se usa `make clean`. Para esto también se provee `make reset` que realiza estos dos al mismo tiempo.

---

## Inputs del usuario

Las entidades que pueden recibir inputs son sólo Leia y los dos Informantes, por lo que estos están en máquinas distintas. Ambas reciben los comandos dispuestos en el enunciado de la misma forma (e.g. "AddCity" y no de otra forma), el cual debe ser preciso ya que si hay errores en la estructura, esto podría hacer que se rompa el programa.

La forma de los comandos es:

`<comando> <nombre_planeta> <nombre_ciudad> [argumento]`

## Estrategia de conflictos

Para solucionar posibles conflictos al momento practicar la eventual consistencia, Se comparara cantidad de cambios y se replicara en los servidores fulcrum los cambios realizados en el servidor que tenga una mayor de comandos realizados. Descartando los cambios de los otros servidores
