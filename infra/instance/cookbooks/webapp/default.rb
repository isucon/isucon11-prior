include_cookbook 'systemd'
include_cookbook 'repository'

remote_file '/home/isucon/env' do
  owner 'isucon'
  group 'isucon'
  mode  '0644'
end

execute 'install webapp' do
  command <<-EOS
  rm -rf /home/isucon/webapp
  cp -a #{node[:isucon11_repository]}/webapp /home/isucon/webapp
  cp #{node[:isucon11_repository]}/REVISION /home/isucon/webapp/REVISION
  chown -R isucon:isucon /home/isucon/webapp
  EOS
  cwd '/home/isucon'
  not_if "test -d /home/isucon/webapp && test -f /home/isucon/webapp/REVISION && test $(cat /home/isucon/webapp/REVISION) = $(cat #{node[:isucon11_repository]}/REVISION)"

  notifies :run, 'execute[/home/isucon/webapp/tools/initdb]', :immediately
  notifies :run, 'execute[bundle install]', :immediately
  notifies :restart, 'service[web-ruby]'
end

execute '/home/isucon/webapp/tools/initdb' do
  action :nothing
  user 'isucon'
  cwd '/home/isucon/webapp'
end

# systemctl

remote_file '/etc/systemd/system/web-ruby.service' do
  owner 'root'
  group 'root'
  mode  '0644'
  notifies :run, 'execute[systemctl daemon-reload]', :immediately
  notifies :restart, 'service[web-ruby]'
end

remote_file '/etc/systemd/system/web-golang.service' do
  owner 'root'
  group 'root'
  mode  '0644'
  notifies :run, 'execute[systemctl daemon-reload]', :immediately
end

service 'web-ruby' do
  action [:enable, :start]
end

service 'web-golang' do
  action [:disable, :stop]
end

# ruby

execute 'bundle install' do
  action :nothing
  command <<-EOS
  /home/isucon/.x bundle config set deployment true
  /home/isucon/.x bundle config set path vendor/bundle
  /home/isucon/.x bundle install -j8
  /home/isucon/.x bundle config set deployment false
  EOS
  user 'isucon'
  cwd '/home/isucon/webapp/ruby'
  not_if 'cd /home/isucon/webapp/ruby && test -e .bundle && /home/isucon/.x bundle check'
  notifies :restart, 'service[web-ruby]'
end

execute '/home/isucon/.x bundle config set deployment false' do
  user 'isucon'
  cwd '/home/isucon/webapp/ruby'
  only_if 'cd /home/isucon/webapp/ruby && /home/isucon/.x bundle config get --parseable deployment | grep "deployment=true"'
end

# golang

execute '/home/isucon/.x make build' do
  user 'isucon'
  cwd '/home/isucon/webapp/golang'
  only_if 'test -x /home/isucon/webapp/golang/bin/webapp'
end
