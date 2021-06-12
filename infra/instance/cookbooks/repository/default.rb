node[:isucon11_repository] = '/home/isuadmin/src/isucon11-prior'

execute 'rm -rf /home/isucon/src'

execute "git clone --depth 1 github.com:isucon/isucon11-prior.git #{node[:isucon11_repository]}" do
  user 'isuadmin'
  not_if "test -d #{node[:isucon11_repository]}"
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
  user 'isuadmin'
  cwd node[:isucon11_repository]
  not_if "cd #{node[:isucon11_repository]} && git fetch origin && test $(git rev-parse origin/main) = $(cat REVISION)"

  notifies :run, 'execute[build frontend]', :immediately
end

execute 'build frontend' do
  action :nothing
  command <<-EOS
  /home/isucon/.x yarn install --frozen-lockfile
  /home/isucon/.x yarn build
  rm -rf node_modules
  EOS
  user 'isuadmin'
  cwd "#{node[:isucon11_repository]}/webapp/frontend"
end
