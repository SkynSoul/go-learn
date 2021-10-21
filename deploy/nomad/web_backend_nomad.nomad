job "web-backend-nomad" {
  region = "local-vm"
  datacenters = ["local-vm-1"]

  type = "service"

  # backend only
  constraint {
    attribute = "${node.unique.name}"
    operator  = "regexp"
    value     = "nomad-client-backend-[0-9]"
  }

  group "web-backend-nomad" {
    scaling {
      min = 1
      max = 10
    }

    network {
      port http {}
      port http_test {}
    }

    task "web-backend-nomad" {
      service {
        name = "web-backend-nomad"

        address_mode = "host"

        port = "http"

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
          "cd /root/services/web-backend-nomad && ./bin/web-backend-nomad --port ${NOMAD_PORT_http}",
        ]
      }

      resources {
        cpu    = 500
        memory = 128
      }
    }

    task "web-backend-nomad-test" {
      service {
        name = "web-backend-nomad"

        address_mode = "host"

        port = "http_test"

        check {
          type      = "http"
          port      = "http_test"
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
          "cd /root/services/web-backend-nomad && ./bin/web-backend-nomad --port ${NOMAD_PORT_http_test}",
        ]
      }

      resources {
        cpu    = 500
        memory = 128
      }
    }
  }
}
