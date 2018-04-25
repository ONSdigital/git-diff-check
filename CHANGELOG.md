# Changelog

## Unreleased

- fix #8 ignore lines in a patch that are being removed rather than added / changed

## 0.4.0 2018-04-25

- move to global installation - requires removal of previous individual repo copies
- add install script for mac os

## 0.3.0 2018-04-03

- add build for `linux/amd64`

## 0.2.0 2018-03-28

- add regex line test for `----BEGIN CERTIFICATE----`
- add entropy check into core (optional)
- add `DC_ENTROPY_EXPERIMENT` environment option to activate entropy checking
- add support for multi-os cross compilation in Makefile

## 0.1.0 2018-02-07

- initial release
