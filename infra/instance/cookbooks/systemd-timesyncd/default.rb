remote_file '/etc/systemd/timesyncd.conf' do
  owner 'root'
  group 'root'
  mode  '0644'
  notifies :restart, 'service[systemd-timesyncd.service]'
end

service 'systemd-timesyncd.service' do
  action [:start, :enable]
end
