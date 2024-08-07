# Justfile

# Variables
build_dir := "build"
src_dir := "."

# Define the Go binary name based on the OS
binary_name := if os_family() == "windows" { "aegis.exe" } else { "aegis" }

# Build the project
build:
    mkdir -p {{build_dir}}
    go build -o {{build_dir}}/{{binary_name}} {{src_dir}}
    @echo "Build complete: {{build_dir}}/{{binary_name}}"

# Clean the build directory
clean:
    -rm -rf {{build_dir}}
    -rm *.enc
    -rm *.dec
    @echo "Clean complete"

# Default recipe
default: build

