.PHONY: tag sub-tag

# 根模块：仅根据远程仓库更新状态决定是否打并推送远程 tag（不提交代码）
tag:
	@python3 scripts/tag_release.py tag

# 多模块：递归检查 go.mod 目录，仅根据远程仓库更新状态为模块打并推送远程 tag（不提交代码）
sub-tag:
	@python3 scripts/tag_release.py sub-tag
