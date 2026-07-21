# Root Makefile — forwards to tools/cmd
# 内容库为本仓库主体；Go 工具在 tools/cmd 下自包含维护。

.PHONY: build run test test-short test-coverage clean install fmt vet lint lint-fix security quality release mod help

# 所有 Go 相关 target 转发到 tools/cmd
build run test test-short test-coverage clean install fmt vet lint lint-fix security quality release mod help:
	@$(MAKE) -C tools/cmd $@
