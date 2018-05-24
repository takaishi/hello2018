
namespace "default" do
end

namespace "kube-system" do
end

namespace "kube-public" do
end

namespace "foo001" do
  labels name: 'staging', test: '123', test2: '456'
end

namespace "foo002" do
  labels name: 'production', service: 'aaaaaaaaaaa'
end

namespace "foo003" do
  labels name: 'development', service: 'fuga'
end
