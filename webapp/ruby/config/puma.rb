root = File.expand_path('..', __dir__)

workers 1

directory root
rackup File.join(root, 'config.ru')
bind 'tcp://0.0.0.0:9292'
environment ENV.fetch('RACK_ENV') { 'development' }
pidfile File.join(root, 'tmp', 'puma.pid')
