execute 'timedatectl set-timezone Asia/Tokyo' do
  not_if 'timedatectl | grep "Asia/Tokyo"'
end
