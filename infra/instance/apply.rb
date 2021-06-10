require 'bundler/setup'
Bundler.require
require "yaml"
require "open3"
require "logger"

MITAMAE_URL = "https://github.com/itamae-kitchen/mitamae/releases/download/v1.12.6/mitamae-x86_64-linux"

instances = YAML.load_file('instances.yml')

Parallel.each(instances) do |ip|
  name = '%03d' % ip.split('.').last.to_i
  logger = Logger.new(STDOUT, progname: name, datetime_format: "%H:%M:%S")
  logger.info("start")

  Open3.popen3("bundle exec itamae ssh --node-yaml node.yml --host #{ip} recipe.rb") do |i, o, e, w|
    i.close
    o.each { |line| logger.info(line.rstrip) }
    e.each { |line| logger.error(line.rstrip) }
    w.value
  end
end
