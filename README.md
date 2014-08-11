# git-vendor

Relaxed git submodules.

## Install

### Binaries

TODO:

### Make

```
go get github.com/shanna/git-vendor
```

## Config

### .gitvendor

Config files use git config with the same basic format as branches and submodules.

```
[vendor "my_repo"]
  path   = vendor/my_repo
  url    = git@localhost:my_repo.git
  branch = master
```

