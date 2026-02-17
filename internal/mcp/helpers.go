package mcp

import (
	"encoding/json"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// parseArgs unmarshals the raw JSON arguments from an MCP request into a map
func parseArgs(raw json.RawMessage) map[string]interface{} {
	var args map[string]interface{}
	if len(raw) > 0 {
		json.Unmarshal(raw, &args)
	}
	if args == nil {
		args = make(map[string]interface{})
	}
	return args
}

// toolResult creates a successful MCP tool result with text content
func toolResult(text string) *gomcp.CallToolResult {
	return &gomcp.CallToolResult{
		Content: []gomcp.Content{&gomcp.TextContent{Text: text}},
	}
}

// toolError creates an MCP tool result indicating an error
func toolError(message string) *gomcp.CallToolResult {
	return &gomcp.CallToolResult{
		IsError: true,
		Content: []gomcp.Content{&gomcp.TextContent{Text: message}},
	}
}

// inputSchema creates a JSON schema for tool input as json.RawMessage
func inputSchema(properties map[string]interface{}, required []string) json.RawMessage {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	data, _ := json.Marshal(schema)
	return data
}

// prop creates a JSON schema property definition
func prop(typ, description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        typ,
		"description": description,
	}
}

// arrayProp creates a JSON schema array property definition
func arrayProp(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "array",
		"description": description,
		"items":       map[string]interface{}{"type": "string"},
	}
}

// getStringArg extracts a string argument from a parsed args map
func getStringArg(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}

// getIntArg extracts an integer argument from a parsed args map (JSON numbers are float64)
func getIntArg(args map[string]interface{}, key string, defaultVal int) int {
	if v, ok := args[key].(float64); ok {
		return int(v)
	}
	return defaultVal
}

// getStringSliceArg extracts a string slice argument from a parsed args map
func getStringSliceArg(args map[string]interface{}, key string) []string {
	if arr, ok := args[key].([]interface{}); ok {
		var result []string
		for _, v := range arr {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}
