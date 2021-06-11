execute 'netplan apply' do
  action :nothing
end

remote_file '/etc/netplan/90-dns.yaml' do
  notifies :run, 'execute[netplan apply]', :immediately
end

remote_file '/etc/systemd/resolved.conf' do
  notifies :restart, 'service[systemd-resolved]'
end

service 'systemd-resolved' do
  action [:enable, :start]
end
