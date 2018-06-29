# runit-conf

A Supervisor inspired configurator for runit. Useful for automation and control
of multiple applications inside a Docker container.

## Installation

Inside a Dockerfile:

```dockerfile
ENV RC_TGZ runit-conf-linux-<ARCH>.tgz
ENV RC_VER v#.#.#
RUN wget -q https://github.com/jarrodsinclair/runit-conf/releases/download/${RC_VER}/${RC_TGZ} && \
    tar -xzf ${RC_TGZ} --strip-components=1 -C /usr/local/bin/ && \
    rm -rf ${RC_TGZ}
```

## Usage

Once runit-conf has been installed, call the bootstrap script with a configuration
file as the first (and only) command at the end of the Dockerfile:

```docker
CMD ["runit-bootstrap", "runit-conf.toml"]
```

## Configuration

A TOML formatted file is used for configuring the runit services and startup commands,
similar to a Supervisor YAML file. First, provide a target path to write all of
the runit scripts. A list of services is defined with boot and foreground (run)
commands as Linux shell calls. Each service can optionally be tied to an environment
variable that will activate when the value is 1, else will not run. If the conditional
field is not present in the service description, the service will always run.

Sample configuration file:

```toml
[runit]
directory = "/usr/src/runit"

[[service]]
name = "my_optional_service_1"
conditional = "MY_SERVICE_1"
directory = "my_dir_1"
boot = """
command_at_boot_1.1
command_at_boot_1.2
"""
run = "command_persistent"

[[service]]
name = "my_optional_service_2"
conditional = "MY_SERVICE_2"
boot = "command_at_boot_2"
run = "command_persistent_2"

[[service]]
name = "always_run"
run = "foreground_service"
```
