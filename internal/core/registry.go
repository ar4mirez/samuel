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
}

// ComponentType represents the type of component
type ComponentType string

const (
	ComponentTypeLanguage  ComponentType = "language"
	ComponentTypeFramework ComponentType = "framework"
	ComponentTypeWorkflow  ComponentType = "workflow"
)

// Languages contains all available language guides
var Languages = []Component{
	{Name: "typescript", Path: ".agent/language-guides/typescript.md", Description: "TypeScript/JavaScript"},
	{Name: "python", Path: ".agent/language-guides/python.md", Description: "Python"},
	{Name: "go", Path: ".agent/language-guides/go.md", Description: "Go"},
	{Name: "rust", Path: ".agent/language-guides/rust.md", Description: "Rust"},
	{Name: "kotlin", Path: ".agent/language-guides/kotlin.md", Description: "Kotlin"},
	{Name: "java", Path: ".agent/language-guides/java.md", Description: "Java"},
	{Name: "csharp", Path: ".agent/language-guides/csharp.md", Description: "C#/.NET"},
	{Name: "php", Path: ".agent/language-guides/php.md", Description: "PHP"},
	{Name: "swift", Path: ".agent/language-guides/swift.md", Description: "Swift"},
	{Name: "cpp", Path: ".agent/language-guides/cpp.md", Description: "C/C++"},
	{Name: "ruby", Path: ".agent/language-guides/ruby.md", Description: "Ruby"},
	{Name: "sql", Path: ".agent/language-guides/sql.md", Description: "SQL"},
	{Name: "shell", Path: ".agent/language-guides/shell.md", Description: "Shell/Bash"},
	{Name: "r", Path: ".agent/language-guides/r.md", Description: "R"},
	{Name: "dart", Path: ".agent/language-guides/dart.md", Description: "Dart"},
	{Name: "html-css", Path: ".agent/language-guides/html-css.md", Description: "HTML/CSS"},
	{Name: "lua", Path: ".agent/language-guides/lua.md", Description: "Lua"},
	{Name: "assembly", Path: ".agent/language-guides/assembly.md", Description: "Assembly"},
	{Name: "cuda", Path: ".agent/language-guides/cuda.md", Description: "CUDA"},
	{Name: "solidity", Path: ".agent/language-guides/solidity.md", Description: "Solidity"},
	{Name: "zig", Path: ".agent/language-guides/zig.md", Description: "Zig"},
}

// Frameworks contains all available framework guides
var Frameworks = []Component{
	// TypeScript/JavaScript
	{Name: "react", Path: ".agent/framework-guides/react.md", Description: "React"},
	{Name: "nextjs", Path: ".agent/framework-guides/nextjs.md", Description: "Next.js"},
	{Name: "express", Path: ".agent/framework-guides/express.md", Description: "Express.js"},
	// Python
	{Name: "django", Path: ".agent/framework-guides/django.md", Description: "Django"},
	{Name: "fastapi", Path: ".agent/framework-guides/fastapi.md", Description: "FastAPI"},
	{Name: "flask", Path: ".agent/framework-guides/flask.md", Description: "Flask"},
	// Go
	{Name: "gin", Path: ".agent/framework-guides/gin.md", Description: "Gin"},
	{Name: "echo", Path: ".agent/framework-guides/echo.md", Description: "Echo"},
	{Name: "fiber", Path: ".agent/framework-guides/fiber.md", Description: "Fiber"},
	// Rust
	{Name: "axum", Path: ".agent/framework-guides/axum.md", Description: "Axum"},
	{Name: "actix-web", Path: ".agent/framework-guides/actix-web.md", Description: "Actix-web"},
	{Name: "rocket", Path: ".agent/framework-guides/rocket.md", Description: "Rocket"},
	// Kotlin
	{Name: "spring-boot-kotlin", Path: ".agent/framework-guides/spring-boot-kotlin.md", Description: "Spring Boot (Kotlin)"},
	{Name: "ktor", Path: ".agent/framework-guides/ktor.md", Description: "Ktor"},
	{Name: "android-compose", Path: ".agent/framework-guides/android-compose.md", Description: "Android Compose"},
	// Java
	{Name: "spring-boot-java", Path: ".agent/framework-guides/spring-boot-java.md", Description: "Spring Boot (Java)"},
	{Name: "quarkus", Path: ".agent/framework-guides/quarkus.md", Description: "Quarkus"},
	{Name: "micronaut", Path: ".agent/framework-guides/micronaut.md", Description: "Micronaut"},
	// C#
	{Name: "aspnet-core", Path: ".agent/framework-guides/aspnet-core.md", Description: "ASP.NET Core"},
	{Name: "blazor", Path: ".agent/framework-guides/blazor.md", Description: "Blazor"},
	{Name: "unity", Path: ".agent/framework-guides/unity.md", Description: "Unity"},
	// PHP
	{Name: "laravel", Path: ".agent/framework-guides/laravel.md", Description: "Laravel"},
	{Name: "symfony", Path: ".agent/framework-guides/symfony.md", Description: "Symfony"},
	{Name: "wordpress", Path: ".agent/framework-guides/wordpress.md", Description: "WordPress"},
	// Swift
	{Name: "swiftui", Path: ".agent/framework-guides/swiftui.md", Description: "SwiftUI"},
	{Name: "uikit", Path: ".agent/framework-guides/uikit.md", Description: "UIKit"},
	{Name: "vapor", Path: ".agent/framework-guides/vapor.md", Description: "Vapor"},
	// Ruby
	{Name: "rails", Path: ".agent/framework-guides/rails.md", Description: "Rails"},
	{Name: "sinatra", Path: ".agent/framework-guides/sinatra.md", Description: "Sinatra"},
	{Name: "hanami", Path: ".agent/framework-guides/hanami.md", Description: "Hanami"},
	// Dart
	{Name: "flutter", Path: ".agent/framework-guides/flutter.md", Description: "Flutter"},
	{Name: "shelf", Path: ".agent/framework-guides/shelf.md", Description: "Shelf"},
	{Name: "dart-frog", Path: ".agent/framework-guides/dart-frog.md", Description: "Dart Frog"},
}

// Workflows contains all available workflows
var Workflows = []Component{
	{Name: "initialize-project", Path: ".agent/workflows/initialize-project.md", Description: "Project setup"},
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
}

// CoreFiles contains essential files always installed
var CoreFiles = []string{
	"CLAUDE.md",
	"AI_INSTRUCTIONS.md",
	".agent/README.md",
	".agent/project.md.template",
	".agent/state.md.template",
	".agent/language-guides/README.md",
	".agent/framework-guides/README.md",
	".agent/workflows/README.md",
	".agent/memory/.gitkeep",
	".agent/tasks/.gitkeep",
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
