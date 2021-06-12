module RecipeHelper
  def include_cookbook(name)
    cookbook, recipe = *name.split('/')
    include_recipe File.join(__dir__, "cookbooks", cookbook, "#{recipe || 'default'}.rb")
  end
end
Itamae::Recipe::EvalContext.include(RecipeHelper)

machines = run_command('vim-cmd vmsvc/getallvms | tail -n +2').stdout.strip.split("\n").map do |vm|
  parts = vm.split(/\s{2,}/)
  { id: parts[0], name: parts[1] }
end.sort_by { |m| m[:name] }

node.reverse_merge!({
  machines: machines,
});

include_cookbook 'power/shutdown'
include_cookbook 'vm'
include_cookbook 'power/on'
