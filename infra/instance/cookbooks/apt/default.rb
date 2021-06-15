execute 'apt update' do
  action :nothing
end

execute 'apt upgrade -y' do
end
