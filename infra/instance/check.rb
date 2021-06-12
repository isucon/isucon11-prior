require 'bundler/setup'
Bundler.require
require 'yaml'
require 'net/http'
require 'net/ping'
require 'cli/ui'
require 'optparse'

timeout = 10

opt = OptionParser.new
opt.on('-t NUM', '--timeout=NUM') {|v| timeout = v.to_i }
opt.parse!(ARGV)

instances = YAML.load_file('instances.yml')

def ping?(host)
  $stderr.puts host
  check = Net::Ping::External.new(host)
  check.ping?
rescue Exception
  false
end

def ssh?(host)
  Net::SSH.start(host, 'isuadmin', timeout: 1) do |ssh|
    ssh.exec!('echo ok').strip == 'ok'
  end
rescue Exception
  false
end

def app?(host)
  nginx = Net::HTTP.start(host, 80) { |http| http.get('/') }.code == '200'
  api = Net::HTTP.start(host, 80) { |http| http.get('/api/schedules') }.code == '200'
  nginx && api
rescue Exception
  false
end

def netdata?(host)
  Net::HTTP.start(host, 19999) { |http| http.get('/') }.code == '200'
rescue Exception
  false
end

CLI::UI::StdoutRouter.enable
CLI::UI::SpinGroup.new(auto_debrief: false).tap do |group|
  instances.each do |ip|
    group.add(ip) do |spinner|
      begin
        block = proc do
          spinner.update_title("#{ip}: ping")
          until ping?(ip)
            sleep 1
          end

          spinner.update_title("#{ip}: ssh")
          until ssh?(ip)
            sleep 1
          end

          spinner.update_title("#{ip}: app")
          until app?(ip)
            sleep 1
          end

          spinner.update_title("#{ip}: netdata")
          until netdata?(ip)
            sleep 1
          end

          spinner.update_title("#{ip}: OK")
        end

        timeout > 0 ? Timeout.timeout(timeout, &block) : block.call
      rescue
        CLI::UI::Spinner::TASK_FAILED
      end
    end
    group.wait
  end
  group.wait
end
