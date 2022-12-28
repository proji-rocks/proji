local function run_command(command)
  local handle = io.popen(command)
  local result = handle:read("*a")
  handle:close()
  return result
end

-- Create a new virtual environment using Python's built-in venv module
run_command("python -m venv venv")

-- Activate the virtual environment
run_command("source venv/bin/activate")

-- Install the required Python packages using pip
run_command("pip install -r requirements.txt")
