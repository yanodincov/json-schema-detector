# JSON Schema Detector

A tool for automatic JSON document analysis and generation of structured schemas with JSON Schema standard support.

## Features

- 🔍 **Automatic Type Analysis** - Detection of primitive and composite data types
- 📋 **JSON Schema Generation** - Creation of standard JSON Schema documents
- 🔄 **Schema Updates** - Merging new data with existing schemas
- ✅ **Validation** - Checking JSON data against schemas
- 📊 **Statistics** - Detailed analytics on data structures
- 🎯 **Enum Types** - Interactive field conversion to enum with value selection
- 🔗 **Polymorphic Types** - Creation of polymorphic objects with oneOf/anyOf
- 🛠️ **Interactive Field Management** - Changing types and descriptions via commands
- 📍 **JSON Path Navigation** - Precise field addressing in complex schemas
- 🔧 **Smart Default Values** - Automatic filling and updating of default values
- 🛡️ **Overwrite Protection** - Mechanism to preserve critical default values
- 📦 **Single Object Support** - JSON analysis without mandatory data field
- 🔄 **Automatic Schema Commits** - Automatic saving of changes to git

## Installation

```bash
go install github.com/yanodincov/json-schema-detector/cmd@latest
```

## Usage

### JSON File Analysis

```bash
# Basic analysis
json-schema-detector analyze examples/sample_data.json

# Analysis with output file specification
json-schema-detector analyze examples/sample_data.json -o user_schema.json

# Analysis with automatic commit of changes
json-schema-detector analyze examples/sample_data.json --auto-commit
```

### Schema Updates

```bash
# Update existing schema with new data
json-schema-detector update user_schema.json -i new_data.json

# Update with automatic commit
json-schema-detector update user_schema.json -i new_data.json --auto-commit
```

### Data Validation

```bash
# Basic validation
json-schema-detector validate data.json user_schema.json

# Verbose validation
json-schema-detector validate data.json user_schema.json -v

# Strict validation
json-schema-detector validate data.json user_schema.json -s
```

### Interactive Field Management

```bash
# View all fields in schema
json-schema-detector list-fields user_schema.json

# View fields with types
json-schema-detector list-fields user_schema.json --types

# Detailed field view
json-schema-detector list-fields user_schema.json --verbose

# Convert field to enum type
json-schema-detector update-field user_schema.json "data.0.role" enum

# Create polymorphic type
json-schema-detector update-field user_schema.json "data.0.user" polymorph

# Update field description
json-schema-detector update-field user_schema.json "data.0.id" description

# Protect default value from overwriting
json-schema-detector update-field user_schema.json "data.0.role" preserve-default

# Interactive mode (operation selection)
json-schema-detector update-field user_schema.json "data.0.status"

# Field update with automatic commit
json-schema-detector update-field user_schema.json "data.0.role" enum --auto-commit
```

### JSON Path Navigation

For working with fields in complex schemas, JSON Path syntax is used:

```bash
# Simple fields
data.name           # name field in data object
data.id             # id field in data object

# Arrays
data.0.name         # name field in first element of data array
users.0.profile.age # age field in first user's profile

# Nested objects
user.profile.settings.theme    # deeply nested field
config.database.connection.host # field in configuration

# Command examples
json-schema-detector list-fields schema.json
json-schema-detector update-field schema.json "data.0.role" enum
json-schema-detector update-field schema.json "users.0.profile.type" polymorph
```

### Smart Default Values

The analyzer automatically fills default values with smart logic:

```bash
# On first analysis, default is filled with current value
json-schema-detector analyze user.json

# On schema update, default is reset if value changed
json-schema-detector update user.schema.json -i user_updated.json

# Protect critical default values from overwriting
json-schema-detector update-field user.schema.json "role" preserve-default
```

**Default filling rules:**
- ✅ Filled on first analysis (if value is not empty)
- ✅ Reset on update if value changed
- ✅ Not filled for empty values (`""`, `0`)
- ✅ Always filled for boolean values
- ✅ Protected from overwriting with `x-preserve-default` flag

### Single Object Support

The analyzer automatically determines data structure:

```bash
# Structure with data array
{
  "data": [
    {"id": 1, "name": "John"}
  ]
}

# Single object (processed as one element)
{
  "id": 1,
  "name": "John",
  "profile": {
    "age": 30
  }
}
```

### Automatic Schema Commits

All commands support automatic commit of changes to git:

```bash
# Analysis with commit
json-schema-detector analyze data.json --auto-commit
# Creates commit: "schema: analyze data.schema.json"

# Update with commit  
json-schema-detector update schema.json -i new_data.json --auto-commit
# Creates commit: "schema: update schema.json"

# Field change with commit
json-schema-detector update-field schema.json "field" enum --auto-commit
# Creates commit: "schema: update-field schema.json"
```

