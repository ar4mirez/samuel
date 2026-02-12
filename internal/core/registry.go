package core

// DefaultRegistry is the default GitHub repository for AICoF
const DefaultRegistry = "https://github.com/ar4mirez/aicof"

// DefaultOwner is the GitHub owner
const DefaultOwner = "ar4mirez"

// DefaultRepo is the GitHub repository name
const DefaultRepo = "aicof"

// TemplatePrefix is the path prefix in the repository where template files are stored
// This allows the CLI to find source files in the downloaded archive
const TemplatePrefix = "template/"

// Component represents an installable component
type Component struct {
	Name        string
	Path        string
	Description string
	Category    string   // Optional: "language", "framework", "skill", ""
	Tags        []string // Optional: additional search terms e.g. ["golang", "backend"]
}

// ComponentType represents the type of component
type ComponentType string

const (
	ComponentTypeLanguage  ComponentType = "language"
	ComponentTypeFramework ComponentType = "framework"
	ComponentTypeWorkflow  ComponentType = "workflow"
	ComponentTypeSkill     ComponentType = "skill"
)

// Languages contains all available language guide skills.
// Language guides are now Agent Skills stored at .agent/skills/<name>-guide/.
var Languages = []Component{
	{Name: "typescript", Path: ".agent/skills/typescript-guide", Description: "TypeScript/JavaScript", Category: "language", Tags: []string{"ts", "js", "javascript", "node"}},
	{Name: "python", Path: ".agent/skills/python-guide", Description: "Python", Category: "language", Tags: []string{"py", "pip", "django", "fastapi"}},
	{Name: "go", Path: ".agent/skills/go-guide", Description: "Go", Category: "language", Tags: []string{"golang", "goroutine"}},
	{Name: "rust", Path: ".agent/skills/rust-guide", Description: "Rust", Category: "language", Tags: []string{"cargo", "crate"}},
	{Name: "kotlin", Path: ".agent/skills/kotlin-guide", Description: "Kotlin", Category: "language", Tags: []string{"kt", "android", "jvm"}},
	{Name: "java", Path: ".agent/skills/java-guide", Description: "Java", Category: "language", Tags: []string{"jvm", "maven", "gradle"}},
	{Name: "csharp", Path: ".agent/skills/csharp-guide", Description: "C#/.NET", Category: "language", Tags: []string{"dotnet", "net", "cs"}},
	{Name: "php", Path: ".agent/skills/php-guide", Description: "PHP", Category: "language", Tags: []string{"composer", "laravel"}},
	{Name: "swift", Path: ".agent/skills/swift-guide", Description: "Swift", Category: "language", Tags: []string{"ios", "macos", "xcode"}},
	{Name: "cpp", Path: ".agent/skills/cpp-guide", Description: "C/C++", Category: "language", Tags: []string{"c", "cplusplus", "cmake"}},
	{Name: "ruby", Path: ".agent/skills/ruby-guide", Description: "Ruby", Category: "language", Tags: []string{"rb", "gem", "rails"}},
	{Name: "sql", Path: ".agent/skills/sql-guide", Description: "SQL", Category: "language", Tags: []string{"postgres", "mysql", "sqlite"}},
	{Name: "shell", Path: ".agent/skills/shell-guide", Description: "Shell/Bash", Category: "language", Tags: []string{"bash", "sh", "zsh", "scripting"}},
	{Name: "r", Path: ".agent/skills/r-guide", Description: "R", Category: "language", Tags: []string{"rstats", "tidyverse", "cran"}},
	{Name: "dart", Path: ".agent/skills/dart-guide", Description: "Dart", Category: "language", Tags: []string{"flutter", "pub"}},
	{Name: "html-css", Path: ".agent/skills/html-css-guide", Description: "HTML/CSS", Category: "language", Tags: []string{"html", "css", "web", "accessibility"}},
	{Name: "lua", Path: ".agent/skills/lua-guide", Description: "Lua", Category: "language", Tags: []string{"luajit", "love2d", "neovim"}},
	{Name: "assembly", Path: ".agent/skills/assembly-guide", Description: "Assembly", Category: "language", Tags: []string{"asm", "x86", "arm"}},
	{Name: "cuda", Path: ".agent/skills/cuda-guide", Description: "CUDA", Category: "language", Tags: []string{"gpu", "nvidia", "parallel"}},
	{Name: "solidity", Path: ".agent/skills/solidity-guide", Description: "Solidity", Category: "language", Tags: []string{"ethereum", "evm", "smart-contract"}},
	{Name: "zig", Path: ".agent/skills/zig-guide", Description: "Zig", Category: "language", Tags: []string{"systems", "comptime"}},
}

