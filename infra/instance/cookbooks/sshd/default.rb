package 'openssh-server'

service 'ssh' do
  action [:enable, :start]
end

remote_file '/etc/ssh/sshd_config' do
  owner 'root'
  group 'root'
  mode '644'
  notifies :restart, 'service[ssh]'
end
