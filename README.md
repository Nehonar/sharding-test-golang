# go-sharding-basic

Sistema de ejemplo en Go que implementa sharding básico usando `hash(username) % N` para distribuir usuarios entre múltiples shards de PostgreSQL.

Este proyecto demuestra:

- Separación clara en capas (API, Handler, Service, Router, Storage, Shards)
- Diseño limpio y mantenible
- Enrutamiento determinístico de escritura y lectura
- Sharding real en múltiples bases de datos PostgreSQL
- Arquitectura pensada para escalar horizontalmente
- Uso correcto de Go modules, interfaces y dependencias

---

## Arquitectura general

Flujo principal (CreateUser / GetUser):

Cliente → Handler → Service → Router → Shard → Base de datos

Cada capa tiene una responsabilidad única:

- **Handler**: HTTP (JSON ↔ dominio)
- **Service**: lógica de negocio
- **Router**: elige el shard correcto
- **Shard**: ejecuta SQL
- **Database**: almacena datos
- **Models**: estructuras de datos puras

---

## Decisión arquitectónica: Hash Mod N

El sharding se realiza aplicando SHA-256 al `username` y haciendo módulo del número de shards:

shardIndex = hash(username)[0] % numShards

Esto asegura que el mismo usuario siempre cae en el mismo shard.

---

## Limitaciones

Este método de sharding **NO es escalable en producción** porque:

- Cambiar el número de shards hace que todos los usuarios cambien de shard.
- No hay redistribución controlada.
- No existen réplicas ni tolerancia a fallos.
- No permite balanceo dinámico.

Este proyecto es un ejemplo educativo, no una implementación productiva.

---

## Ejecución

1. Levantar los shards:

```bash
docker compose up -d
```

2. Ejecutar el servidor:

```bash
go run cmd/server/main.go
```

3. Crear usuario:

```bash
curl -X POST http://localhost:8080/create-user

-H "Content-Type: application/json"
-d '{"user":"usuario","password":"1234"}'
```

4. Consultar usuario:

```bash
curl "http://localhost:8080/get-user?user=usuario"
```

---

## Objetivo educativo

Este proyecto forma parte de un portafolio orientado a demostrar:

- Habilidades de arquitectura de software
- Conocimiento de sistemas distribuidos
- Dominio de patrones avanzados
- Claridad en diseño y documentación