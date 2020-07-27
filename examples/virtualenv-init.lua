--- Smaller wrapper function to reduce code redundancy when dealing with shell commands.
local function execShellCommand(command, err_msg)
	if not os.execute(command) then
		io.stderr:write("Error: " .. err_msg)
		os.exit(1)
	end
end

--- Check if os is *nix
local function isOSUnix()
	if package.config:sub(1, 1) == "/" then
		return true
	end
	return false
end

--- Main wrapper function
function main()
	--- Check if shell is usable
	if not os.execute() then
		io.stderr:write("Error: shell is not available.\n")
		os.exit(1)
	end

	io.stdout:write("Initializing virtualenv...\n")
	execShellCommand("virtualenv --quiet .env", "failed to initialize virtualenv")

	io.stdout:write("Adding .env/ directory to .gitignore file...\n")
	execShellCommand("echo .env/ >>.gitignore", "failed to add .env/ directory to .gitignore file")

	--- Activation command differs depending on OS
	local err_msg = "failed to activate virtualenv."
	local command
	if isOSUnix() then
		command = "source .env/bin/activate"
	else
		command = ".env\\Scripts\\activate"
	end

	io.stdout:write("Activating the virtualenv...\n")
	execShellCommand(command, err_msg)

	io.stdout:write("Installing dependencies with pip...\n")
	execShellCommand("pip install --quiet pylint pep8 black pytest", "failed to install dependencies")

	io.stdout:write("Done...\n")
end

main()
