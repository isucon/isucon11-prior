%w(
  mysql-server-8.0
  mysql-server-core-8.0
  mysql-client-8.0
  mysql-client-core-8.0

  libmysqlclient-dev

  mysql-client
  mysql-common
  mysql-server
).each do |_|
  package _
end

remote_file "/etc/mysql/mysql.conf.d/mysqld.cnf" do
  owner 'root'
  group 'root'
  mode  '0644'
  notifies :restart, 'service[mysql.service]'
end

service "mysql.service" do
  action [:enable, :start]
end

execute %|mysql -uroot -e "create user isucon@localhost identified with mysql_native_password by 'isucon'; grant all privileges on *.* to isucon@localhost; flush privileges;"| do
  user 'root'
  not_if %(mysql -uroot -e "select User,Host from mysql.user"|grep -q isucon)
end
