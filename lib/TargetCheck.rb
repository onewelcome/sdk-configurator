# Find a root folder in the users Xcode Project called Pods, or make one
require 'rubygems'
require 'xcodeproj'

xcodeproj_filepath = ARGV[0]
target_name = ARGV[1]

project = Xcodeproj::Project.open(xcodeproj_filepath)

target_exists = 0
unless project.nil?
	project.targets.each do |target|
		if target.name == target_name
			target_exists = 1		
		end
	end
end

puts target_exists