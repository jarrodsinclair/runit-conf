# runit-conf

A Supervisor inspired configurator for runit.

## Usage inside a Docker container

Once runit-conf has been installed, call the bootstrap script with a configuration
file as the first (and only) command:

```docker
CMD ["runit-bootstrap", "config/runit-conf.toml"]
```
