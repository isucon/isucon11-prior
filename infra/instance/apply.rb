require 'bundler/setup'
Bundler.require
require "yaml"
require "open3"
require "logger"
require 'net/http'
require 'uri'
require 'optparse'

parallelism = 20

opt = OptionParser.new
opt.on('-p NUM', '--parallelism=NUM') {|v| parallelism = v.to_i }
opt.parse!(ARGV)

instances = YAML.load_file('instances.yml') rescue []
contestants = YAML.load_file('contestants.yml') rescue []
node = YAML.load_file('node.yml') rescue {}
keys_cache = (YAML.load_file('keys.yml') rescue {})

client = Net::HTTP.new('github.com', '443')
client.use_ssl = true

github_keys = {}
github_users = ((node['admins'] || []) + (node['contestants'] || {}).values.flatten + contestants).uniq
github_users.each do |username|
  username = username.strip
  keys = (keys_cache['ssh_keys'] || {})[username] || []
  if keys.empty?
    puts "fetch #{username}'s keys..."
    client.start do |http|
      res = http.request_get("/#{username}.keys")
      if res.code == '200'
        keys = res.read_body.strip.split("\n").map(&:strip).sort.uniq.reject(&:empty?).map { |k| "#{k} #{username}" }
        puts "WARNING: #{username}'s keys: empty" if keys.empty?
      else
        puts "WARNING: #{username}'s keys: #{res.code}"
      end
    end
  end
  github_keys[username] = keys
end
File.write 'keys.yml', YAML.dump({ 'ssh_keys' => github_keys })

servers = node['contestants']
server_names = servers.keys.sort
contestants.each_with_index do |user, idx|
  server = servers[server_names[idx]]
  raise 'Server not found' if server.nil?
  server << user
end

File.write 'apply.yml', YAML.dump(node.merge({ 'ssh_keys' => github_keys }))

exit 0

Parallel.each(instances, in_processes: parallelism) do |ip|
  name = '%03d' % ip.split('.').last.to_i
  logger = Logger.new(STDOUT, progname: name, datetime_format: "%H:%M:%S")
  logger.info("start")

  Open3.popen3("bundle exec itamae ssh --node-yaml apply.yml --host #{ip} recipe.rb") do |i, o, e, w|
    i.close
    o.each { |line| logger.info(line.rstrip) }
    e.each { |line| logger.error(line.rstrip) }
    w.value
  end
end
