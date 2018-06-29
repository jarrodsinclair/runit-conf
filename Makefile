NAME=runit-conf
export GOPATH=$(CURDIR)

# call do_build(1:name, 2:os, 3:arch, 4:arm_ver)
define do_build
	$(eval base := $(1)-$(2)-$(3)$(4))
	@cp -r out/tmpl out/$(base)
	@cd out/$(base); GOOS=$2 GOARCH=$3 GOARM=$4 CGO_ENABLED=0 go build ../../src/$(1)/
	@cd out; tar -zcf $(base).tgz $(base)/*
	@rm -rf out/$(base)
endef

all:
	@rm -rf out; mkdir -p out/tmpl
	@cp src/$(NAME)/runit-bootstrap.sh out/tmpl/runit-bootstrap
	@chmod 0755 out/tmpl/runit-bootstrap
	@cd src/$(NAME); rm -rf Gopkg.lock vendor
	@cd src/$(NAME); dep ensure
	$(call do_build,$(NAME),linux,arm,5)
	$(call do_build,$(NAME),linux,arm,6)
	$(call do_build,$(NAME),linux,arm,7)
	$(call do_build,$(NAME),linux,arm64,)
	$(call do_build,$(NAME),linux,amd64,)
	@rm -rf out/tmpl
	@ls -al out
