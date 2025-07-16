#!/bin/bash

# Test script for Hue MCP Server
# This script tests basic MCP server functionality

echo "Testing Hue MCP Server..."

# Set test environment variables
export HUE_BRIDGE_IP="192.168.87.51"
export HUE_USERNAME="test-username"

# Build the server
echo "Building server..."
go build -o hue-mcp . || exit 1

echo "Build successful!"

# Test basic MCP handshake
echo "Testing MCP handshake..."
echo '{"jsonrpc":"2.0","method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0.0"},"capabilities":{}},"id":1}' | ./hue-mcp

echo ""
echo "To test with your actual bridge:"
echo "1. Set HUE_USERNAME to your actual Hue bridge username"
echo "2. Run: ./hue-mcp"
echo "3. Send MCP commands via stdin"
echo ""
echo "For Claude Desktop integration, add to your config:"
echo '  ~/Library/Application Support/Claude/claude_desktop_config.json'