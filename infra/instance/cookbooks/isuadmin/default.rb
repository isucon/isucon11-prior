require 'net/http'
require 'uri'

node.reverse_merge!({
  admins: []
})

ssh_keys = node[:admins].map do |username|
  Net::HTTP.get(URI.parse("https://github.com/#{username}.keys")).strip
end

file '/home/isuadmin/.ssh/authorized_keys' do
  owner 'isuadmin'
  group 'isuadmin'
  mode '0600'
  content ssh_keys.sort.uniq.join("\n").strip + "\n"
end
