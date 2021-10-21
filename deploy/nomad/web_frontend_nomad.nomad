job "web-frontend-nomad" {
  region = "local-vm"
  datacenters = ["local-vm-1"]

  type = "system"

  # backend only
  constraint {
    attribute = "${node.unique.name}"
    operator  = "regexp"
    value     = "nomad-client-frontend-[0-9]"
  }

  group "web-frontend-nomad" {
    network {
      port "http" {
        static = 81
      }
    }

    task "web-frontend-nomad" {
      service {
        name = "web-frontend-nomad"

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
          "cd /root/services/web-frontend-nomad && ./bin/web-frontend-nomad --port ${NOMAD_PORT_http}",
        ]
      }

      resources {
        cpu    = 500
        memory = 128
      }
    }
  }
}
