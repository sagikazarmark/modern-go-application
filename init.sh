#!/bin/bash

function prompt() {
    echo -n -e "\033[1;32m?\033[0m \033[1m$1\033[0m ($2) "
}

function replace() {
    if [[ ! -z "${TEST}" ]]; then
        mkdir -p $(dirname "tmp/inittest/$2")
        sed -E -e "$1" $2 > tmp/inittest/$2
    else
        sed -E -e "$1" $2 > $2.new
        mv -f $2.new $2
    fi
}

function move() {
    if [[ ! -z "${TEST}" ]]; then
        mkdir -p $(dirname "tmp/inittest/$2")
        cp -r "$1" tmp/inittest/$2
    else
        mv $@
    fi
}

function remove() {
    if [[ -z "${TEST}" ]]; then
        rm $@
    fi
}

if [[ ! -z "${TEST}" ]]; then
    mkdir -p tmp/inittest
    echo "." >> tmp/.gitignore
fi

originalProjectName="project"
originalPackageName="github.com/sagikazarmark/modern-go-application"
originalBinaryName="modern-go-application"
originalServiceName="mga"
originalFriendlyServiceName="Modern Go Application"

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

prompt "Service name" ${projectName}
read serviceName
serviceName=$(echo "${serviceName:-${projectName}}" | sed 's/[[:space:]]//g')

defaultFriendlyServiceName=$(echo "${serviceName}" | sed -e 's/-/ /g;' | awk '{for(i=1;i<=NF;i++){ $i=toupper(substr($i,1,1)) substr($i,2) }}1')
prompt "Friendly service name" "${defaultFriendlyServiceName}"
read friendlyServiceName
friendlyServiceName=${friendlyServiceName:-${defaultFriendlyServiceName}}

prompt "Remove init script" "y/N"
read removeInit
removeInit=${removeInit:-n}

# IDE configuration
move .idea/${originalProjectName}.iml .idea/${projectName}.iml
replace "s|.idea/${originalProjectName}.iml|.idea/${projectName}.iml|g" .idea/modules.xml

# Run configurations
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/All_tests.xml
replace 's|name="project"|name="'${projectName}'"|; s|value="\$PROJECT_DIR\$\/cmd\/'${originalBinaryName}'\/"|value="$PROJECT_DIR$/cmd/'${binaryName}'/"|' .idea/runConfigurations/Debug.xml
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/Integration_tests.xml
replace 's|name="project"|name="'${projectName}'"|' .idea/runConfigurations/Tests.xml
replace "s|${originalBinaryName}|${binaryName}|" .vscode/launch.json

# Update variables
replace "s|${originalServiceName}|${serviceName}|; s|${originalFriendlyServiceName}|${friendlyServiceName}|" cmd/${originalBinaryName}/vars.go

# Binary name
move cmd/${originalBinaryName} cmd/${binaryName}

# Makefile
replace "s|^PACKAGE = .*|PACKAGE = ${packageName}|; s|^BUILD_PACKAGE \??= .*|BUILD_PACKAGE = \${PACKAGE}/cmd/${binaryName}|; s|^BINARY_NAME \?= .*|BINARY_NAME \?= ${binaryName}|" Makefile

# Other project files
declare -a files=(".circleci/config.yml" ".gitlab-ci.yml" "CHANGELOG.md" "Dockerfile")
for file in "${files[@]}"; do
    if [[ -f "${file}" ]]; then
        replace "s|${originalPackageName}|${packageName}|" ${file}
    fi
done

# Update source code
find cmd/ -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done
find internal/ -type f | while read file; do replace "s|${originalPackageName}|${packageName}|" "$file"; done

if [[ "${removeInit}" != "n" && "${removeInit}" != "N" ]]; then
    remove "$0"
fi

# Spotguide
replace "/^ *path: src\/.*/d" .banzaicloud/pipeline.yaml
