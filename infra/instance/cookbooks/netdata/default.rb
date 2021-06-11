package 'netdata'
package 'python3'
package 'python3-mysqldb'

execute 'systemctl restart netdata' do
  action :nothing
end

remote_file '/etc/netdata/netdata.conf' do
  owner 'root'
  group 'root'
  mode '644'
  notifies :run, 'execute[systemctl restart netdata]'
end

remote_file '/etc/netdata/python.d/web_log.conf' do
  owner 'root'
  group 'root'
  mode '644'
  notifies :run, 'execute[systemctl restart netdata]'
end

execute %|mysql -uroot -e "create user netdata@localhost identified with mysql_native_password by 'netdata'; grant all privileges on *.* to netdata@localhost; flush privileges;"| do
  user 'root'
  not_if %(mysql -uroot -e "select User,Host from mysql.user"|grep -q netdata)
end

remote_file '/etc/netdata/python.d/mysql.conf' do
  owner 'root'
  group 'root'
  mode '644'
  notifies :run, 'execute[systemctl restart netdata]'
end
