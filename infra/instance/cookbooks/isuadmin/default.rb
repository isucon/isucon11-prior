node.reverse_merge!({
  admins: [],
  ssh_keys: {},
})

admins = (node[:admins] || [])

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

remote_file '/home/isuadmin/.ssh/deploy' do
  owner 'isuadmin'
  group 'isuadmin'
  mode '600'
end

remote_file '/home/isuadmin/.ssh/deploy.pub' do
  owner 'isuadmin'
  group 'isuadmin'
  mode '644'
end
