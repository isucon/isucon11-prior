execute 'curl -s https://install.speedtest.net/app/cli/install.deb.sh | bash' do
  not_if 'test -f /etc/apt/sources.list.d/ookla_speedtest-cli.list'
end

package 'speedtest'
