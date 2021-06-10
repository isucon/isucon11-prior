execute 'netplan apply' do
  action :nothing
end

remote_file '/etc/netplan/90-dns.yaml' do
  notifies :run, 'execute[netplan apply]', :immediately
end

execute 'systemctl restart systemd-resolved' do
  action :nothing
end

remote_file '/etc/systemd/resolved.conf' do
  notifies :run, 'execute[systemctl restart systemd-resolved]', :immediately
end