// Frameworks contains all available framework guide skills.
// Framework guides are now Agent Skills stored at .agent/skills/<name>/.
var Frameworks = []Component{
	// TypeScript/JavaScript
	{Name: "react", Path: ".agent/skills/react", Description: "React", Category: "framework", Tags: []string{"reactjs", "jsx", "tsx", "frontend"}},
	{Name: "nextjs", Path: ".agent/skills/nextjs", Description: "Next.js", Category: "framework", Tags: []string{"next", "ssr", "react", "vercel"}},
	{Name: "express", Path: ".agent/skills/express", Description: "Express.js", Category: "framework", Tags: []string{"node", "rest", "api", "middleware"}},
	// Python
	{Name: "django", Path: ".agent/skills/django", Description: "Django", Category: "framework", Tags: []string{"python", "orm", "drf", "admin"}},
	{Name: "fastapi", Path: ".agent/skills/fastapi", Description: "FastAPI", Category: "framework", Tags: []string{"python", "async", "pydantic", "openapi"}},
	{Name: "flask", Path: ".agent/skills/flask", Description: "Flask", Category: "framework", Tags: []string{"python", "wsgi", "jinja", "lightweight"}},
	// Go
	{Name: "gin", Path: ".agent/skills/gin", Description: "Gin", Category: "framework", Tags: []string{"go", "golang", "rest", "api"}},
	{Name: "echo", Path: ".agent/skills/echo", Description: "Echo", Category: "framework", Tags: []string{"go", "golang", "rest", "labstack"}},
	{Name: "fiber", Path: ".agent/skills/fiber", Description: "Fiber", Category: "framework", Tags: []string{"go", "golang", "express-like", "fasthttp"}},
	// Rust
	{Name: "axum", Path: ".agent/skills/axum", Description: "Axum", Category: "framework", Tags: []string{"rust", "tower", "async", "tokio"}},
	{Name: "actix-web", Path: ".agent/skills/actix-web", Description: "Actix-web", Category: "framework", Tags: []string{"rust", "actor", "async", "http"}},
	{Name: "rocket", Path: ".agent/skills/rocket", Description: "Rocket", Category: "framework", Tags: []string{"rust", "type-safe", "macros"}},
	// Kotlin
	{Name: "spring-boot-kotlin", Path: ".agent/skills/spring-boot-kotlin", Description: "Spring Boot (Kotlin)", Category: "framework", Tags: []string{"kotlin", "spring", "jvm", "coroutines"}},
	{Name: "ktor", Path: ".agent/skills/ktor", Description: "Ktor", Category: "framework", Tags: []string{"kotlin", "coroutines", "dsl", "jvm"}},
	{Name: "android-compose", Path: ".agent/skills/android-compose", Description: "Android Compose", Category: "framework", Tags: []string{"kotlin", "android", "material3", "jetpack"}},
	// Java
	{Name: "spring-boot-java", Path: ".agent/skills/spring-boot-java", Description: "Spring Boot (Java)", Category: "framework", Tags: []string{"java", "spring", "jvm", "jpa"}},
	{Name: "quarkus", Path: ".agent/skills/quarkus", Description: "Quarkus", Category: "framework", Tags: []string{"java", "graalvm", "reactive", "cloud-native"}},
	{Name: "micronaut", Path: ".agent/skills/micronaut", Description: "Micronaut", Category: "framework", Tags: []string{"java", "compile-time", "di", "cloud-native"}},
	// C#
	{Name: "aspnet-core", Path: ".agent/skills/aspnet-core", Description: "ASP.NET Core", Category: "framework", Tags: []string{"csharp", "dotnet", "minimal-api", "efcore"}},
	{Name: "blazor", Path: ".agent/skills/blazor", Description: "Blazor", Category: "framework", Tags: []string{"csharp", "dotnet", "wasm", "signalr"}},
	{Name: "unity", Path: ".agent/skills/unity", Description: "Unity", Category: "framework", Tags: []string{"csharp", "game", "3d", "scripting"}},
	// PHP
	{Name: "laravel", Path: ".agent/skills/laravel", Description: "Laravel", Category: "framework", Tags: []string{"php", "eloquent", "blade", "artisan"}},
	{Name: "symfony", Path: ".agent/skills/symfony", Description: "Symfony", Category: "framework", Tags: []string{"php", "doctrine", "twig", "components"}},
	{Name: "wordpress", Path: ".agent/skills/wordpress", Description: "WordPress", Category: "framework", Tags: []string{"php", "cms", "themes", "plugins"}},
	// Swift
	{Name: "swiftui", Path: ".agent/skills/swiftui", Description: "SwiftUI", Category: "framework", Tags: []string{"swift", "ios", "macos", "declarative"}},
	{Name: "uikit", Path: ".agent/skills/uikit", Description: "UIKit", Category: "framework", Tags: []string{"swift", "ios", "storyboard", "programmatic"}},
	{Name: "vapor", Path: ".agent/skills/vapor", Description: "Vapor", Category: "framework", Tags: []string{"swift", "server", "fluent", "async"}},
	// Ruby
	{Name: "rails", Path: ".agent/skills/rails", Description: "Rails", Category: "framework", Tags: []string{"ruby", "activerecord", "mvc", "hotwire"}},
	{Name: "sinatra", Path: ".agent/skills/sinatra", Description: "Sinatra", Category: "framework", Tags: []string{"ruby", "lightweight", "dsl", "rack"}},
	{Name: "hanami", Path: ".agent/skills/hanami", Description: "Hanami", Category: "framework", Tags: []string{"ruby", "clean-architecture", "dry-rb"}},
	// Dart
	{Name: "flutter", Path: ".agent/skills/flutter", Description: "Flutter", Category: "framework", Tags: []string{"dart", "mobile", "material", "riverpod"}},
	{Name: "shelf", Path: ".agent/skills/shelf", Description: "Shelf", Category: "framework", Tags: []string{"dart", "http", "server", "middleware"}},
	{Name: "dart-frog", Path: ".agent/skills/dart-frog", Description: "Dart Frog", Category: "framework", Tags: []string{"dart", "server", "file-based", "routing"}},
}

