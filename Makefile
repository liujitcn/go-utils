# 运行要求：Linux/macOS，或 Windows 下使用 WSL/Git Bash（需具备 make、python3、go）

.PHONY: help tag

# 统一打 tag：默认扫描根目录及子目录的 go.mod；可通过 MODULE=auth 指定起始目录递归扫描（不提交代码）
tag:
	@python3 scripts/tag_release.py $(if $(MODULE),--path $(MODULE),)

# 显示帮助
help:
	@echo ""
	@echo "Usage:"
	@echo " make [target]"
	@echo ""
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
