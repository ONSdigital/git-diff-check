#!/usr/bin/env sh

repo='ONSdigital/git-diff-check'

get_latest_release() { # From https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

if [[ "$OSTYPE" == "darwin"* ]]; then
  binary='pre-commit_darwin-amd64'
  target=${HOME}/.githooks
else
  echo "OS '${OSTYPE}' not currently supported by installer - please refer to manual instructions in the README."
  exit 0
fi

release_version=$(get_latest_release ${repo})
release="https://github.com/${repo}/releases/download/${release_version}/${binary}"

# Create the target location if it doesn't already exist
[ ! -d ${target} ] &&
  {
    echo "Creating global hooks folder at ${target}";
    mkdir -p ${target}
  }

# Check if we're up to date
echo "Check for previous versions ..."
[ -f ${target}/pre-commit.version ] &&
  {
    existing=$(cat ${target}/pre-commit.version)
    echo "-- found existing version ${existing}"
    [ "${existing}" = "$release_version" ] &&
      {
        echo "-- already up to date!"
        exit 0
      }
    echo "-- new version available ${release_version}"
  }

# Fetch the tool
echo "Fetching Git Diff precommit hook ${release_version} ..."
echo "-- from ${release} ..."
curl -L --progress-bar ${release} -o "${target}/pre-commit"
chmod +x "${target}/pre-commit"

# Update the git config
echo "Updating git config ..."
git config --global core.hooksPath ${target}

# Store the installed version
echo ${release_version} > ${target}/pre-commit.version

echo "Add done!"
