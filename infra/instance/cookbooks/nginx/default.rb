package 'nginx'

service 'nginx' do
  action [:enable, :start]
end

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
  notifies :restart, 'service[nginx]'
end

remote_file '/etc/nginx/nginx.conf' do
  owner 'root'
  group 'root'
  mode '0644'
  notifies :restart, 'service[nginx]'
end

directory '/usr/local/ssl' do
  owner 'root'
  group 'root'
  mode '0755'
end

# remote_file '/usr/local/ssl/fullchain.pem' do
#   owner 'root'
#   group 'root'
#   mode '0644'
#   notifies :restart, 'service[nginx]'
# end
#
# remote_file '/usr/local/ssl/privkey.pem' do
#   owner 'root'
#   group 'root'
#   mode '0644'
#   notifies :restart, 'service[nginx]'
# end

remote_file '/etc/nginx/sites-available/default' do
  owner 'root'
  group 'root'
  mode '0644'
  notifies :restart, 'service[nginx]'
end
