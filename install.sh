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

release="https://github.com/${repo}/releases/download/$(get_latest_release ${repo})/${binary}"
echo ${release}

# Create the target location if it doesn't already exist
[ ! -d ${target} ] &&
  {
    echo "Creating global hooks folder at ${target}";
    mkdir -p ${target}
  }

# Fetch the tool
curl -L ${release} -o "${target}/pre-commit"
chmod +x "${target}/pre-commit"

# Update the git config
echo "Updating git config"
git config --global core.hooksPath ${target}

echo "Add done!"
