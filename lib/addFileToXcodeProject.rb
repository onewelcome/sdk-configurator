require 'rubygems'
require 'xcodeproj'
require 'fileutils'

xcodeproj_filepath = ARGV[0]
file_name = ARGV[1]
app_name = ARGV[2]
group_name = ARGV[3]

# Find group
project = Xcodeproj::Project.open(xcodeproj_filepath)
xcodeproj_group = project.main_group[group_name]

# Create group
unless xcodeproj_group
  xcodeproj_group = project.main_group.new_group(group_name, nil)
end

# Add file to group
file_refs = []
file_ref = xcodeproj_group.files.find{|file|file.real_path.to_s==file_name}
unless file_ref
  file_ref = xcodeproj_group.new_file(file_name)
  file_refs<<file_ref
end

# Add file to target
project.targets.each do |target|
  if target.name == app_name
    if group_name == "Resources"
      target.add_resources(file_refs)
    else
      target.add_file_references(file_refs)
    end
  end
end

project.save