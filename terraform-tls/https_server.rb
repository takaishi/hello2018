require 'webrick'
require 'webrick/https'

cert = ARGV[0]
key = ARGV[1]
 
server = WEBrick::HTTPServer.new(
    :BindAddress => '127.0.0.1',
    :Port => 3000,
    :DocumentRoot => '.',
    :SSLEnable  => true,
    :SSLCertificate => OpenSSL::X509::Certificate.new(open(cert).read),
    :SSLPrivateKey => OpenSSL::PKey::RSA.new(open(key).read))
Signal.trap('INT') { server.shutdown }
server.start
