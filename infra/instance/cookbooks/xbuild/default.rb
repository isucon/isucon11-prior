include_cookbook 'user'

directory "/opt/xbuild" do
  owner 'isucon'
  group 'isucon'
  mode  '0755'
end

execute "git clone --depth 1 https://github.com/tagomoris/xbuild /opt/xbuild" do
  user 'isucon'
  not_if "test -e /opt/xbuild/.git"
end

directory '/home/isucon/bin' do
  owner 'isucon'
  group 'isucon'
  mode  '0755'
end

directory '/home/isucon/local' do
  owner 'isucon'
  group 'isucon'
  mode  '0755'
end

remote_file '/home/isucon/.local.env' do
  owner 'isucon'
  group 'isucon'
  mode  '0644'
end

remote_file '/home/isucon/.x' do
  owner 'isucon'
  group 'isucon'
  mode  '0755'
end
