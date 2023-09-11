# Usar una imagen de Go como base
FROM golang:latest

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el contenido de la carpeta actual (donde se encuentra tu código) al directorio de trabajo en el contenedor
COPY . /app

# Compilar tu aplicación
RUN go build -o main

# Exponer el puerto en el que tu aplicación escucha
EXPOSE 9000

# Comando para ejecutar tu aplicación cuando el contenedor se inicie
CMD ["./main"]