# Git Diff Check

A simple library for checking git diff output for potentially sensitive information

## Pre-commit hook

*Requires git 2.9+*

A pre-commit hook script is provided for convenience that uses this library
to test changes before you commit.

- [From binary](#from-binary)
- [From source](#from-source)

### From Binary

- **For Mac OS**

1. Run the installer:

```sh
$ curl -L https://raw.githubusercontent.com/ONSdigital/git-diff-check/master/install.sh | sh
```

- **For other platforms**

1. Download the latest [release](https://github.com/ONSdigital/git-diff-check/releases) for your platform
1. Create (if not already) a folder to store global git hooks (e.g. `${HOME}/.githooks`)
1. Unzip the release and place the `pre-commit` script in the global hooks folder (ensure it's executable)
1. Configure git to use the hooks:

```sh
$ git config --global core.hooksPath <path-to-global-hooks-folder>
```

### From Source

(requires Go 1.11+)

```sh
$ go get github.com/ONSdigital/git-diff-check
# or ..
$ cd ${GOPATH}
$ git clone https://github.com/ONSdigital/git-diff-check.git src/github.com/ONSdigital/git-diff-check
```

Then build:

```sh
$ cd ${GOPATH}/src/github.com/ONSdigital/git-diff-check
$ go build -o pre-commit cmd/pre-commit/main.go
```

Then follow the steps in *From Binary (other platforms)* using your compiled binary
in place of a downloaded one

### Usage

Once installed, the binary will run each time you use `git commit`.

If it finds things it thinks could be sensitive it will flag a warning and stop
the commit proceeding, e.g.:

```sh
$ git add questionableCode.py
$ git commit
Running precommit diff check
WARNING! Potential sensitive data found:
Found in (questionableCode.py)
    > [line] Possible AWS Access Key (line 6)

If you're VERY SURE these files are ok, rerun commit with --no-verify
```

**NB** Currently if you update the pre-commit script in your templates, you will
need to manually re-copy it into each repo that uses it.

## Experimental Entropy Checking

By default, the `pre-commit` tool won't use entropy checking on patch strings. If you
wish to enable this functionality, please set the `DC_ENTROPY_EXPERIMENT` environment
variable.

```sh
$ export DC_ENTROPY_EXPERIMENT=1
```

## License

Copyright (c) 2017 Crown Copyright (Office for National Statistics)

Released under MIT license, see [LICENSE](LICENSE) for details.
