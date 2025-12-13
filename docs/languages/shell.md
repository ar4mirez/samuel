# Shell/Bash Guide

> **Applies to**: Bash 4+, POSIX sh, Zsh, Shell Scripts, CI/CD Pipelines

---

## Core Principles

1. **Fail Fast**: Use strict mode, check errors
2. **Portability**: Prefer POSIX when possible
3. **Security**: Quote variables, validate input
4. **Readability**: Clear naming, comments, functions
5. **Idempotency**: Scripts should be safe to re-run

---

## Language-Specific Guardrails

### Strict Mode (Always Use)
- ✓ Start scripts with `set -euo pipefail`
- ✓ Use `#!/usr/bin/env bash` for portability
- ✓ Check return codes for critical commands
- ✓ Use `trap` for cleanup on exit

### Quoting
- ✓ Always quote variables: `"$variable"`
- ✓ Use `"${variable}"` in strings
- ✓ Quote command substitutions: `"$(command)"`
- ✓ Use arrays for lists, not space-separated strings

### Style
- ✓ Use snake_case for variables and functions
- ✓ Use SCREAMING_SNAKE_CASE for constants/env vars
- ✓ 2-space indentation
- ✓ Max line length: 80-100 characters
- ✓ Use `[[` instead of `[` in Bash
- ✓ Use `(( ))` for arithmetic

### Functions
- ✓ Declare functions with `function_name() { }`
- ✓ Use `local` for function variables
- ✓ Document with comments above function
- ✓ Return meaningful exit codes

### Security
- ✓ Validate all user input
- ✓ Don't use `eval` unless absolutely necessary
- ✓ Quote all variables (prevents word splitting)
- ✓ Use full paths for commands in cron/scripts
- ✓ Don't store secrets in scripts

---

## Script Template

```bash
#!/usr/bin/env bash
#
# Description: Brief description of what this script does
# Usage: ./script.sh [options] <arguments>
#
# Author: Your Name
# Date: 2024-01-01

set -euo pipefail
IFS=$'\n\t'

# Constants
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_NAME="$(basename "${BASH_SOURCE[0]}")"

# Default values
DEBUG="${DEBUG:-false}"
VERBOSE="${VERBOSE:-false}"

# Colors (optional)
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly NC='\033[0m' # No Color

#######################################
# Print error message to stderr
# Arguments:
#   Message to print
#######################################
error() {
    echo -e "${RED}ERROR: $*${NC}" >&2
}

#######################################
# Print info message
# Arguments:
#   Message to print
#######################################
info() {
    echo -e "${GREEN}INFO: $*${NC}"
}

#######################################
# Print warning message
# Arguments:
#   Message to print
#######################################
warn() {
    echo -e "${YELLOW}WARN: $*${NC}" >&2
}

#######################################
# Print debug message if DEBUG is true
# Arguments:
#   Message to print
#######################################
debug() {
    if [[ "$DEBUG" == "true" ]]; then
        echo "DEBUG: $*" >&2
    fi
}

#######################################
# Print usage information
#######################################
usage() {
    cat << EOF
Usage: ${SCRIPT_NAME} [OPTIONS] <argument>

Description:
    Brief description of what the script does.

Options:
    -h, --help      Show this help message
    -v, --verbose   Enable verbose output
    -d, --debug     Enable debug output

Arguments:
    argument        Description of the argument

Examples:
    ${SCRIPT_NAME} -v myarg
    ${SCRIPT_NAME} --debug myarg

EOF
}

#######################################
# Cleanup function called on exit
#######################################
cleanup() {
    local exit_code=$?
    debug "Cleaning up..."
    # Add cleanup logic here (remove temp files, etc.)
    exit "$exit_code"
}

#######################################
# Main function
# Arguments:
#   All command line arguments
#######################################
main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                usage
                exit 0
                ;;
            -v|--verbose)
                VERBOSE="true"
                shift
                ;;
            -d|--debug)
                DEBUG="true"
                shift
                ;;
            -*)
                error "Unknown option: $1"
                usage
                exit 1
                ;;
            *)
                break
                ;;
        esac
    done

    # Validate required arguments
    if [[ $# -lt 1 ]]; then
        error "Missing required argument"
        usage
        exit 1
    fi

    local arg="$1"
    debug "Argument: $arg"

    # Main script logic here
    info "Processing: $arg"
}

# Set trap for cleanup
trap cleanup EXIT

# Run main function with all arguments
main "$@"
```

