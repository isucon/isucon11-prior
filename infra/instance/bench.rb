require 'bundler/setup'
Bundler.require
require 'yaml'
require 'cli/ui'
require 'optparse'

num = 10

opt = OptionParser.new
opt.on('-n NUM', '--num=NUM') {|v| num = v.to_i }
opt.parse!(ARGV)

instances = YAML.load_file('instances.yml')[0...num]

def benchmark(host)
  Net::SSH.start(host, 'isuadmin', timeout: 5) do |ssh|
    ssh.exec!('sudo -u isucon bash -c "cd /home/isucon && /home/isucon/bin/benchmarker"').strip
  end
end

CLI::UI::StdoutRouter.enable
CLI::UI::SpinGroup.new(auto_debrief: true).tap do |group|
  instances.each do |ip|
    group.add(ip) do |spinner|
      begin
        stdout = benchmark(ip)
        score = stdout.split(/\n/).find { |l| l =~ /score: / } || ''
        spinner.update_title("#{ip}: #{score}")
      end
    end
    # group.wait
  end
  group.wait
end
