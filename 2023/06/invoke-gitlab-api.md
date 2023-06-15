# 调用 GitLab API

### 如何在 GitLab CI 中调用 GitLab API？
1. **获取 API Token**：GitLab API 需要一个有效的 API Token 进行身份验证。你可以使用个人访问令牌（Personal Access Token）或者 CI_JOB_TOKEN。个人访问令牌具有更广泛的权限，可以访问 GitLab 中的大部分 API，而 CI_JOB_TOKEN 的权限较少，只能访问与当前 CI/CD 作业相关的一些 API。（注：实际测试发现使用 CI_JOB_TOKEN 会返回 `{"message":"401 Unauthorized"}`）
    
    如果你需要使用个人访问令牌，你可以在 GitLab 的用户设置中创建一个。然后，你应该将这个令牌作为一个 secret variable 存储在你的 GitLab 项目的 CI/CD 设置中，以保护它的安全。
    
2. 发起 API 请求，在请求头中包含 token。可以用 project id 查，也可以用 group 拼装 project name 的形式。

```bash
# default branch
curl --insecure --header "PRIVATE-TOKEN: v_5Usxzc9FBjxcB7c-Vn" "https://git.i.ncmps.com/api/v4/projects/nlc%2Fnlc_script_api/repository/commits?per_page=3"
# specific branch
curl --insecure --header "PRIVATE-TOKEN: BF-EVJSBCuzT-HxFrC8w" "https://git.i.ncmps.com/api/v4/projects/f5%2Fmain/repository/commits?ref_name=luomeng&per_page=3"

# tag 特殊一点
curl --insecure --header "PRIVATE-TOKEN: v_5Usxzc9FBjxcB7c-Vn" "https://git.i.ncmps.com/api/v4/projects/f5%2Fmain/repository/tags/v1.4.10"
```

### 针对查询 tag，如何显示 3 条提交信息？
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