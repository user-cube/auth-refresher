# Auth Refresher

Auth Refresher is a command-line tool designed to simplify the process of managing Docker and ECR registry logins. It provides an intuitive interface for selecting registries from a configuration file and handles login operations with support for AWS and Helm registries.

## Features

- **Docker/ECR Registry Login**: Easily log in to Docker and AWS ECR registries.
- **Helm Registry Login**: Seamlessly log in to Helm registries.
- **Graceful Cancellation**: Cancel operations gracefully without leaving incomplete states.
- **Spinner Integration**: Visual feedback during login operations.
- **YAML Configuration**: Manage registries through a simple YAML configuration file.

## Installation

### Go install to the win

```
go install github.com/user-cube/auth-refresher@latest
```

### Build it yourself

1. Clone the repository:
   ```bash
   git clone https://github.com/user-cube/auth-refresher.git
   cd auth-refresher
   ```

2. Build the project:
   ```bash
   go build -o auth-refresher
   ```

3. Run the binary:
   ```bash
   ./auth-refresher
   ```

## Usage

### Login to a Registry

Use the `login` command to log in to a registry:
```bash
./auth-refresher login
```

Follow the prompts to select a registry and log in.

### Logout from a Registry

Use the `logout` command to log out from a registry:
```bash
./auth-refresher logout
```

Follow the prompts to select a registry and log out. This command supports Docker, AWS ECR, and Helm registries.

### List Registries

Use the `list` command to view all configured registries:
```bash
./auth-refresher list
```

The output includes the registry name, type, URL, and timestamps for the last login and logout operations.

### Example Helm Login Command

For Helm registries, the tool executes the following commands internally:
```bash
AWS_REGION="us-west-2"
AWS_ACCOUNT_ID="123456789012"

# Get ECR credentials
export HELM_ECR_PASSWORD=$(aws ecr get-login-password --region $AWS_REGION)

# Log in to ECR via Helm registry
helm registry login \
    $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com \
    --username AWS \
    --password $HELM_ECR_PASSWORD
```

### Configuration

The tool uses a YAML configuration file located at `~/.auth-refresher/config.yaml`. Example:
```yaml
last_used_registry: my-docker-registry
registries:
  my-docker-registry:
    name: My Docker Registry
    type: docker
    url: https://index.docker.io/v1/
  my-aws-ecr:
    name: My AWS ECR
    type: aws
    url: 123456789012.dkr.ecr.us-west-2.amazonaws.com
    region: us-west-2
  my-helm-registry:
    name: My Helm Registry
    type: helm
    url: 123456789012.dkr.ecr.us-west-2.amazonaws.com
    region: us-west-2
```

## Development

### Prerequisites

- Go 1.18 or later

### Run Locally

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run main.go
   ```