--- Smaller wrapper function to reduce code redundancy when dealing with shell commands.
local function execShellCommand(command, err_msg)
	if not os.execute(command) then
		io.stderr.write("Error: " .. err_msg)
		os.exit(1)
	end
end

--- Check if a file or directory exists in this path
local function exists(file)
	local ok, err, code = os.rename(file, file)
	if not ok then
		if code == 13 then
			--- Permission denied, but it exists
			return true
		end
	end
	return ok, err
end

--- Check if a directory exists in this path
local function doesDirExist(path)
	--- "/" works on both Unix and Windows
	return exists(path .. "/")
end

--- Initialize the git repository
local function gitInit()
	--- Check if git was already initialized
	io.stdout:write("Checking for existing repository...\n")
	if doesDirExist(".git") then
		io.stderr:write("Error: failed to create git repo. Found an existing .git directory.\n")
		os.exit(1)
	end

	--- Initialize git repo
	io.stdout:write("Executing git init...\n")
	execShellCommand("git init .", "failed to initialize git.")
end

--- Add a remote to the git repository
local function gitRemoteAdd()
	io.stdout:write("Remote name (defaults to origin): ")
	local remote_name = io.stdin:read()
	if remote_name == "n" or remote_name == "" then
		remote_name = "origin"
	end

	--- Add an upload url
	io.stdout:write("Remote url: ")
	local remote_url = io.stdin:read()
	if remote_url == "n" or remote_url == "" then
		io.stderr:write("Error: failed to set the git remote url. Remote url may not be empty.\n")
		os.exit(1)
	end
	if not remote_url:match(".git$") then
		remote_url = remote_url .. ".git"
	end

	--- Add the remote to git
	execShellCommand("git remote add " .. remote_name .. " " .. remote_url .. "", "failed to add git remote.")
	return remote_name
end

--- Commit the changes done by this script.
local function gitCommitChanges()
	io.stdout:write("Staging changes...\n")
	execShellCommand("git add .", "failed to stage changes.")

	io.stdout:write("Committing changes...\n")
	execShellCommand(
		"git commit -m 'Initialize project with proji' -m '[proji](https://github.com/nikoksr/proji)'",
		"failed to commit changes."
	)
end

--- Add a tag to the last commit of the current git repo.
local function gitAddTag()
	io.stdout:write("Adding tag v0.1.0 to latest commit...\n")
	execShellCommand("git tag -a v0.1.0 -m 'Version 0.1.0'", "failed to add tag to latest commit.")
end

--- Push the committed changes to a remote repository.
local function gitPushChanges(remote_name)
	io.stdout:write("Pushing the changes to the remote repository...\n")
	execShellCommand(
		"git push ---quiet -u " .. remote_name .. " master",
		"failed to push changes to remote repository."
	)
end

--- Main wrapper
local function main()
	io.stdout:write("Initializing git repository...\n")

	--- Check if shell is usable
	if not os.execute() then
		io.stderr:write("Error: shell is not available.\n")
		os.exit(1)
	end

	--- Initialize the git repository
	gitInit()

	--- Stage and commit the changes
	gitCommitChanges()

	--- Add tag v0.1.0 to latest commit
	gitAddTag()

	--- Check if git remote should be added
	local input
	repeat
		io.stdout:write("Add a git remote? [y/N] ")
		input = string.lower(io.stdin:read())
		if input == "n" or input == "" then
			io.stdout:write("Done...\n")
			os.exit(0)
		end
	until input == "y"

	--- Add a remote name; typically 'origin'
	local remote_name = gitRemoteAdd()

	--- Push the changes to the remote repository
	gitPushChanges(remote_name)
	io.stdout:write("Done...\n")
end

main()
