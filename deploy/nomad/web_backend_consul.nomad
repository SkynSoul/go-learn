job "web-backend-consul" {
  region = "local-vm"
  datacenters = ["local-vm-1"]

  type = "service"

  # backend only
  constraint {
    attribute = "${node.unique.name}"
    operator  = "regexp"
    value     = "nomad-client-backend-[0-9]"
  }

  group "web-backend-consul" {
    scaling {
      min = 1
      max = 10
    }

    network {
      port http {}
    }

    task "web-backend-consul" {
      service {
        name = "web-backend-consul"

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
          "cd /root/services/web-backend-consul && ./bin/web-backend-consul --port ${NOMAD_PORT_http}",
        ]
      }

      resources {
        cpu    = 500
        memory = 128
      }
    }
  }
}