// Workflows contains all available workflows
var Workflows = []Component{
	{Name: "initialize-project", Path: ".agent/workflows/initialize-project.md", Description: "Project setup"},
	{Name: "create-rfd", Path: ".agent/workflows/create-rfd.md", Description: "Technical decision documents"},
	{Name: "create-prd", Path: ".agent/workflows/create-prd.md", Description: "Requirements documents"},
	{Name: "generate-tasks", Path: ".agent/workflows/generate-tasks.md", Description: "Task breakdown"},
	{Name: "code-review", Path: ".agent/workflows/code-review.md", Description: "Pre-commit quality review"},
	{Name: "security-audit", Path: ".agent/workflows/security-audit.md", Description: "Security assessment"},
	{Name: "testing-strategy", Path: ".agent/workflows/testing-strategy.md", Description: "Test planning"},
	{Name: "cleanup-project", Path: ".agent/workflows/cleanup-project.md", Description: "Prune unused guides"},
	{Name: "refactoring", Path: ".agent/workflows/refactoring.md", Description: "Technical debt remediation"},
	{Name: "dependency-update", Path: ".agent/workflows/dependency-update.md", Description: "Safe dependency updates"},
	{Name: "update-framework", Path: ".agent/workflows/update-framework.md", Description: "AICoF version updates"},
	{Name: "troubleshooting", Path: ".agent/workflows/troubleshooting.md", Description: "Debugging workflow"},
	{Name: "generate-agents-md", Path: ".agent/workflows/generate-agents-md.md", Description: "Cross-tool compatibility"},
	{Name: "document-work", Path: ".agent/workflows/document-work.md", Description: "Capture patterns"},
	{Name: "create-skill", Path: ".agent/workflows/create-skill.md", Description: "Create Agent Skills"},
}