---

## Variables and Data Types

### Variable Declaration
```bash
# Simple assignment (no spaces around =)
name="John"
count=42

# Read-only constant
readonly CONFIG_FILE="/etc/app/config"

# Export for child processes
export PATH="/usr/local/bin:$PATH"

# Default value if not set
name="${NAME:-default_value}"

# Error if not set
required_var="${REQUIRED_VAR:?Error: REQUIRED_VAR not set}"

# Command substitution
current_date="$(date +%Y-%m-%d)"
file_count="$(ls -1 | wc -l)"
```

### Arrays
```bash
# Declare array
declare -a fruits=("apple" "banana" "cherry")

# Add element
fruits+=("date")

# Access element
echo "${fruits[0]}"  # apple

# All elements
echo "${fruits[@]}"  # apple banana cherry date

# Length
echo "${#fruits[@]}"  # 4

# Iterate
for fruit in "${fruits[@]}"; do
    echo "$fruit"
done

# Associative arrays (Bash 4+)
declare -A user
user[name]="John"
user[email]="john@example.com"
echo "${user[name]}"
```

### String Operations
```bash
string="Hello, World!"

# Length
echo "${#string}"  # 13

# Substring
echo "${string:0:5}"  # Hello

# Replace first occurrence
echo "${string/World/Bash}"  # Hello, Bash!

# Replace all occurrences
echo "${string//l/L}"  # HeLLo, WorLd!

# Remove prefix
filename="document.txt"
echo "${filename%.txt}"  # document

# Remove suffix
path="/home/user/file.txt"
echo "${path##*/}"  # file.txt (basename)
echo "${path%/*}"   # /home/user (dirname)

# Uppercase/lowercase (Bash 4+)
echo "${string^^}"  # HELLO, WORLD!
echo "${string,,}"  # hello, world!
```

---

## Control Flow

### Conditionals
```bash
# If statement
if [[ "$name" == "John" ]]; then
    echo "Hello, John"
elif [[ "$name" == "Jane" ]]; then
    echo "Hello, Jane"
else
    echo "Hello, stranger"
fi

# Test operators
[[ -z "$var" ]]      # True if empty
[[ -n "$var" ]]      # True if not empty
[[ "$a" == "$b" ]]   # String equality
[[ "$a" != "$b" ]]   # String inequality
[[ "$a" =~ ^[0-9]+$ ]]  # Regex match

# Numeric comparison (use (( )) or -eq, -lt, etc.)
if (( count > 10 )); then
    echo "Count is greater than 10"
fi

if [[ "$count" -gt 10 ]]; then
    echo "Count is greater than 10"
fi

# File tests
[[ -f "$file" ]]     # True if file exists
[[ -d "$dir" ]]      # True if directory exists
[[ -r "$file" ]]     # True if readable
[[ -w "$file" ]]     # True if writable
[[ -x "$file" ]]     # True if executable
[[ -s "$file" ]]     # True if file size > 0
```

### Loops
```bash
# For loop
for i in 1 2 3 4 5; do
    echo "$i"
done

# C-style for loop
for ((i = 0; i < 10; i++)); do
    echo "$i"
done

# Loop over array
for item in "${array[@]}"; do
    echo "$item"
done

# Loop over files
for file in *.txt; do
    [[ -f "$file" ]] || continue
    echo "Processing: $file"
done

# While loop
count=0
while [[ $count -lt 5 ]]; do
    echo "$count"
    ((count++))
done

# Read file line by line
while IFS= read -r line; do
    echo "$line"
done < "$file"

# Until loop
until [[ -f "$file" ]]; do
    echo "Waiting for $file..."
    sleep 1
done
```

