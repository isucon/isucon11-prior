package 'nginx'

execute 'systemctl restart nginx' do
  action :nothing
end

execute '/usr/sbin/logrotate -f /etc/logrotate.d/nginx' do
  action :nothing
end

remote_file '/etc/logrotate.d/nginx' do
  owner 'root'
  group 'root'
  mode '0644'
  notifies :run, 'execute[/usr/sbin/logrotate -f /etc/logrotate.d/nginx]', :delayed
end

remote_file '/etc/nginx/nginx.conf' do
  owner 'root'
  group 'root'
  mode '0644'
  notifies :run, 'execute[systemctl restart nginx]', :delayed
end

remote_file '/etc/nginx/sites-available/default' do
  owner 'root'
  group 'root'
  mode '0644'
  notifies :run, 'execute[systemctl restart nginx]', :delayed
end
