include_cookbook 'apt'

execute 'add-apt-repository -y ppa:redislabs/redis' do
  notifies :run, 'execute[apt update]', :immediately
  not_if 'test -f /etc/apt/sources.list.d/redislabs-ubuntu-redis-focal.list'
end

package 'redis'
