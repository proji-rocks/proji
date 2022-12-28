local function run_command(command)
  local handle = io.popen(command)
  local result = handle:read("*a")
  handle:close()
  return result
end

-- Initialize a new Git repository
run_command("git init")

-- Add all the files in the current directory to the staging area
run_command("git add .")

-- Commit the staged changes
run_command("git commit -m 'Initial commit'")
