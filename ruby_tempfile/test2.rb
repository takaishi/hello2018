require 'tempfile'
foo = Tempfile.open('nyah.yml') do |f|
  f.write 'hello world'
  ObjectSpace.each_object(Tempfile) {|f| p [f, f.path] }
  f.path
end
p foo
p File.read(foo)
ObjectSpace.each_object(Tempfile) {|f| p [f, f.path] }
GC.start
puts "Run GC!!!"
ObjectSpace.each_object(Tempfile) {|f| p [f, f.path] }
p foo
p File.read(foo)
