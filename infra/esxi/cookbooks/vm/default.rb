require 'net/http'
require 'uri'

node.reverse_merge!({
  password: 'ubuntu',
  vmx_root: '/',
  source_vmdk_path: '',
  ip_prefix: '127.0.0',
  gateway: "127.0.0.254",
  restore: false,
  instance: {
    cpu: 2,
    mem: 4096,
    hdd: "8g",
  },
  admins: [],
})

ssh_keys = node[:admins].map do |username|
  Net::HTTP.get(URI.parse("https://github.com/#{username}.keys")).strip.split("\n").map(&:strip)
end.flatten

node[:machines].each do |machine|
  execute "Reload: #{machine[:name]}" do
    command "vim-cmd vmsvc/reload #{machine[:id]}"
    action :nothing
  end

  vm_dir = "#{node[:vmx_root]}/#{machine[:name]}"
  vmx = "#{vm_dir}/#{machine[:name]}.vmx"
  vm_data = {}
  run_command("cat #{vmx} | sort").stdout.strip.split("\n").each do |line|
    parts = line.split(' = ')
    key = parts.shift
    val = parts.join(' = ').gsub(/^"|"$/, '')
    vm_data[key] = val
  end
  vm_ip = "#{node[:ip_prefix]}.#{machine[:name].sub(/^isucon-/, '').to_i}"

  vm_data['guestOS'] = 'ubuntu-64'
  vm_data['numvcpus'] = node[:instance][:cpu]
  vm_data['memSize'] = node[:instance][:mem]

  if node[:restore]
    vm_data['ide1:0.deviceType'] = 'cdrom-image'
    vm_data['ide1:0.fileName'] = './cloud-init.iso'
    vm_data['ide1:0.present'] = 'TRUE'
  else
    vm_data.keys.each do |key|
      vm_data.delete(key) if key.start_with?('ide1:')
    end
  end

  execute "Restore VMDK: #{machine[:name]}" do
    command <<-EOS
    rm #{vm_dir}/#{machine[:name]}*.vmdk
    vmkfstools -i #{node[:source_vmdk_path]} #{vm_dir}/#{machine[:name]}.vmdk
    vmkfstools -X #{node[:instance][:hdd]} #{vm_dir}/#{machine[:name]}.vmdk
    EOS

    notifies :run, "execute[Reload: #{machine[:name]}]"
  end if node[:restore]

  file "#{vm_dir}/user-data" do
    content <<-EOS.rstrip + "\n"
#cloud-config
user: isuadmin
password: #{node[:password]}
chpasswd: {expire: False}
ssh_pwauth: False
ssh_authorized_keys:
#{ssh_keys.sort.map {|key| "  - #{key}" }.join("\n")}
    EOS

    notifies :run, "execute[Cloud Init: #{machine[:name]}]"
  end

  file "#{vm_dir}/meta-data" do
    content <<-EOS.rstrip + "\n"
instance-id: #{machine[:name]}
local-hostname: #{machine[:name]}
network-interfaces: |
  auto ens160
  iface ens160 inet static
  address #{vm_ip}
  gateway #{node[:gateway]}
  dns-nameservers 8.8.8.8 8.8.4.4
    EOS

    notifies :run, "execute[Cloud Init: #{machine[:name]}]"
  end

  execute "Cloud Init: #{machine[:name]}" do
    action :nothing

    command <<-EOS
      genisoimage -output cloud-init.iso -volid cidata -joliet -rock user-data meta-data
    EOS
    cwd vm_dir

    notifies :run, "execute[Reload: #{machine[:name]}]"
  end

  file vmx do
    content vm_data.keys.sort.map { |k| %{#{k} = "#{vm_data[k]}"} }.join("\n") + "\n"

    notifies :run, "execute[Reload: #{machine[:name]}]"
  end
end