### Case Statement
```bash
case "$command" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        stop_service
        start_service
        ;;
    status|info)  # Multiple patterns
        show_status
        ;;
    *)
        echo "Unknown command: $command"
        exit 1
        ;;
esac
```

---

## Functions

### Function Definition
```bash
#######################################
# Calculate sum of numbers
# Arguments:
#   Numbers to sum
# Returns:
#   Sum via stdout
#######################################
calculate_sum() {
    local sum=0
    local num

    for num in "$@"; do
        ((sum += num))
    done

    echo "$sum"
}

# Call function
result="$(calculate_sum 1 2 3 4 5)"
echo "Sum: $result"
```

### Return Values
```bash
# Return exit code (0 = success, 1-255 = error)
is_valid_email() {
    local email="$1"

    if [[ "$email" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
        return 0
    else
        return 1
    fi
}

if is_valid_email "test@example.com"; then
    echo "Valid email"
fi

# Return value via stdout
get_config_value() {
    local key="$1"
    local config_file="${2:-/etc/app/config}"

    grep "^${key}=" "$config_file" | cut -d'=' -f2
}

value="$(get_config_value "database_host")"
```

---

## Input/Output

### Reading Input
```bash
# Read single line
read -r -p "Enter your name: " name

# Read with timeout
if read -r -t 10 -p "Enter value (10s timeout): " value; then
    echo "You entered: $value"
else
    echo "Timeout!"
fi

# Read password (no echo)
read -r -s -p "Enter password: " password
echo  # New line after password

# Read into array
read -r -a items <<< "one two three"
```

### Output
```bash
# Stdout
echo "Normal output"
printf "Formatted: %s - %d\n" "$name" "$count"

# Stderr
echo "Error message" >&2

# Redirect stdout to file
echo "content" > file.txt    # Overwrite
echo "content" >> file.txt   # Append

# Redirect both stdout and stderr
command > output.log 2>&1
command &> output.log  # Bash shorthand

# Discard output
command > /dev/null 2>&1
```

### Here Documents
```bash
# Multi-line string
cat << EOF
This is a multi-line
string with variable expansion: $variable
EOF

# No variable expansion
cat << 'EOF'
This preserves $literal text
EOF

# Indent-friendly (<<-)
cat <<- EOF
	This ignores leading tabs
	for cleaner code
EOF
```

---

## Error Handling

### Exit Codes
```bash
# Check command success
if command; then
    echo "Success"
else
    echo "Failed with exit code: $?"
fi

# Chain commands
command1 && command2  # Run command2 only if command1 succeeds
command1 || command2  # Run command2 only if command1 fails

# Exit on error
set -e
command_that_might_fail || true  # Continue even if fails

# Custom error handling
set -e
trap 'echo "Error on line $LINENO"; exit 1' ERR
```

### Trap
```bash
# Cleanup on exit
cleanup() {
    rm -f "$temp_file"
    echo "Cleaned up"
}
trap cleanup EXIT

# Handle signals
trap 'echo "Interrupted"; exit 130' INT
trap 'echo "Terminated"; exit 143' TERM

# Multiple signals
trap 'cleanup' EXIT INT TERM

# Ignore signal
trap '' SIGINT
```

---

## File Operations

### Common Operations
```bash
# Check existence
[[ -f "$file" ]] && echo "File exists"
[[ -d "$dir" ]] && echo "Directory exists"

# Create directory (with parents)
mkdir -p "$dir/subdir"

# Copy
cp "$source" "$dest"
cp -r "$source_dir" "$dest_dir"

# Move/rename
mv "$old_name" "$new_name"

# Delete
rm "$file"
rm -rf "$dir"  # Careful!

# Temporary file
temp_file="$(mktemp)"
trap 'rm -f "$temp_file"' EXIT
```

