# Find a root folder in the users Xcode Project called Pods, or make one
require 'rubygems'
require 'xcodeproj'
require 'fileutils'

xcodeproj_filepath = ARGV[0]
target_name = ARGV[1]
file_path = ARGV[2]
group_name = ARGV[3]

file_name = File.basename(file_path)
xcodeproj_folderpath = File.dirname(xcodeproj_filepath)
xcodeproj_group_path = xcodeproj_folderpath+"/"+group_name

#create certificate folder
FileUtils.mkdir_p xcodeproj_group_path

#copy file
if File.exists?(file_path)
	puts file_path+" already exists. Do you want to overwrite? (y/n)"
	answer = STDIN.gets.chomp
	if answer != "y"
		puts file_path+" was not overwritten."
		exit
	end
end
FileUtils.cp(file_path, xcodeproj_group_path)

#create group
project = Xcodeproj::Project.open(xcodeproj_filepath)
xcodeproj_group = project.main_group[group_name]
unless xcodeproj_group
  	xcodeproj_group = project.main_group.new_group(group_name, nil)
end

#add file to group
file_refs = []
file_ref = xcodeproj_group.files.find{|file|file.real_path.to_s==xcodeproj_group_path+'/'+file_name}
unless file_ref
	file_ref = xcodeproj_group.new_file(xcodeproj_group_path+'/'+file_name)
	file_refs<<file_ref
end

#add file to target and resources/compile files
project.targets.each do |target|
	if target.name == target_name
		if File.extname(file_path)==".cer" || File.extname(file_path)==".pem"
			target.add_resources(file_refs)
		else
			target.add_file_references(file_refs)
		end
	end
end

project.save
