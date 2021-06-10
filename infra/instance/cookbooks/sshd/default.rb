execute 'systemctl restart sshd' do
  action :nothing
end

remote_file '/etc/ssh/sshd_config' do
  owner 'root'
  group 'root'
  mode '644'
  notifies :run, 'execute[systemctl restart sshd]', :delayed
end
