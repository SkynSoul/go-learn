set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

::service with consul
go build -o ./bin/linux/web-backend-consul ./cmd/consul/backend

go build -o ./bin/linux/web-frontend-consul ./cmd/consul/frontend

::service with nomad
go build -o ./bin/linux/web-backend-nomad ./cmd/nomad/backend

go build -o ./bin/linux/web-frontend-nomad ./cmd/nomad/frontend