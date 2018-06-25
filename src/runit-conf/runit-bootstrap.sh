#!/bin/sh

# usage: runit-bootstrap config.toml

# execute runit-conf with the parsed config file, writing all runit shell
# scripts in the expected directory structure, returning the absolute base path
# for future processes
cmd="runit-conf $1"
path=$($cmd) || exit $?

# save environment variables to file, in order to source into each run script
export > $path/envvars

# execute boot scripts
# (similar to runit stage 1 scripts)
cd /
$path/boot_all

# execute run scripts using runit process management
# (runit stage 2 scripts)
cd /
exec runsvdir -P $path/service
