node[:machines].each do |machine|
  execute "Power off: #{machine[:name]}" do
    command "vim-cmd vmsvc/power.off #{machine[:id]}"
    not_if "vim-cmd vmsvc/power.getstate #{machine[:id]} | grep 'Powered off'"
  end
end
