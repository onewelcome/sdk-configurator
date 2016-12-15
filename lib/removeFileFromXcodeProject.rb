require 'rubygems'
require 'xcodeproj'
require 'fileutils'

xcodeproj_filepath = ARGV[0]
file_name = ARGV[1]
group_name = ARGV[2]

# Find group
project = Xcodeproj::Project.open(xcodeproj_filepath)
xcodeproj_group = project.main_group[group_name]

# Remove file from group
if xcodeproj_group != nil
  xcodeproj_group.files.find{ |file|
    if file.real_path.to_s.include? file_name
      file.referrers.each do |ref|
        if ref.isa == "PBXBuildFile"
          ref.remove_from_project
        end
      end
      file.remove_from_project
    end
  }
end
project.save