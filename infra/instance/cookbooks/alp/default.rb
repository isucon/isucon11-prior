package 'unzip'

execute 'install alp' do
  command <<-EOS
  wget -O /tmp/alp.zip https://github.com/tkuchiki/alp/releases/download/v1.0.3/alp_linux_amd64.zip
  cd /tmp && unzip alp.zip
  install -m 755 /tmp/alp /usr/local/bin/alp
  EOS
  not_if 'which alp'
end
