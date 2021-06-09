require 'bundler/setup'
Bundler.require

require 'simplecov'
SimpleCov.start
require_relative 'app'

run App.new
