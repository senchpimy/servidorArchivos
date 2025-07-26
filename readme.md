
# 🗂️ Servidor de Archivos con Autenticación

![License](https://img.shields.io/badge/License-MIT-blue) ![Go](https://img.shields.io/badge/Built%20with-Go-brightgreen)

Este proyecto es un **servidor HTTP en Go** que permite compartir archivos de forma controlada. Ofrece autenticación básica (por nombre y contraseña) y la posibilidad de descargar archivos individuales o directorios comprimidos en `.zip`. También incluye una interfaz web simple.

---

## 🚀 Características

* 🔐 Acceso restringido a archivos mediante nombre y contraseña.
* 📁 Descarga de archivos y carpetas completas (como `.zip`).
* 🧩 Configuración dinámica mediante `config.json`.
* 🖼️ Endpoint especial para servir una imagen si se solicita `neco`.
* 🌐 Interfaz web generada con `index.html`.
* 📦 Manejo de archivos mediante JSON, GJSON y templates en Go.

---

## ⚙️ Configuración (`config.json`)

El archivo define los accesos permitidos y los archivos que puede ver cada usuario:

```json
{
  "usuario1": {
    "password": "none",
    "files": ["archivo1.txt", "carpeta/"]
  },
  "admin": {
    "password": "1234",
    "files": ["/ruta/segura/secreta.txt"]
  }
}
```

---

## 🧪 Uso

### Iniciar el servidor:

```bash
go run main.go
```

El servidor se inicia en [`http://localhost:8000`](http://localhost:8000)

### Realizar una solicitud de archivo:

Envía una solicitud POST con cuerpo JSON a `/`:

```json
{
  "name": "usuario1",
  "password": "none",
  "file": "archivo1.txt"
}
```

Si la contraseña es correcta y el archivo está autorizado en `config.json`, se devolverá como respuesta.

### Solicitud especial (Easter Egg)

```json
{ "name": "neco" }
```

Esto devuelve una imagen `GIF` especial.

---

## 🖥️ Interfaz Web

Al acceder a la raíz (`/`), se renderiza `index.html`, el cual puede mostrar un formulario o interfaz para el usuario final. Se inyectan automáticamente las variables `PORT` y `URL`.

---

## 🔐 Seguridad

* Solo los archivos listados explícitamente pueden ser servidos.
* Las carpetas solicitadas se comprimen al vuelo y se sirven como `.zip`.

---

## 📌 Notas Técnicas

* El servidor evita el acceso concurrente malicioso usando un contador de intentos (no implementado aún, `TRIES_TO_GET_A_FILE` es decorativo por ahora).
* Se ignoran subdirectorios dentro de carpetas compartidas.
* Todos los logs se imprimen en consola.

---

## 📄 Licencia

Este proyecto está licenciado bajo la [MIT License](LICENSE).
