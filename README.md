# git-vendor

Relaxed git submodules.

## Install

### Make

```
sudo make install
```

## Config

### .gitvendor

Config files use git config with the same basic format as branches and submodules.

```
[vendor "my_repo"]
  path = vendor/
  url  = git@localhost:my_repo.git
```

