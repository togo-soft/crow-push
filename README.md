# Push - Git Repository Multi-Platform Sync Plugin

A [Crow CI](https://crowci.dev) plugin that synchronizes your git repository code to multiple platforms simultaneously. Push your code changes to multiple git hosting services (GitHub, Gitea, Codeberg, Gitee, Bitbucket, CNB, etc.) in a single CI/CD workflow.

## Features

- 🚀 **Multi-platform support**: Push code to multiple git repositories at once (GitHub, Gitea, Codeberg, Gitee, Bitbucket, CNB)
- 🤖 **Auto repository creation**: Automatically create repositories on target platforms if they don't exist
- 🔐 **Multiple authentication methods**: Support for SSH keys, HTTPS tokens, and username/password
- 🌍 **Cross-platform compatibility**: Built with pure Go, works on Linux, macOS, and Windows
- 🔗 **Simple configuration**: YAML-based configuration with environment variables
- 📝 **Batch operations**: Push both branches and tags across all platforms

## Supported Platforms

- **GitHub** - `platform_type: github`
- **Gitea**/**Codeberg**/**Codefloe** - `platform_type: gitea`
- **gitee.com** - `platform_type: gitee`
- **bitbucket.org** - `platform_type: bitbucket`
- **cnb.cool** - `platform_type: cnb`

## Installation

This plugin is designed to work with Crow CI. Add it to your `.crow/push.yaml` file:

```yaml
steps:
  push:
    image: codefloe.com/actions/push:v1.0.3
    settings:
      platforms:
        - name: codeberg
          enabled: true
          url: ssh://git@codeberg.org/openhub/push.git
          ssh_key:
            from_secret: ssh_key
          remote_name: codeberg
```

## Usage

### Configuration

Configure the plugin via YAML settings in your Crow CI workflow. The plugin reads the `PLUGIN_PLATFORMS` environment variable to determine target repositories.

#### Platform Configuration

Each platform entry supports the following fields:

| Field                | Type    | Required | Description                                                   |
|----------------------|---------|----------|---------------------------------------------------------------|
| `name`               | string  | ✅        | Platform identifier (e.g., "github", "codeberg")              |
| `enabled`            | boolean | ✅        | Enable/disable pushing to this platform                       |
| `url`                | string  | ✅        | Git repository URL (HTTPS or SSH format)                      |
| `platform_type`      | string  | ❌        | Platform type: `github`, `gitea`, `gitee`, `bitbucket`, `cnb` |
| `username`           | string  | ❌        | Username for HTTPS authentication                             |
| `password`           | string  | ❌        | Password for HTTPS authentication                             |
| `token`              | string  | ❌        | Access token for authentication                               |
| `ssh_key`            | string  | ❌        | SSH private key for SSH authentication                        |
| `ssh_key_passphrase` | string  | ❌        | Passphrase for encrypted SSH key                              |
| `ssh_user`           | string  | ❌        | SSH username (default: "git")                                 |
| `owner`              | string  | ❌        | Repository owner/organization (for auto-creation)             |
| `repository`         | string  | ❌        | Repository name (for auto-creation)                           |
| `is_organization`    | boolean | ❌        | Whether owner is an organization (for auto-creation)          |
| `auto_create`        | boolean | ❌        | Auto-create repository if it doesn't exist                    |
| `is_private`         | boolean | ❌        | Create as private repository (for auto-creation)              |
| `endpoint`           | string  | ❌        | API endpoint for Gitea-like platforms                         |
| `remote_name`        | string  | ❌        | Custom git remote name (default: "push-plugin-{name}")        |

### Authentication Methods

Choose one of the following authentication methods per platform:

#### 1. SSH Key (Recommended for CI)

```yaml
platforms:
  - name: github
    enabled: true
    url: ssh://git@github.com/owner/repo.git
    ssh_key:
      from_secret: ssh_key
    ssh_user: git  # Optional, defaults to "git"
    remote_name: github
```

#### 2. HTTPS with Token

**NOTE**: Most platforms have already blocked HTTPS push code.

```yaml
platforms:
  - name: github
    enabled: true
    url: https://github.com/owner/repo.git
    username: xxxx
    token:
      from_secret: github_token
    remote_name: github
```

#### 3. HTTPS with Username/Password

**NOTE**: Most platforms have already blocked HTTPS push code.

```yaml
platforms:
  - name: gitee
    enabled: true
    url: https://gitee.com/owner/repo.git
    username: myusername
    password:
      from_secret: gitee_password
    remote_name: gitee
```

### Auto-Create Repository

Enable automatic repository creation with the `auto_create` option:

```yaml
platforms:
  - name: github
    enabled: true
    owner: my-organization
    repository: my-repo
    url: ssh://git@github.com/my-organization/my-repo.git
    ssh_key:
      from_secret: ssh_key
    platform_type: github
    auto_create: true
    is_organization: true
    is_private: false
    remote_name: github
```

**Requirements for auto-creation:**
- `owner`: Repository owner or organization name
- `repository`: Repository name to create
- `platform_type`: Must be specified (github, gitea, gitee, bitbucket, cnb)
- `is_organization`: Set to `true` if owner is an organization
- `token`: Required, must have repository creation permission

### Complete Multi-Platform Example

```yaml
variables:
  platforms: &platforms
    - name: github
      enabled: true
      owner: my-org
      repository: my-repo
      url: ssh://git@github.com/my-org/my-repo.git
      ssh_key:
        from_secret: ssh_key
      platform_type: github
      auto_create: true
      is_organization: true
      is_private: false
      remote_name: github

    - name: codeberg
      enabled: true
      owner: my-org
      repository: my-repo
      url: ssh://git@codeberg.org/my-org/my-repo.git
      ssh_key:
        from_secret: ssh_key
      platform_type: gitea
      endpoint: https://codeberg.org
      auto_create: true
      is_organization: true
      is_private: false
      remote_name: codeberg

    - name: gitee
      enabled: true
      owner: my-user
      repository: my-repo
      url: https://gitee.com/my-user/my-repo.git
      username: my-username
      password:
        from_secret: gitee_password
      platform_type: gitee
      auto_create: true
      is_private: false
      remote_name: gitee

when:
  - event: [push, pull_request, tag]

clone:
  git:
    image: codeberg.org/crow-plugins/clone:1.0.3
    settings:
      tags: true

steps:
  push:
    image: codefloe.com/actions/push:v1.0.3
    settings:
      PLATFORMS: *platforms
```

## How It Works

1. **Configuration Loading**: Reads platform configuration from `PLUGIN_PLATFORMS` environment variable
2. **Repository Opening**: Opens the git repository at `CI_WORKSPACE` (set by Crow CI)
3. **Auto-Create**: If `auto_create` is enabled, creates repositories on target platforms using their APIs
4. **Remote Creation**: Creates git remotes for each enabled platform
5. **Authentication**: Sets up the appropriate authentication method (SSH, HTTP Basic Auth)
6. **Push Operation**: Pushes all branches and tags to each platform with `--force` flag
7. **Error Reporting**: Reports which platforms succeeded and which failed

## Environment Variables

The plugin reads the following Crow CI environment variables:

- `CI_WORKSPACE`: The repository workspace directory (set by Crow CI)
- `PLUGIN_PLATFORMS`: Platform configuration (JSON array)
- `PLUGIN_ACCESS_TOKEN`: Optional access token (for future use)

## Development

### Prerequisites

- Go 1.26 or later
- go-git v5
- Client libraries for each platform API (github.com/google/go-github, etc.)

### Building

```bash
go mod tidy
go build -o build/push .
```

### Architecture

- **`main.go`**: Entry point, reads configuration from environment
- **`config.go`**: Configuration structures and validation
- **`plugin.go`**: Core push logic, authentication handling, and repository auto-creation
- **`client/`**: Platform-specific API clients for repository creation
  - `client.go`: Client interface definition
  - `github.go`: GitHub API client
  - `gitea.go`: Gitea API client
  - `gitee.go`: Gitee API client
  - `bitbucket.go`: Bitbucket API client
  - `cnb.go`: CNB API client

## Troubleshooting

### SSH Authentication Errors

**Error**: `invalid auth method` or `cannot create known hosts callback`

**Solution**:
- Ensure SSH URL format: `ssh://git@platform.com/org/repo.git` or `git@platform.com:org/repo.git`
- Verify SSH key is valid OpenSSH format
- Check `ssh_user` is correct (defaults to "git")

### HTTPS Authentication Errors

**Error**: `authentication required: Credentials are incorrect or have expired`

**Solution**:
- Verify token/password is correct and not expired
- Ensure token has `write:repository` permission
- Check that the URL format is correct (https://.../repo.git)

### Repository Auto-Creation Failures

**Error**: `create repository failed`

**Solution**:
- Verify `platform_type` is correct
- Check that `owner` and `repository` fields are set
- Ensure token has repository creation permission
- For Gitea-like platforms, verify `endpoint` is set correctly
- Verify `is_organization` is set correctly

### Push Failures

**Error**: `already up to date` - This is not an error, it means there are no new changes

**Solution**: Push operations should succeed with this message. All refs are already synchronized.

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
