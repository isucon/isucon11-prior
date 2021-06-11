include_cookbook 'repository'

execute 'build benchmarker' do
  command <<-EOS
  /home/isucon/.x make build -B
  install -m 755 ./bin/benchmarker /home/isucon/bin/benchmarker
  EOS
  user 'isucon'
  cwd '/home/isucon/src/github.com/isucon/isucon11-sc/benchmarker'
  not_if 'test -x /home/isucon/bin/benchmarker && test $(/home/isucon/bin/benchmarker -version) = $(cat /home/isucon/src/github.com/isucon/isucon11-sc/REVISION)'
end
