#!/bin/bash

# List of target platforms
platforms=("linux/amd64" "windows/amd64" "darwin/amd64") # Add more platforms if needed

# Create the folder to store the binaries if it doesn't exist
directory="bin-out"  # Change this to the folder name you want to create

app_name="cubecobraprinter" # Change this to the name of your application
build_nr=$(git rev-list --count HEAD)

if [ ! -d "$directory" ]; then
    mkdir "$directory"
    echo "Folder '$directory' created."
else
    echo "Folder '$directory' already exists."
fi

# Loop through each platform and build the binary
for platform in "${platforms[@]}"
do
    export GOOS=$(echo $platform | cut -d'/' -f1)
    export GOARCH=$(echo $platform | cut -d'/' -f2)


    output_name="$directory/$app_name-$build_nr-$GOOS-$GOARCH"
    zipFileName="$output_name.zip"

    if [ "$GOOS" == "windows" ]; then
        output_name="$output_name.exe" # Append .exe for Windows
    fi

    go build -o $output_name
    echo "Built $output_name"

    # create a zip file of the binary
    zip -j $zipFileName $output_name
done

echo "Finished building all binaries."
