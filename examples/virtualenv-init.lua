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

	io.stdout:write("Hallo Welt, Proji erstellt grade ein neues Projekt\nIch pinge mal google an..\n")
	execShellCommand("ping google.com", "failed to ping google")

end

main()
