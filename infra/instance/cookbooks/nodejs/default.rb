include_cookbook 'xbuild'

version = '16.3.0'

execute "rm -rf /home/isucon/local/nodejs; /opt/xbuild/node-install v#{version} /home/isucon/local/nodejs" do
  user 'isucon'
  not_if "/home/isucon/local/nodejs/bin/node --version | grep -q '^v#{version}$'"
end

execute 'curl -o- -L https://yarnpkg.com/install.sh | /home/isucon/.x bash' do
  user 'isucon'
  not_if "test -e /home/isucon/.yarn"
end
