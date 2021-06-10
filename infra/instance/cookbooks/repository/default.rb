execute 'git clone --depth 1 git@github.com:isucon/isucon11-sc.git /home/isucon/src/github.com/isucon/isucon11-sc' do
  user 'isucon'
  not_if 'test -d /home/isucon/src/github.com/isucon/isucon11-sc'
end

execute 'update repository' do
  command <<-EOC
  git reset --
  git checkout -- .
  git clean -f
  git checkout main
  git pull origin
  git rev-parse HEAD > REVISION
  EOC
  user 'isucon'
  cwd '/home/isucon/src/github.com/isucon/isucon11-sc'
  not_if 'cd /home/isucon/src/github.com/isucon/isucon11-sc && git fetch origin && test $(git rev-parse origin/main) = $(cat REVISION)'

  notifies :run, 'execute[build frontend]', :immediately
end

execute 'build frontend' do
  action :nothing
  command <<-EOS
  /home/isucon/.x yarn install --frozen-lockfile
  /home/isucon/.x yarn build
  rm -rf node_modules
  EOS
  user 'isucon'
  cwd '/home/isucon/src/github.com/isucon/isucon11-sc/webapp/frontend'
end
