
#!/usr/bin/env bash
set -e

# Install swag CLI tool
echo "Installing swag CLI..."
go install github.com/swaggo/swag/cmd/swag@latest

# Get required Go packages
echo "Getting Go dependencies..."
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/swag/cmd/swag

echo "Setup complete."

