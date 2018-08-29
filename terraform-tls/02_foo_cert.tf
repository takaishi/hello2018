resource "tls_private_key" "foo" {
  algorithm   = "RSA"
}

resource "tls_cert_request" "foo" {
  key_algorithm   = "${tls_private_key.foo.algorithm}"
  private_key_pem = "${tls_private_key.foo.private_key_pem}"

  subject {
    common_name  = "foo.hello2018.repl.info"
  }

  dns_names = [
    "localhost",
  ]
}

resource "tls_locally_signed_cert" "foo" {
  ca_key_algorithm   = "${tls_private_key.root.algorithm}"

  cert_request_pem   = "${tls_cert_request.foo.cert_request_pem}"

  ca_private_key_pem = "${tls_private_key.root.private_key_pem}"
  ca_cert_pem        = "${tls_self_signed_cert.root.cert_pem}"

  validity_period_hours = 12

  allowed_uses = [
    "server_auth",
  ]
}

resource "local_file" "foo_cert_pem" {
  filename = "foo.pem"
  content  = "${tls_locally_signed_cert.foo.cert_pem}"
}

resource "local_file" "foo_cert_key" {
  filename = "foo.key"
  content  = "${tls_private_key.foo.private_key_pem}"
}
