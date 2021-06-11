node.reverse_merge!({
  contestants: {},
  ssh_keys: {},
})

users = (node[:contestants][node[:hostname]] || [])

user 'isucon' do
  home '/home/isucon'
  shell '/bin/bash'
  create_home true
end

remote_file '/home/isucon/.gitconfig' do
  owner 'isucon'
  group 'isucon'
  mode '644'
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

remote_file '/home/isucon/.ssh/isucon' do
  owner 'isucon'
  group 'isucon'
  mode '600'
end

remote_file '/home/isucon/.ssh/isucon.pub' do
  owner 'isucon'
  group 'isucon'
  mode '644'
end

file '/home/isucon/.ssh/authorized_keys' do
  owner 'isucon'
  group 'isucon'
  mode '0600'
  content users.map {|u| node[:ssh_keys][u] || [] }.flatten.sort.uniq.join("\n") + "\n"
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

execute 'ssh-keygen -b 2048 -t rsa -f /home/isucon/.ssh/id_rsa -q -N ""' do
  user 'isucon'
  not_if 'test -f /home/isucon/.ssh/id_rsa'
end
