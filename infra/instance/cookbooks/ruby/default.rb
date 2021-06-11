include_cookbook 'xbuild'

%w[
  autoconf bison build-essential libssl-dev libyaml-dev libreadline-dev zlib1g-dev libncurses5-dev libffi-dev libgdbm6 libgdbm-dev libdb-dev
].each do |_|
  package _
end

version = '3.0.1'

execute "rm -rf /home/isucon/local/ruby; /opt/xbuild/ruby-install #{version} /home/isucon/local/ruby" do
  user 'isucon'
  not_if "/home/isucon/local/ruby/bin/ruby -v | grep '^ruby #{version}p'"
end

execute "/home/isucon/.x gem install bundler -v '~> 2.2.3' --no-doc" do
  user 'isucon'
  not_if "/home/isucon/local/ruby/bin/bundle version | grep -q 'version 2.2'"
end
