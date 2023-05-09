# specify the base image to be user for the application, alpine or ubuntu
FROM golang:1.19-alpine

# create a working directory inside the image
WORKDIR /app

# copy Go modules and depedencies
COPY go.mod ./
COPY go.sum ./

# download Go modules and depedencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . ./

# compile application
RUN go build -o /genuine-id

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8000

# command to be used to execute when the image is used to start a container
CMD [ "/genuine-id" ]