# Usa una imagen base de Golang
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /go/src/app

# Copia los archivos del proyecto al directorio de trabajo del contenedor
COPY ./api .
COPY ./.env .

# Descarga las dependencias del proyecto
RUN go get -d -v ./...

# Compila la aplicación
RUN go install -v ./...

# Expone el puerto en el que la aplicación se ejecutará
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["app"]