### Text Processing
```bash
# Search
grep "pattern" file.txt
grep -r "pattern" directory/
grep -E "regex" file.txt

# Replace
sed 's/old/new/g' file.txt
sed -i 's/old/new/g' file.txt  # In-place

# Extract columns
cut -d',' -f1,3 file.csv
awk -F',' '{print $1, $3}' file.csv

# Sort and unique
sort file.txt
sort -u file.txt  # Unique
sort -n file.txt  # Numeric

# Count
wc -l file.txt  # Lines
wc -w file.txt  # Words
```

---

## Best Practices

### Do This
```bash
# Quote variables
echo "$variable"

# Use [[ ]] for tests
if [[ -f "$file" ]]; then

# Use arrays for lists
files=("file1.txt" "file 2.txt" "file3.txt")
for f in "${files[@]}"; do

# Use local in functions
my_func() {
    local var="value"
}

# Check command exists
if command -v docker &> /dev/null; then
    docker run ...
fi

# Use parameter expansion defaults
name="${1:-default}"
```

### Don't Do This
```bash
# Unquoted variables
echo $variable  # Word splitting issues

# [ ] instead of [[ ]]
if [ -f $file ]; then  # Breaks with spaces

# Space-separated "arrays"
files="file1 file2 file3"
for f in $files; do  # Breaks with spaces in names

# Global variables in functions
my_func() {
    var="value"  # Pollutes global scope
}

# Parsing ls output
for f in $(ls); do  # Don't do this
```

---

## Common Patterns

### Argument Parsing with getopts
```bash
while getopts ":hv:o:" opt; do
    case $opt in
        h)
            usage
            exit 0
            ;;
        v)
            verbose="$OPTARG"
            ;;
        o)
            output="$OPTARG"
            ;;
        \?)
            error "Invalid option: -$OPTARG"
            exit 1
            ;;
        :)
            error "Option -$OPTARG requires an argument"
            exit 1
            ;;
    esac
done
shift $((OPTIND - 1))
```

### Configuration File
```bash
# Load config file
if [[ -f "$config_file" ]]; then
    # shellcheck source=/dev/null
    source "$config_file"
fi

# Or parse key=value format
while IFS='=' read -r key value; do
    [[ "$key" =~ ^#.*$ ]] && continue  # Skip comments
    [[ -z "$key" ]] && continue         # Skip empty lines
    declare "$key=$value"
done < "$config_file"
```

### Logging
```bash
readonly LOG_FILE="/var/log/myapp.log"

log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp
    timestamp="$(date '+%Y-%m-%d %H:%M:%S')"

    echo "[$timestamp] [$level] $message" >> "$LOG_FILE"

    if [[ "$VERBOSE" == "true" ]]; then
        echo "[$level] $message"
    fi
}

log "INFO" "Starting application"
log "ERROR" "Something went wrong"
```

---

## Testing

### ShellCheck
```bash
# Install: apt install shellcheck / brew install shellcheck

# Run on script
shellcheck script.sh

# Disable specific check
# shellcheck disable=SC2034
unused_variable="value"
```

### Unit Testing with Bats
```bash
#!/usr/bin/env bats

setup() {
    # Run before each test
    source ./my_script.sh
}

@test "calculate_sum adds numbers correctly" {
    result="$(calculate_sum 1 2 3)"
    [ "$result" -eq 6 ]
}

@test "is_valid_email validates correct email" {
    run is_valid_email "test@example.com"
    [ "$status" -eq 0 ]
}

@test "is_valid_email rejects invalid email" {
    run is_valid_email "invalid"
    [ "$status" -eq 1 ]
}
```

---

## References

- [Bash Manual](https://www.gnu.org/software/bash/manual/)
- [ShellCheck](https://www.shellcheck.net/)
- [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html)
- [Pure Bash Bible](https://github.com/dylanaraps/pure-bash-bible)
- [Bats Testing](https://github.com/bats-core/bats-core)
