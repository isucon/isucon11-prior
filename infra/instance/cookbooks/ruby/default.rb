include_cookbook 'xbuild'

version = '3.0.1'

execute "rm -rf /home/isucon/local/ruby; /opt/xbuild/ruby-install #{version} /home/isucon/local/ruby" do
  user 'isucon'
  not_if "/home/isucon/local/ruby/bin/ruby -v | grep '^ruby #{version}p'"
end

execute "/home/isucon/.x gem install bundler -v '~> 2.2.3' --no-doc" do
  user 'isucon'
  not_if "/home/isucon/local/ruby/bin/bundle version | grep -q 'version 2.2'"
end