**Requirements:**
- Git must be installed and available in PATH
- Working directory must be a git repository
- Schema file will be automatically added to staging area

**Commit message format:**
```
schema: <operation> <schema_file_name>
```

## Configuration

The tool works without configuration files and uses sensible defaults. 

Main behavior parameters:
- JSON Schema draft-07 format
- Automatic data type detection
- Smart default values for non-empty fields
- Support for enum and polymorphic types via interactive commands

## Usage Examples

### Interactive Field Management

```bash
# View all fields in schema
$ json-schema-detector list-fields examples/sample_data.schema.json

🔍 Fields in schema: examples/sample_data.schema.json
├── data (array)
│   ├── 0.active (boolean)
│   ├── 0.created_at (string)
│   ├── 0.id (number)
│   ├── 0.name (string)
│   ├── 0.role (string - enum: admin, user, manager)
│   └── 0.permissions (array)

# Convert field to enum type
$ json-schema-detector update-field examples/sample_data.schema.json "data.0.role" enum

🔧 Updating field in schema
📄 Schema file: examples/sample_data.schema.json
🎯 Field path: data.0.role
🔄 Operation: enum

📝 Enter possible values for enum (one per line):
💡 Finish input with empty line

Value: admin
Value: user
Value: manager
Value: 

📝 Field description (optional): User role in system
✅ Field converted to enum with 3 values
🎯 Values: [admin user manager]
✅ Field successfully updated: data.0.role
```

### Working with Default Values

```bash
# Analyze single object with automatic default filling
$ json-schema-detector analyze examples/user_simple.json

# Result includes default values
{
  "role": {
    "type": "string",
    "default": "admin"
  },
  "active": {
    "type": "boolean", 
    "default": true
  }
}

# Protect critical default value
$ json-schema-detector update-field examples/user_simple.schema.json "role" preserve-default

🔧 Updating field in schema
📄 Schema file: examples/user_simple.schema.json
🎯 Field path: role
🔄 Operation: preserve-default

🔒 Protecting default value from overwriting
✅ Default value protected: admin
✅ Field protected from default overwriting: role

# Now field contains protection flag
{
  "role": {
    "type": "string",
    "default": "admin",
    "x-preserve-default": true
  }
}
```

### Input Data (sample_data.json)

```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "role": "admin",
      "permissions": ["read", "write", "delete"],
      "active": true
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "role": "user",
      "active": true
    }
  ]
}
```

### Generated Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "data": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {
            "type": "number",
            "description": "Unique identifier"
          },
          "name": {
            "type": "string",
            "description": "User name"
          },
          "role": {
            "type": "string",
            "enum": ["admin", "user"],
            "description": "Role in system"
          },
          "permissions": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Access permissions"
          },
          "active": {
            "type": "boolean",
            "description": "Activity status"
          }
        },
        "required": ["id", "name", "role", "active"]
      }
    }
  },
  "required": ["data"]
}
```

## Building from Source

```bash
git clone https://github.com/yanodincov/json-schema-detector.git
cd json-schema-detector
go mod tidy
go build -o json-schema-detector cmd/main.go
```

## Development

### Project Structure

```
├── cmd/                    # CLI commands
│   ├── main.go            # Entry point
│   ├── root/              # Root command
│   ├── analyze/           # Analyze command
│   ├── update/            # Update command
│   ├── validate/          # Validate command
│   ├── update-field/      # Interactive field management
│   └── list-fields/       # Schema field viewer
├── pkg/                   # Core packages
│   ├── types/             # Data types
│   ├── analyzer/          # JSON analyzer
│   ├── validator/         # Schema validator
│   └── fieldmanager/      # Schema field manager
├── examples/              # Example data
└── schemas/               # Generated schemas
```

## Development Roadmap

### In Progress
- 🔄 **Polymorphic Types** - Creating oneOf/anyOf schemas for different object variants
- 🧪 **Extended Testing** - Automated tests for all components
- 📈 **Usage Statistics** - Analytics on fields and types

### Planned
- 🌐 **Web Interface** - Graphical interface for schema management
- 🔌 **API Interface** - REST API for integration with other systems
- 📊 **Extended Analytics** - Detailed reports on data structures
- 🎨 **Schema Customization** - Themes and output settings
- 🔍 **Schema Search** - Quick search for fields and types

### Running Tests

```bash
go test ./...
```

### Testing with Examples

```bash
# Analyze sample data
./json-schema-detector analyze examples/sample_data.json

# Validate against generated schema
./json-schema-detector validate examples/valid_data.json examples/sample_data.schema.json

# Interactive field updates
./json-schema-detector update-field examples/sample_data.schema.json "data.0.role" enum
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details. 