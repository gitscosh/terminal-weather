# Setup name variables for the package/tool
NAME := terminal-weather
PKG := github.com/gitscosh/$(NAME)

CGO_ENABLED := 0

# Set any default go build tags.
BUILDTAGS :=

include basic.mk

.PHONY: prebuild
prebuild:
