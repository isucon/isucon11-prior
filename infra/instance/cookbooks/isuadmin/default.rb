node.reverse_merge!({
  admins: [],
  ssh_keys: {},
})

admins = (node[:admins] || [])

user 'isuadmin' do
  home '/home/isuadmin'
  shell '/bin/bash'
  create_home true
end

file '/etc/sudoers.d/isuadmin' do
  content "isuadmin ALL=(ALL) NOPASSWD:ALL\n"
  owner 'root'
  group 'root'
  mode '440'
end

directory '/home/isuadmin/.ssh' do
  owner 'isuadmin'
  group 'isuadmin'
  mode '700'
end

file '/home/isuadmin/.ssh/authorized_keys' do
  owner 'isuadmin'
  group 'isuadmin'
  mode '0600'
  content admins.map {|u| node[:ssh_keys][u] || [] }.flatten.sort.uniq.join("\n") + "\n"
end

remote_file '/home/isuadmin/.ssh/config' do
  owner 'isuadmin'
  group 'isuadmin'
  mode '600'
end