// Skills contains bundled skills (community skills loaded dynamically).
// Language guide skills are also registered in the Languages slice for backward compatibility.
var Skills = []Component{
	{Name: "commit-message", Path: ".agent/skills/commit-message", Description: "Generate commit messages"},
	// Language guide skills (mirrored from Languages with skill naming)
	{Name: "go-guide", Path: ".agent/skills/go-guide", Description: "Go language guardrails and patterns", Category: "language", Tags: []string{"go", "golang"}},
	{Name: "typescript-guide", Path: ".agent/skills/typescript-guide", Description: "TypeScript/JavaScript guardrails and patterns", Category: "language", Tags: []string{"typescript", "ts", "javascript", "js"}},
	{Name: "python-guide", Path: ".agent/skills/python-guide", Description: "Python guardrails and patterns", Category: "language", Tags: []string{"python", "py"}},
	{Name: "rust-guide", Path: ".agent/skills/rust-guide", Description: "Rust guardrails and patterns", Category: "language", Tags: []string{"rust", "cargo"}},
	{Name: "kotlin-guide", Path: ".agent/skills/kotlin-guide", Description: "Kotlin guardrails and patterns", Category: "language", Tags: []string{"kotlin", "kt"}},
	{Name: "java-guide", Path: ".agent/skills/java-guide", Description: "Java guardrails and patterns", Category: "language", Tags: []string{"java", "jvm"}},
	{Name: "csharp-guide", Path: ".agent/skills/csharp-guide", Description: "C#/.NET guardrails and patterns", Category: "language", Tags: []string{"csharp", "dotnet"}},
	{Name: "php-guide", Path: ".agent/skills/php-guide", Description: "PHP guardrails and patterns", Category: "language", Tags: []string{"php", "composer"}},
	{Name: "swift-guide", Path: ".agent/skills/swift-guide", Description: "Swift guardrails and patterns", Category: "language", Tags: []string{"swift", "ios"}},
	{Name: "cpp-guide", Path: ".agent/skills/cpp-guide", Description: "C/C++ guardrails and patterns", Category: "language", Tags: []string{"cpp", "c", "cmake"}},
	{Name: "ruby-guide", Path: ".agent/skills/ruby-guide", Description: "Ruby guardrails and patterns", Category: "language", Tags: []string{"ruby", "rb"}},
	{Name: "sql-guide", Path: ".agent/skills/sql-guide", Description: "SQL guardrails and patterns", Category: "language", Tags: []string{"sql", "postgres", "mysql"}},
	{Name: "shell-guide", Path: ".agent/skills/shell-guide", Description: "Shell/Bash guardrails and patterns", Category: "language", Tags: []string{"shell", "bash", "sh"}},
	{Name: "r-guide", Path: ".agent/skills/r-guide", Description: "R guardrails and patterns", Category: "language", Tags: []string{"r", "rstats"}},
	{Name: "dart-guide", Path: ".agent/skills/dart-guide", Description: "Dart guardrails and patterns", Category: "language", Tags: []string{"dart", "flutter"}},
	{Name: "html-css-guide", Path: ".agent/skills/html-css-guide", Description: "HTML/CSS guardrails and patterns", Category: "language", Tags: []string{"html", "css"}},
	{Name: "lua-guide", Path: ".agent/skills/lua-guide", Description: "Lua guardrails and patterns", Category: "language", Tags: []string{"lua", "luajit"}},
	{Name: "assembly-guide", Path: ".agent/skills/assembly-guide", Description: "Assembly guardrails and patterns", Category: "language", Tags: []string{"assembly", "asm"}},
	{Name: "cuda-guide", Path: ".agent/skills/cuda-guide", Description: "CUDA guardrails and patterns", Category: "language", Tags: []string{"cuda", "gpu"}},
	{Name: "solidity-guide", Path: ".agent/skills/solidity-guide", Description: "Solidity guardrails and patterns", Category: "language", Tags: []string{"solidity", "ethereum"}},
	{Name: "zig-guide", Path: ".agent/skills/zig-guide", Description: "Zig guardrails and patterns", Category: "language", Tags: []string{"zig", "comptime"}},
	// Framework skills (mirrored from Frameworks)
	{Name: "react", Path: ".agent/skills/react", Description: "React framework guardrails and patterns", Category: "framework", Tags: []string{"react", "reactjs", "jsx", "tsx"}},
	{Name: "nextjs", Path: ".agent/skills/nextjs", Description: "Next.js framework guardrails and patterns", Category: "framework", Tags: []string{"nextjs", "next", "ssr"}},
	{Name: "express", Path: ".agent/skills/express", Description: "Express.js framework guardrails and patterns", Category: "framework", Tags: []string{"express", "node", "rest"}},
	{Name: "django", Path: ".agent/skills/django", Description: "Django framework guardrails and patterns", Category: "framework", Tags: []string{"django", "python", "orm"}},
	{Name: "fastapi", Path: ".agent/skills/fastapi", Description: "FastAPI framework guardrails and patterns", Category: "framework", Tags: []string{"fastapi", "python", "async"}},
	{Name: "flask", Path: ".agent/skills/flask", Description: "Flask framework guardrails and patterns", Category: "framework", Tags: []string{"flask", "python", "wsgi"}},
	{Name: "gin", Path: ".agent/skills/gin", Description: "Gin framework guardrails and patterns", Category: "framework", Tags: []string{"gin", "go", "golang"}},
	{Name: "echo", Path: ".agent/skills/echo", Description: "Echo framework guardrails and patterns", Category: "framework", Tags: []string{"echo", "go", "golang"}},
	{Name: "fiber", Path: ".agent/skills/fiber", Description: "Fiber framework guardrails and patterns", Category: "framework", Tags: []string{"fiber", "go", "golang"}},
	{Name: "axum", Path: ".agent/skills/axum", Description: "Axum framework guardrails and patterns", Category: "framework", Tags: []string{"axum", "rust", "tower"}},
	{Name: "actix-web", Path: ".agent/skills/actix-web", Description: "Actix-web framework guardrails and patterns", Category: "framework", Tags: []string{"actix-web", "rust", "actor"}},
	{Name: "rocket", Path: ".agent/skills/rocket", Description: "Rocket framework guardrails and patterns", Category: "framework", Tags: []string{"rocket", "rust", "macros"}},
	{Name: "spring-boot-kotlin", Path: ".agent/skills/spring-boot-kotlin", Description: "Spring Boot (Kotlin) framework guardrails and patterns", Category: "framework", Tags: []string{"spring", "kotlin", "jvm"}},
	{Name: "ktor", Path: ".agent/skills/ktor", Description: "Ktor framework guardrails and patterns", Category: "framework", Tags: []string{"ktor", "kotlin", "coroutines"}},
	{Name: "android-compose", Path: ".agent/skills/android-compose", Description: "Android Compose framework guardrails and patterns", Category: "framework", Tags: []string{"compose", "android", "kotlin"}},
	{Name: "spring-boot-java", Path: ".agent/skills/spring-boot-java", Description: "Spring Boot (Java) framework guardrails and patterns", Category: "framework", Tags: []string{"spring", "java", "jvm"}},
	{Name: "quarkus", Path: ".agent/skills/quarkus", Description: "Quarkus framework guardrails and patterns", Category: "framework", Tags: []string{"quarkus", "java", "graalvm"}},
	{Name: "micronaut", Path: ".agent/skills/micronaut", Description: "Micronaut framework guardrails and patterns", Category: "framework", Tags: []string{"micronaut", "java", "di"}},
	{Name: "aspnet-core", Path: ".agent/skills/aspnet-core", Description: "ASP.NET Core framework guardrails and patterns", Category: "framework", Tags: []string{"aspnet", "csharp", "dotnet"}},
	{Name: "blazor", Path: ".agent/skills/blazor", Description: "Blazor framework guardrails and patterns", Category: "framework", Tags: []string{"blazor", "csharp", "wasm"}},
	{Name: "unity", Path: ".agent/skills/unity", Description: "Unity framework guardrails and patterns", Category: "framework", Tags: []string{"unity", "csharp", "game"}},
	{Name: "laravel", Path: ".agent/skills/laravel", Description: "Laravel framework guardrails and patterns", Category: "framework", Tags: []string{"laravel", "php", "eloquent"}},
	{Name: "symfony", Path: ".agent/skills/symfony", Description: "Symfony framework guardrails and patterns", Category: "framework", Tags: []string{"symfony", "php", "doctrine"}},
	{Name: "wordpress", Path: ".agent/skills/wordpress", Description: "WordPress framework guardrails and patterns", Category: "framework", Tags: []string{"wordpress", "php", "cms"}},
	{Name: "swiftui", Path: ".agent/skills/swiftui", Description: "SwiftUI framework guardrails and patterns", Category: "framework", Tags: []string{"swiftui", "swift", "ios"}},
	{Name: "uikit", Path: ".agent/skills/uikit", Description: "UIKit framework guardrails and patterns", Category: "framework", Tags: []string{"uikit", "swift", "ios"}},
	{Name: "vapor", Path: ".agent/skills/vapor", Description: "Vapor framework guardrails and patterns", Category: "framework", Tags: []string{"vapor", "swift", "server"}},
	{Name: "rails", Path: ".agent/skills/rails", Description: "Rails framework guardrails and patterns", Category: "framework", Tags: []string{"rails", "ruby", "activerecord"}},
	{Name: "sinatra", Path: ".agent/skills/sinatra", Description: "Sinatra framework guardrails and patterns", Category: "framework", Tags: []string{"sinatra", "ruby", "dsl"}},
	{Name: "hanami", Path: ".agent/skills/hanami", Description: "Hanami framework guardrails and patterns", Category: "framework", Tags: []string{"hanami", "ruby", "dry-rb"}},
	{Name: "flutter", Path: ".agent/skills/flutter", Description: "Flutter framework guardrails and patterns", Category: "framework", Tags: []string{"flutter", "dart", "mobile"}},
	{Name: "shelf", Path: ".agent/skills/shelf", Description: "Shelf framework guardrails and patterns", Category: "framework", Tags: []string{"shelf", "dart", "http"}},
	{Name: "dart-frog", Path: ".agent/skills/dart-frog", Description: "Dart Frog framework guardrails and patterns", Category: "framework", Tags: []string{"dart-frog", "dart", "server"}},
}

