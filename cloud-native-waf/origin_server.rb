require 'webrick'

srv = WEBrick::HTTPServer.new(:Port => 8000, :DocumentRoot => Dir.pwd)

srv.mount_proc '/' do |req, res|
  res.body = "Hello, World!"
end


srv.start
