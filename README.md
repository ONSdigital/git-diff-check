Git Diff Check
==============

A simple library for checking git diff output for potentially sensitive information

Pre-commit hook
===============

A pre-commit hook script is provided for convenience that uses this library
to test changes before you commit.

## Installing

(requires Go 1.8+)

```shell
$ git clone https://github.com/ONSdigital/git-diff-check.git
$ cd git-diff-check
$ go build -o pre-commit cmd/pre-commit/main.go
$ mkdir -p ${HOME}/.git-templates/hooks
$ mv pre-commit ${HOME}/.git-templates/hooks/pre-commit
```

The hook will now be installed into each repository you subsequently create or
clone locally. If you want to add to an existing repository you can copy the 
`pre-commit` binary into `.git/hooks/pre-commit` in the local repository.

Once installed, the binary will run each time you use `git commit`.

If it finds things it thinks could be sensitive it will flag a warning and stop
the commit proceeding, e.g.:

```shell
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

License
=======

Copyright (c) 2017 Crown Copyright (Office for National Statistics)

Released under MIT license, see [LICENSE](LICENSE) for details.