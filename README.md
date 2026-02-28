# Push - Git Repository Multi-Platform Sync Plugin

A [Crow CI](https://crowci.dev) plugin that synchronizes your git repository code to multiple platforms simultaneously. Push your code changes to multiple git hosting services (Codeberg, Gitea, GitHub, etc.) in a single CI/CD workflow.

## Features

- 🚀 **Multi-platform support**: Push code to multiple git repositories at once
- 🔐 **Multiple authentication methods**: Support for SSH keys, HTTPS tokens, and username/password
- 🌍 **Cross-platform compatibility**: Built with pure Go, works on Linux, macOS, and Windows
- 🔗 **Simple configuration**: YAML-based configuration with environment variables
- 📝 **Batch operations**: Push both branches and tags across all platforms

## Installation

This plugin is designed to work with Crow CI. Add it to your `.crow/push.yaml` file:

```yaml
steps:
  push:
    image: codefloe.com/actions/push:v0.0.2
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

Configure the plugin via environment variables or YAML settings. The plugin reads the `PLUGIN_PLATFORMS` environment variable to determine target repositories.

#### Platform Configuration

Each platform entry supports the following fields:

| Field                | Type    | Required | Description                                            |
|----------------------|---------|----------|--------------------------------------------------------|
| `name`               | string  | ✅        | Platform identifier (e.g., "codeberg", "github")       |
| `enabled`            | boolean | ✅        | Enable/disable pushing to this platform                |
| `url`                | string  | ✅        | Git repository URL (HTTPS or SSH format)               |
| `username`           | string  | ❌        | Username for HTTPS authentication                      |
| `password`           | string  | ❌        | Password for HTTPS authentication                      |
| `token`              | string  | ❌        | Access token for HTTPS authentication                  |
| `ssh_key`            | string  | ❌        | SSH private key for SSH authentication                 |
| `ssh_key_passphrase` | string  | ❌        | Passphrase for encrypted SSH key                       |
| `organization`       | string  | ❌        | Organization name (for future repository creation)     |
| `repository`         | string  | ❌        | Repository name (for future repository creation)       |
| `remote_name`        | string  | ❌        | Custom git remote name (default: "push-plugin-{name}") |

### Authentication Methods

Choose one of the following authentication methods per platform:

#### 1. SSH Key (Recommended for CI)

```yaml
platforms:
  - name: codeberg
    enabled: true
    url: ssh://git@codeberg.org/openhub/push.git
    ssh_key: |
      -----BEGIN OPENSSH PRIVATE KEY-----
      MIIEpAIBAAKCAQEA...
      ...
      -----END OPENSSH PRIVATE KEY-----
    remote_name: codeberg
```

Or use Crow CI secrets:

```yaml
platforms:
  - name: codeberg
    enabled: true
    url: ssh://git@codeberg.org/openhub/push.git
    ssh_key:
      from_secret: codeberg_ssh_key
    remote_name: codeberg
```

**Note**: For CI environments, SSH host key verification is disabled automatically. For local development, ensure `~/.ssh/known_hosts` is properly configured.

#### 2. HTTPS with Token

```yaml
platforms:
  - name: github
    enabled: true
    url: https://github.com/openhub/push.git
    username: git
    password:
      from_secret: github_token
    remote_name: github
```

#### 3. HTTPS with Username/Password

```yaml
platforms:
  - name: gitea
    enabled: true
    url: https://gitea.example.com/openhub/push.git
    username: myusername
    password:
      from_secret: gitea_password
    remote_name: gitea
```

### Complete Example

`.crow/push.yaml`:

```yaml
variables:
  platforms: &platforms
    - name: codeberg
      enabled: true
      organization: openhub
      repository: push
      url: ssh://git@codeberg.org/openhub/push.git
      ssh_key:
        from_secret: ssh_key
      remote_name: codeberg
    - name: github
      enabled: true
      url: https://github.com/openhub/push.git
      username: git
      password:
        from_secret: github_token
      remote_name: github

when:
  - event: [push, pull_request, tag]

clone:
  git:
    image: codeberg.org/crow-plugins/clone:1.0.3
    settings:
      tags: true

steps:
  push:
    image: codefloe.com/actions/push:v0.0.2
    settings:
      platforms: *platforms
```

## How It Works

1. **Repository Detection**: Opens the git repository at `CI_WORKSPACE` (set by Crow CI)
2. **Remote Creation**: Creates git remotes for each enabled platform
3. **Authentication**: Sets up the appropriate authentication method (SSH, HTTP Basic Auth)
4. **Push Operation**: Pushes all branches and tags to each platform
5. **Error Handling**: Reports which platforms succeeded and which failed

## Environment Variables

The plugin reads the following Crow CI environment variables:

- `CI_WORKSPACE`: The repository workspace directory (set by Crow CI)
- `PLUGIN_PLATFORMS`: Platform configuration (JSON array)
- `PLUGIN_ACCESS_TOKEN`: Optional access token (for future use)

Additional environment variables for CI:

- `SSH_KNOWN_HOSTS`: SSH host verification (set to empty string to disable in CI)

## Development

### Prerequisites

- Go 1.26 or later
- go-git v5

### Building

```bash
go mod tidy
go build -o build/push .
```

### Testing

```bash
go test ./...
```

### Architecture

- **`main.go`**: Entry point, reads configuration from environment
- **`config.go`**: Configuration structures and validation
- **`plugin.go`**: Core push logic and authentication handling

## Troubleshooting

### SSH Authentication Errors

**Error**: `invalid auth method` or `cannot create known hosts callback`

**Solution**:
- Ensure SSH URL format: `ssh://git@platform.com/org/repo.git` or `git@platform.com:org/repo.git`
- Set `SSH_KNOWN_HOSTS=""` environment variable in CI
- Verify SSH key is valid OpenSSH format

### HTTPS Authentication Errors

**Error**: `authentication required: Credentials are incorrect or have expired`

**Solution**:
- Verify token/password is correct and not expired
- For Codeberg/Forgejo, use format: `username: git`, `password: <token>`
- Ensure token has `write:repository` permission

### Push Failures

**Error**: `already up to date` - This is not an error, it means there are no new changes

**Solution**: Push operations should succeed with this message

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
