module RecipeHelper
  def include_cookbook(name)
    cookbook, recipe = *name.split('/')
    include_recipe File.join(__dir__, "cookbooks", cookbook, "#{recipe || 'default'}.rb")
  end
end

Itamae::Recipe::EvalContext.include(RecipeHelper)

include_cookbook 'resolv'
include_cookbook 'sshd'
include_cookbook 'isuadmin'
include_cookbook 'user'
include_cookbook 'ruby'
include_cookbook 'golang'
include_cookbook 'nodejs'
include_cookbook 'nginx'
include_cookbook 'mysql'
include_cookbook 'systemd-timesyncd'
include_cookbook 'redis'
include_cookbook 'monitor-tools'
include_cookbook 'benchmarker'
include_cookbook 'webapp'
