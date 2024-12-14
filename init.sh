#!/bin/bash

# Destination directory of modifications
DEST="."

# Original project variables
originalProjectName="project"
originalPackageName="github.com/sagikazarmark/modern-go-application"
originalBinaryName="modern-go-application"
originalAppName="mga"
originalFriendlyAppName="Modern Go Application"

# Prepare testing
if [[ ! -z "${TEST}" ]]; then
    #set -xe
    DEST="tmp/inittest"
    mkdir -p ${DEST}
    echo "." > tmp/.gitignore
fi

function prompt() {
    echo -n -e "\033[1;32m?\033[0m \033[1m$1\033[0m ($2) "
}

function replace() {
    if [[ ! -z "${TEST}" ]]; then
        dest=$(echo $2 | sed "s|^${DEST}/||")
        mkdir -p $(dirname "${DEST}/${dest}")
        if [[ "$2" == "${DEST}/${dest}" ]]; then
            sed -E -e "$1" $2 > ${DEST}/${dest}.new
            mv -f ${DEST}/${dest}.new ${DEST}/${dest}
        else
            sed -E -e "$1" $2 > ${DEST}/${dest}
        fi
    else
        if [[ -f "$2" ]]; then  # Only try to replace if file exists
            sed -E -e "$1" $2 > $2.new
            mv -f $2.new $2
        fi
    fi
}

function move() {
    if [[ ! -z "${TEST}" ]]; then
        dest=$(echo $2 | sed "s|^${DEST}/||")
        mkdir -p $(dirname "${DEST}/${dest}")
        cp -r "$1" ${DEST}/${dest}
    else
        if [[ -e "$1" ]]; then  # Only try to move if source exists
            mkdir -p $(dirname "$2")  # Create target directory if it doesn't exist
            mv $@
        fi
    fi
}

function remove() {
    if [[ -z "${TEST}" ]]; then
        if [[ -e "$1" ]]; then  # Only try to remove if file exists
            rm $@
        fi
    fi
}

defaultPackageName=${PWD##*src/}
prompt "Package name" ${defaultPackageName}
read packageName
packageName=$(echo "${packageName:-${defaultPackageName}}" | sed 's/[[:space:]]//g')

defaultProjectName=$(basename ${packageName})
prompt "Project name" ${defaultProjectName}
read projectName
projectName=$(echo "${projectName:-${defaultProjectName}}" | sed 's/[[:space:]]//g')

prompt "Binary name" ${projectName}
read binaryName
binaryName=$(echo "${binaryName:-${projectName}}" | sed 's/[[:space:]]//g')

prompt "Application name" ${projectName}
read appName
appName=$(echo "${appName:-${projectName}}" | sed 's/[[:space:]]//g')

defaultFriendlyAppName=$(echo "${appName}" | sed -e 's/-/ /g;' | awk '{for(i=1;i<=NF;i++){ $i=toupper(substr($i,1,1)) substr($i,2) }}1')
prompt "Friendly application name" "${defaultFriendlyAppName}"
read friendlyAppName
friendlyAppName=${friendlyAppName:-${defaultFriendlyAppName}}

prompt "Update README" "Y/n"
read updateReadme
updateReadme=${updateReadme:-y}

prompt "Remove init script" "y/N"
read removeInit
removeInit=${removeInit:-n}

# Create necessary directories
mkdir -p .idea/runConfigurations
mkdir -p .vscode

# IDE configuration
if [[ -f ".idea/${originalProjectName}.iml" ]]; then
    move ".idea/${originalProjectName}.iml" ".idea/${projectName}.iml"
    replace "s|${originalProjectName}.iml|${projectName}.iml|g" .idea/modules.xml
fi

# Run configurations
if [[ -f ".idea/runConfigurations/All_tests.xml" ]]; then
    replace 's|name="'${originalProjectName}'"|name="'${projectName}'"|; s|module name="'${originalProjectName}'"|module name="'${projectName}'"|; s|value="'${originalPackageName}'"|value="'${packageName}'"|' .idea/runConfigurations/All_tests.xml
fi

if [[ -f ".idea/runConfigurations/Debug.xml" ]]; then
    replace 's|name="'${originalProjectName}'"|name="'${projectName}'"|; s|module name="'${originalProjectName}'"|module name="'${projectName}'"|; s|value="'${originalPackageName}'/cmd/'${originalBinaryName}'"|value="'${packageName}'/cmd/'${binaryName}'"|' .idea/runConfigurations/Debug.xml
fi

if [[ -f ".idea/runConfigurations/Integration_tests.xml" ]]; then
    replace 's|name="'${originalProjectName}'"|name="'${projectName}'"|; s|module name="'${originalProjectName}'"|module name="'${projectName}'"|; s|value="'${originalPackageName}'"|value="'${packageName}'"|' .idea/runConfigurations/Integration_tests.xml
fi

if [[ -f ".idea/runConfigurations/Tests.xml" ]]; then
    replace 's|name="'${originalProjectName}'"|name="'${projectName}'"|; s|module name="'${originalProjectName}'"|module name="'${projectName}'"|; s|value="'${originalPackageName}'"|value="'${packageName}'"|' .idea/runConfigurations/Tests.xml
fi

# VSCode configuration
if [[ -f ".vscode/launch.json" ]]; then
    replace "s|${originalBinaryName}|${binaryName}|" .vscode/launch.json
fi

# Binary changes:
#   - binary name
#   - source code
#   - variables
if [[ -d "cmd/${originalBinaryName}" ]]; then
    move cmd/${originalBinaryName} cmd/${binaryName}
fi

if [[ -d "cmd/${binaryName}" ]]; then
    replace "s|${originalAppName}|${appName}|; s|${originalFriendlyAppName}|${friendlyAppName}|" ${DEST}/cmd/${binaryName}/main.go
    find ${DEST}/cmd -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done
fi

# Other project files
declare -a files=("CHANGELOG.md" "prototool.yaml" "go.mod" ".golangci.yml" "gqlgen.yml")
for file in "${files[@]}"; do
    if [[ -f "${file}" ]]; then
        replace "s|${originalPackageName}|${packageName}|" ${file}
    fi
done
declare -a files=("prototool.yaml")
for file in "${files[@]}"; do
    if [[ -f "${file}" ]]; then
        replace "s|${originalProjectName}|${projectName}|" ${file}
    fi
done
declare -a files=("Dockerfile")
for file in "${files[@]}"; do
    if [[ -f "${file}" ]]; then
        replace "s|${originalBinaryName}|${binaryName}|" ${file}
    fi
done

# Update source code
if [[ -d ".gen" ]]; then  # Only try to update .gen if it exists
    find .gen -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done
fi
if [[ -d "internal" ]]; then  # Only try to update internal if it exists
    find internal -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done
fi

if [[ "${removeInit}" != "n" && "${removeInit}" != "N" ]]; then
    remove "$0"
fi

# Update readme
if [[ "${updateReadme}" == "y" || "${updateReadme}" == "Y" ]]; then
    echo -e "# FRIENDLY_PROJECT_NAME\n\n**Project description.**" | sed "s/FRIENDLY_PROJECT_NAME/${friendlyAppName}/" > ${DEST}/README.md
fi
