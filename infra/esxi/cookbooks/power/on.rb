node[:machines].each do |machine|
  execute "Power on: #{machine[:name]}" do
    command "vim-cmd vmsvc/power.on #{machine[:id]}"
    not_if "vim-cmd vmsvc/power.getstate #{machine[:id]} | grep 'Powered on'"
  end
end
