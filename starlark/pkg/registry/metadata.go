package registry

// ModuleMetadata describes a complete module
type ModuleMetadata struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Version     string             `json:"version"`
	Category    string             `json:"category"` // "satellite", "ground", "data", "system"
	Author      string             `json:"author"`
	Functions   []FunctionMetadata `json:"functions"`
	Examples    []Example          `json:"examples"`
}

// FunctionMetadata describes a single function within a module
type FunctionMetadata struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  []ParameterMetadata `json:"parameters"`
	ReturnType  TypeMetadata        `json:"return_type"`
	Examples    []Example           `json:"examples"`
	Deprecated  bool                `json:"deprecated,omitempty"`
	Since       string              `json:"since,omitempty"` // Version when added
}

// ParameterMetadata describes a function parameter
type ParameterMetadata struct {
	Name        string       `json:"name"`
	Type        TypeMetadata `json:"type"`
	Description string       `json:"description"`
	Required    bool         `json:"required"`
	Default     string       `json:"default,omitempty"`
}

// TypeMetadata describes a type (for params/returns)
type TypeMetadata struct {
	Name       string  `json:"name"` // "string", "int", "float", "dict", "bool", "None"
	IsDict     bool    `json:"is_dict,omitempty"`
	DictSchema *Schema `json:"dict_schema,omitempty"` // For dict parameters
}

// Schema describes the structure of a dict parameter
type Schema struct {
	Fields []SchemaField `json:"fields"`
}

// SchemaField represents a field in a dictionary schema
type SchemaField struct {
	Name        string       `json:"name"`
	Type        TypeMetadata `json:"type"`
	Description string       `json:"description"`
	Required    bool         `json:"required"`
}

// Example shows usage of a function or module
type Example struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Code        string `json:"code"`
	Output      string `json:"output,omitempty"`
}
