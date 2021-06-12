include_cookbook 'repository'

execute 'build benchmarker' do
  command <<-EOS
  /home/isucon/.x make build -B
  sudo install -m 755 ./bin/benchmarker /home/isucon/bin/benchmarker
  EOS
  user 'isuadmin'
  cwd "#{node[:isucon11_repository]}/benchmarker"
  not_if "test -x /home/isucon/bin/benchmarker && test $(/home/isucon/bin/benchmarker -version) = $(cat #{node[:isucon11_repository]}/REVISION)"
end
