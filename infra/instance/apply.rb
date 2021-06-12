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

instances = YAML.load_file('instances.yml')
node = YAML.load_file('node.yml')

github_keys = {}
github_users = ((node['admins'] || []) + (node['contestants'] || {}).values.flatten).uniq
github_users.each do |username|
  keys = Net::HTTP.get(URI.parse("https://github.com/#{username}.keys")).strip.split("\n").map(&:strip).sort.uniq
  github_keys[username] = keys.map { |k| "#{k} #{username}" }
end
File.write 'keys.yml', YAML.dump({ 'ssh_keys' => github_keys })

File.write 'apply.yml', YAML.dump(node.merge({ 'ssh_keys' => github_keys }))


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
