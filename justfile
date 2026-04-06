#!/usr/bin/env -S just --justfile

# https://just.systems/man/en/just-scripts.html#just-scripts
# https://just.systems/man/en/settings.html#positional-arguments
# https://just.systems/man/en/prerequisites.html#prerequisites

set shell := ["bash", "-c"]
set windows-shell := ['powershell.exe', '-ExecutionPolicy', 'RemoteSigned', '-Command']
#set quiet := true

# ============================================================================ #

# File extension of any executable.
_g_exe_ext := if os_family() == "windows" { ".exe" } else { "" }

# ============================================================================ #

# Golang package name, the same as in go.mod file in "module" directive
go_mod := "grpc-gateway-rnd"

# Golang compilation tags
go_tags := ""

# Addresses of PRIVATE repositories, that is used in project
#repo_addresses := "<site_1>,<site_2>,<site_3>"
repo_addresses := ""

exe_out := "./bin/upstream" + _g_exe_ext # Output of binary compiled file
exe_go := "go" + _g_exe_ext # Main go binary
exe_protoc := "protoc" + _g_exe_ext # Main protoc binary
exe_curl := "curl" + _g_exe_ext # Main curl binary

# Main golang source file
main_source_file := "./main.go"

# ============================================================================ #

export GONOPROXY := repo_addresses
export GONOSUMDB := repo_addresses
export GOPRIVATE := repo_addresses

# ============================================================================ #

_g_date := if os_family() == "windows" { `Get-Date -Format "yyyyMMddHHmmss"` } else { `date '+%Y%m%d%H%M%S'` }
_g_git_hash := if os_family() == "windows" {
	`$gitCommit = git rev-parse --short HEAD 2>$null;
	if ($gitCommit) { echo $gitCommit } else { echo '0000000' }`
} else {
	`git rev-parse --short HEAD || echo '0000000'`
}
_g_git_tag := if os_family() == "windows" {
	`$gitTag = git describe --tags --abbrev=0 2>$null;
	if ($gitTag) { echo $gitTag } else { echo '' }`
} else {
	`git describe --tags --abbrev=0 2>/dev/null || echo ''`
}

_g_go_const_package := go_mod + "/internal/pkg/constant."
go_flags := "-ldflags=\"" +\
	"-X '" + go_mod + "/internal/core/constant.gCommit=" + _g_git_hash + "' " +\
	"-X '" + go_mod + "/internal/core/constant.gBuildTime=" + _g_date + "' " +\
	"-X '" + go_mod + "/internal/core/constant.gLastTag=" + _g_git_tag + "' " +\
	"-s\""

# ============================================================================ #

env:
	{{exe_go}} env

protogen:
	{{exe_curl}} --create-dirs -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto \
		-o ./api/include/google/api/annotations.proto
	{{exe_curl}} --create-dirs -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto \
		-o ./api/include/google/api/http.proto
	{{exe_curl}} --create-dirs -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto \
		-o ./api/include/google/api/field_behavior.proto
	{{exe_protoc}} \
		-I ./api \
		-I ./api/include \
		--include_imports \
		--go_out=./internal/gen/proto/ \
		--go-grpc_out=./internal/gen/proto/ \
		--descriptor_set_out=./config/envoy/descriptor.pb \
		./api/user/v1/*.proto

codegen:
	{{exe_go}} generate ./...

deps:
	{{exe_go}} mod tidy
	{{exe_go}} mod download

lint:
	golangci-lint run --fix

build: env deps protogen codegen
	{{exe_go}} build -v --tags={{go_tags}} {{go_flags}} -o {{exe_out}} {{main_source_file}}

[default]
run: build
	{{exe_out}}
