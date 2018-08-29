resource "tls_private_key" "root" {
  algorithm   = "RSA"
}

resource "tls_self_signed_cert" "root" {
  key_algorithm   = "RSA"
  private_key_pem = "${tls_private_key.root.private_key_pem}"

  validity_period_hours = 26280
  early_renewal_hours   = 8760

  is_ca_certificate = true

  allowed_uses = ["cert_signing"]

  subject {
    common_name  = "hello2018.repl.info"
  }
}

resource "local_file" "root_ca_pem" {
  filename = "root_ca.pem"
  content  = "${tls_self_signed_cert.root.cert_pem}"
}

resource "local_file" "root_ca_key" {
  filename = "root_ca.key"
  content  = "${tls_private_key.root.private_key_pem}"
}
