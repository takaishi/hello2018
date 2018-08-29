require 'cool.io'

class FileWatcher < Cool.io::StatWatcher

  def initialize(path)
    super path, 0.1
  end

  def on_change(previous, current)
    p previous
    p current
    if current.nlink == 0
      puts 'deleted'
    end
  end
end

reactor = Cool.io::Loop.new
watcher = FileWatcher.new ARGV[0] 

reactor.attach watcher
reactor.run
