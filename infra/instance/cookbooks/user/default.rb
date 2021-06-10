require 'net/http'
require 'uri'

node.reverse_merge!({
  contestants: {},
})

users = (node[:contestants][node[:hostname]] || [])
ssh_keys = users.map do |username|
  Net::HTTP.get(URI.parse("https://github.com/#{username}.keys")).strip
end

user 'isucon' do
  home '/home/isucon'
  shell '/bin/bash'
  create_home true
end

directory '/home/isucon/.ssh' do
  owner 'isucon'
  group 'isucon'
  mode '700'
end

remote_file '/home/isucon/.ssh/config' do
  owner 'isucon'
  group 'isucon'
  mode '600'
end

remote_file '/home/isucon/.ssh/id_ed25519' do
  owner 'isucon'
  group 'isucon'
  mode '600'
end

remote_file '/home/isucon/.ssh/id_ed25519.pub' do
  owner 'isucon'
  group 'isucon'
  mode '644'
end

file '/home/isucon/.ssh/authorized_keys' do
  owner 'isucon'
  group 'isucon'
  mode '0600'
  content ssh_keys.sort.uniq.join("\n").strip + "\n"
end

remote_file '/home/isucon/.bashrc' do
  owner 'isucon'
  group 'isucon'
  mode '644'
end

file '/etc/sudoers.d/isucon' do
  content "isucon ALL=(ALL) NOPASSWD:ALL\n"
  owner 'root'
  group 'root'
  mode '440'
end
