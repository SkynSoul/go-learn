job "web-frontend-consul" {
  region = "local-vm"
  datacenters = ["local-vm-1"]

  type = "system"

  # backend only
  constraint {
    attribute = "${node.unique.name}"
    operator  = "regexp"
    value     = "nomad-client-frontend-[0-9]"
  }

  group "web-frontend-consul" {
    network {
      port "http" {
        static = 80
      }
    }

    task "web-frontend-consul" {
      service {
        name = "web-frontend-consul"

        check {
          type      = "http"
          port      = "http"
          path      = "/health"
          interval  = "10s"
          timeout   = "2s"
        }
      }

      driver = "raw_exec"

      config {
        command = "/bin/bash"
        args = [
          "-c",
          "cd /root/services/web-frontend-consul && ./bin/web-frontend-consul --port ${NOMAD_PORT_http}",
        ]
      }

      resources {
        cpu    = 500
        memory = 128
      }
    }
  }
}