// CoreFiles contains essential files always installed
var CoreFiles = []string{
	"CLAUDE.md",
	"AI_INSTRUCTIONS.md",
	".agent/README.md",
	".agent/project.md.template",
	".agent/state.md.template",
	".agent/skills/README.md",
	".agent/workflows/README.md",
	".agent/memory/.gitkeep",
	".agent/tasks/.gitkeep",
	".agent/rfd/.gitkeep",
}

// Template represents a predefined set of components
type Template struct {
	Name        string
	Description string
	Languages   []string
	Frameworks  []string
	Workflows   []string
}

// Templates contains predefined installation templates
var Templates = []Template{
	{
		Name:        "full",
		Description: "All guides and workflows (~160 files)",
		Languages:   getAllNames(Languages),
		Frameworks:  getAllNames(Frameworks),
		Workflows:   []string{"all"},
	},
	{
		Name:        "starter",
		Description: "Core files + common languages (TypeScript, Python, Go)",
		Languages:   []string{"typescript", "python", "go"},
		Frameworks:  []string{},
		Workflows:   []string{"all"},
	},
	{
		Name:        "minimal",
		Description: "Just CLAUDE.md and workflows",
		Languages:   []string{},
		Frameworks:  []string{},
		Workflows:   []string{"all"},
	},
}

