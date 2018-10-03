require 'tempfile'

foo = Tempfile.open('nyah.yml') do |f|
  f.write 'hello world'
  f
end

puts "GCされていないオブジェクト: #{ObjectSpace.each_object(Tempfile).map {|f| [f, f.path].inspect }}"
GC.start
puts "GCしたよ"
puts "GCされていないオブジェクト: #{ObjectSpace.each_object(Tempfile).map {|f| [f, f.path].inspect }}"
p File.read(foo.path)
