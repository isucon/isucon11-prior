include_cookbook 'curl'
include_cookbook 'git'

execute 'rm -rf /home/isucon/.asdf' do
  only_if 'test -d /home/isucon/.asdf'
end

execute 'rm -rf /home/isucon/.asdfrc' do
  only_if 'test -f /home/isucon/.asdfrc'
end

execute 'rm -rf /home/isucon/.tool-versions' do
  only_if 'test -f /home/isucon/.tool-versions'
end
