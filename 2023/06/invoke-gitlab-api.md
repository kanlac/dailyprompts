# 调用 GitLab API

### Q
- 如何在 GitLab CI 中调用 GitLab API？
- 针对查询 tag，如何显示 3 条提交信息？

### A
GitLab API 的 "Repository Tags" 接口只能获取到特定 tag 对应的单个 commit 信息，而不能直接获取到该 tag 下的多个 commit 信息。

如果你想要获取特定 tag 下的多个 commit 信息，你需要先使用 "Repository Tags" 接口获取到该 tag 对应的 commit 的 SHA，然后使用 "Repository Commits" 接口获取到该 commit 之前的 commit 信息。

以下是一个示例脚本，它首先获取特定 tag 对应的 commit 的 SHA，然后获取该 commit 之前的 3 条 commit 信息：

```bash
#!/bin/bash

# Set your project name and namespace
NAMESPACE="your_namespace"
PROJECT_NAME="your_project_name"
TAG_NAME="your_tag_name"

# URL encode the project name
PROJECT_NAME_ENCODED=$(echo -n "$NAMESPACE/$PROJECT_NAME" | jq -sRr @uri)

# Get the tag's commit SHA
TAG_COMMIT_SHA=$(curl --header "PRIVATE-TOKEN: $CI_API_TOKEN" "<https://gitlab.example.com/api/v4/projects/$PROJECT_NAME_ENCODED/repository/tags/$TAG_NAME>" | jq -r '.commit.id')

# Get the last 3 commits before the tag's commit
curl --header "PRIVATE-TOKEN: $CI_API_TOKEN" "<https://gitlab.example.com/api/v4/projects/$PROJECT_NAME_ENCODED/repository/commits?ref_name=$TAG_COMMIT_SHA&per_page=3>" | jq '.[] | {id: .id, title: .title, message: .message}'

```

在这个脚本中，`CI_API_TOKEN` 是你在 GitLab 的 secret variables 中设置的 API Token，`gitlab.example.com` 是你的 GitLab 服务器的地址，`$PROJECT_NAME_ENCODED` 是你的 URL 编码后的项目名称，`$TAG_NAME` 是你的 tag 名。这个脚本将会使用这个 token 来访问 GitLab API，并获取你的 tag 对应的 commit 之前的 3 条 commit 信息。