// FindLanguage finds a language by name
func FindLanguage(name string) *Component {
	for _, lang := range Languages {
		if lang.Name == name {
			return &lang
		}
	}
	return nil
}

// FindFramework finds a framework by name
func FindFramework(name string) *Component {
	for _, fw := range Frameworks {
		if fw.Name == name {
			return &fw
		}
	}
	return nil
}

// FindWorkflow finds a workflow by name
func FindWorkflow(name string) *Component {
	for _, wf := range Workflows {
		if wf.Name == name {
			return &wf
		}
	}
	return nil
}

// FindSkill finds a skill by name
func FindSkill(name string) *Component {
	for _, s := range Skills {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

// FindTemplate finds a template by name
func FindTemplate(name string) *Template {
	for _, t := range Templates {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// GetComponentPaths returns all paths for a given set of components
func GetComponentPaths(languages, frameworks, workflows []string) []string {
	var paths []string

	// Add core files
	paths = append(paths, CoreFiles...)

	// Add languages
	for _, name := range languages {
		if lang := FindLanguage(name); lang != nil {
			paths = append(paths, lang.Path)
		}
	}

	// Add frameworks
	for _, name := range frameworks {
		if fw := FindFramework(name); fw != nil {
			paths = append(paths, fw.Path)
		}
	}

	// Add workflows
	if len(workflows) == 1 && workflows[0] == "all" {
		for _, wf := range Workflows {
			paths = append(paths, wf.Path)
		}
	} else {
		for _, name := range workflows {
			if wf := FindWorkflow(name); wf != nil {
				paths = append(paths, wf.Path)
			}
		}
	}

	return paths
}

func getAllNames(components []Component) []string {
	names := make([]string, len(components))
	for i, c := range components {
		names[i] = c.Name
	}
	return names
}

// GetAllLanguageNames returns all language names
func GetAllLanguageNames() []string {
	return getAllNames(Languages)
}

// GetAllFrameworkNames returns all framework names
func GetAllFrameworkNames() []string {
	return getAllNames(Frameworks)
}

// GetAllWorkflowNames returns all workflow names
func GetAllWorkflowNames() []string {
	return getAllNames(Workflows)
}

// GetAllSkillNames returns all skill names
func GetAllSkillNames() []string {
	return getAllNames(Skills)
}

// GetLanguageSkills returns skills with category "language"
func GetLanguageSkills() []Component {
	var result []Component
	for _, s := range Skills {
		if s.Category == "language" {
			result = append(result, s)
		}
	}
	return result
}

// LanguageToSkillName converts a language name to its skill name (e.g., "go" -> "go-guide")
func LanguageToSkillName(langName string) string {
	return langName + "-guide"
}

// SkillToLanguageName converts a skill name to its language name (e.g., "go-guide" -> "go")
func SkillToLanguageName(skillName string) string {
	if len(skillName) > 6 && skillName[len(skillName)-6:] == "-guide" {
		return skillName[:len(skillName)-6]
	}
	return skillName
}

// FrameworkToSkillName converts a framework name to its skill name.
// For frameworks, the skill name IS the framework name (no suffix).
func FrameworkToSkillName(fwName string) string {
	return fwName
}

// GetFrameworkSkills returns skills with category "framework"
func GetFrameworkSkills() []Component {
	var result []Component
	for _, s := range Skills {
		if s.Category == "framework" {
			result = append(result, s)
		}
	}
	return result
}

// GetSourcePath returns the source path in the repository for a given destination path
// This prepends the TemplatePrefix to locate files in the downloaded archive
func GetSourcePath(destPath string) string {
	return TemplatePrefix + destPath
}

// GetSourcePaths returns source paths for a slice of destination paths
func GetSourcePaths(destPaths []string) []string {
	srcPaths := make([]string, len(destPaths))
	for i, p := range destPaths {
		srcPaths[i] = GetSourcePath(p)
	}
	return srcPaths
}
