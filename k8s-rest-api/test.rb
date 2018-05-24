require 'kubeclient'
ssl_options = {
  client_cert: OpenSSL::X509::Certificate.new(File.read('/Users/r_takaishi/.minikube/client.crt')),
  client_key: OpenSSL::PKey::RSA.new(File.read('/Users/r_takaishi/.minikube/client.key')),
  ca_file: '/Users/r_takaishi/.minikube/ca.crt',
  verify_ssl:  OpenSSL::SSL::VERIFY_PEER
}
client = Kubeclient::Client.new('https://192.168.99.100:8443/api/', "v1", ssl_options: ssl_options)

actual_namespaces = client.get_namespaces


$namespaces = []
def labels(param)
 param
end

def namespace(name, &block)
 params = block.call
 $namespaces << Kubeclient::Resource.new({metadata: {name: name, labels: params}})
end

load './namespaces.rb'


$namespaces.each do |namespace|
  namespace.metadata.name
  if actual_namespaces.any?{|an| an.metadata.name == namespace.metadata.name}
    # 追加: 定義にあるけどクラスターにはない
    namespace.metadata.labels.to_h.keys.each do |key|
      an = actual_namespaces.find{|an| an.metadata.name == namespace.metadata.name}
      unless an.metadata.labels.to_h.has_key?(key)
        puts "#{namespace.metadata.name}にラベル #{key}: #{namespace.metadata.labels[key]}を追加"
      end
    end
    # 削除: 定義にないけどクラスターにはある
    an = actual_namespaces.find{|an| an.metadata.name == namespace.metadata.name}
    an.metadata.labels.to_h.keys.each do |key|
      unless namespace.metadata.labels.to_h.has_key?(key)
        puts "#{namespace.metadata.name}からラベル #{key}: #{an.metadata.labels[key]}を削除"
      end
    end
    # 更新定義にもクラスターにもあるけど値が違う
    namespace.metadata.labels.to_h.keys.each do |key|
      an = actual_namespaces.find{|an| an.metadata.name == namespace.metadata.name}
      if  an.metadata.labels.to_h.has_key?(key) && namespace.metadata.labels[key] != an.metadata.labels[key]
        puts "#{namespace.metadata.name}のラベルを更新 #{key}: #{namespace.metadata.labels[key]}"
      end
    end
    client.update_namespace(namespace)
  else
    puts "create: #{namespace.metadata.name}"
    client.create_namespace(namespace)
  end
end

actual_namespaces.each do |namespace|
  unless $namespaces.any?{|ns| ns.metadata.name == namespace.metadata.name}
    puts "delete #{namespace.metadata.name}"
    client.delete_namespace(namespace.metadata.name)
  end
end
