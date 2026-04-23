// Package output provides formatters for rendering parse results as Markdown.
package output

// Mode controls the level of detail in the output.
type Mode int

const (
	// ModeSummary outputs function signatures, types, imports, and docstrings.
	ModeSummary Mode = iota
	// ModeFullCST outputs the complete concrete syntax tree.
	ModeFullCST
)

// FileInfo holds extracted metadata about a parsed file.
type FileInfo struct {
	File          string       `yaml:"file"`
	Language      string       `yaml:"language"`
	Package       string       `yaml:"package,omitempty"`
	Lines         int          `yaml:"lines"`
	FunctionCount int          `yaml:"function_count"`
	TypeCount     int          `yaml:"type_count"`
	ImportCount   int          `yaml:"import_count"`
	Tags          []string     `yaml:"tags,omitempty"`
	Functions     []FuncInfo   `yaml:"-"`
	Types         []TypeInfo   `yaml:"-"`
	Imports       []ImportInfo `yaml:"-"`
	Constants     []ConstInfo  `yaml:"-"`
}

// FuncInfo describes a function or method.
type FuncInfo struct {
	Name      string
	Signature string
	StartLine int
	EndLine   int
	Doc       string
	Receiver  string // Go: receiver type, empty for plain functions
	Exported  bool
}

// TypeInfo describes a type, struct, class, or interface.
type TypeInfo struct {
	Name      string
	Kind      string // "struct", "interface", "class", "type_alias"
	StartLine int
	EndLine   int
	Doc       string
	Fields    []FieldInfo
	Methods   []string
	Exported  bool
}

// FieldInfo describes a field in a struct/class.
type FieldInfo struct {
	Name string
	Type string
}

// ImportInfo describes an import statement.
type ImportInfo struct {
	Path  string
	Alias string
	Line  int
}

// ConstInfo describes a constant or variable declaration.
type ConstInfo struct {
	Name      string
	Value     string
	Type      string
	StartLine int
	Exported  bool
}

// ProjectInfo holds metadata for the index.md frontmatter.
type ProjectInfo struct {
	Project        string                  `yaml:"project"`
	Root           string                  `yaml:"root"`
	Generated      string                  `yaml:"generated"`
	SunsetVersion  string                  `yaml:"sunset_version"`
	Languages      map[string]LanguageStat `yaml:"languages"`
	TotalFiles     int                     `yaml:"total_files"`
	TotalFunctions int                     `yaml:"total_functions"`
	TotalTypes     int                     `yaml:"total_types"`
	Modules        []ModuleInfo            `yaml:"modules,omitempty"`
}

// LanguageStat holds file count and percentage for a language.
type LanguageStat struct {
	Files      int `yaml:"files"`
	Percentage int `yaml:"percentage"`
}

// ModuleInfo describes a directory/package grouping.
type ModuleInfo struct {
	Path  string `yaml:"path"`
	Files int    `yaml:"files"`
}
