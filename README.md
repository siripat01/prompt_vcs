# 🚀 Prompt Version Control System

A developer portfolio project for managing LLM (Large Language Model) prompts with Git-like version control. Built with Go (CLI + API) and React/TypeScript (Web UI).

## 🎯 Features

### ✅ Implemented Features
| Feature | CLI | Web |
|---------|-----|-----|
| Create prompts | ✅ | ✅ |
| List prompts | ✅ | ✅ |
| Version history | 🔄 | 🔄 |
| Git-like commits | 🔄 | 🔄 |
| Prompt collections | ✅ | 🔄 |
| Tagging system | 🔄 | 🔄 |
| Diff viewer | 🔄 | 🔄 |
| Search & filter | 🔄 | 🔄 |

- **Hydrate & Rehydrate**: manipulate the database states based on the files states if needed.

---

## 🏗️ Architecture

```
prompt-vcs/
├── cmd/prompt-vcs/         # CLI entry point
│   └── main.go
├── internal/
│   ├── api/                 # HTTP API server (Echo framework)
│   ├── cli/                 # CLI commands (Cobra)
│   ├── models/              # Data structures
│   └── storage/db/          # SQLite database layer
├── web/                     # React + TypeScript frontend
│   ├── src/components/      # UI components
│   └── src/App.tsx          # Main application
└── go.mod                   # Go module definition
```

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.22+**
- **Node.js 18+** (for web UI)
- **npm or yarn**

### Backend Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd prompt-vcs

# Install dependencies
go mod tidy

# Initialize database and run CLI
go run cmd/prompt-vcs/main.go init
go run cmd/prompt-vcs/main.go add hello-chat "Hello {name}, how are you?"
go run cmd/prompt-vcs/main.go list

# Start API server
go run cmd/prompt-vcs/main.go serve
```

### Frontend Setup
```bash
cd web
npm install
npm run dev
```

---

## 📚 Usage Examples

### CLI
```bash
# Initialize repository
prompt-vcs init

# Add a new prompt
prompt-vcs add "Code Review" "Review this code: {language}\n\n{code}"

# List all prompts
prompt-vcs list

# Show prompt details
prompt-vcs show <prompt-id>

# Commit a new version
prompt-vcs commit <prompt-id> "Added error handling instructions"

# See differences
prompt-vcs diff <prompt-id> <version1> <version2>

# Rollback to previous version
prompt-vcs rollback <prompt-id> <version>
```

### Web UI
Navigate to `http://localhost:5173` to:
- ✨ Create prompts via forms
- 📋 View all prompts in a grid
- 🏷️ Filter by tags and collections
- 📝 Edit and version prompts
- 📊 View version history

---

## 🛠️ Technology Stack

| Layer | Technology | Reason |
|-------|-----------|--------|
| **Backend** | Go + Echo | Fast, simple, concurrent-safe |
| **CLI** | Cobra | Industry standard, robust |
| **Storage** | SQLite + files | Lightweight, versionable |
| **Frontend** | React 18 + TypeScript | Type-safe, modern |
| **Styling** | Tailwind CSS (planned) | Utility-first, rapid development |
| **Build** | Vite | Fast dev server, optimized builds |

---

## 🤝 Contributing

Contributions welcome! Areas for improvement:
- [ ] Full CRUD operations in web UI
- [ ] Real-time collaboration
- [ ] Git integration (optional)
- [ ] Cloud sync (S3, GitHub)
- [ ] Plugin system for custom prompts

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

**Built with 💙 for the developer community**