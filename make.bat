::service with consul
go build -o ./bin/web-backend-consul.exe ./cmd/consul/backend

go build -o ./bin/web-frontend-consul.exe ./cmd/consul/frontend

::service with nomad
go build -o ./bin/web-backend-nomad.exe ./cmd/nomad/backend

go build -o ./bin/web-frontend-nomad.exe ./cmd/nomad/frontend