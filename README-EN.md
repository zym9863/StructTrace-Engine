[English](README-EN.md) | [简体中文](README.md)

# 🌳 StructTrace Engine

**Data Structure Evolution Tracking Engine** - Real-time visualization of data structure operations for deep understanding of algorithm principles

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?logo=typescript)

---

## ✨ Features

- 🎬 **Frame-by-Frame Animation** - Observe every step of algorithm execution
- 🌲 **Multiple Data Structures** - Red-Black Tree, AVL Tree, Graph, HashMap
- ⚡ **Concurrent Benchmarking** - Real-time performance comparison of different data structures
- 📊 **Interactive Charts** - Intuitive display of execution time, memory usage, and operation speed
- 🎨 **Modern UI** - Beautiful dark theme with smooth animations

## 🖼️ Preview

### Algorithm Visualization
Watch Red-Black Tree insertions, rotations, and recoloring operations in real-time

### Performance Benchmarking
Concurrent testing of insertion/search performance across multiple data structures

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+
- **Node.js** 18+
- **pnpm** 8+

### Installation & Running

```bash
# Clone the repository
git clone https://github.com/zym9863/StructTrace-Engine.git
cd StructTrace-Engine

# Start the backend
cd backend
go run main.go

# In a new terminal, start the frontend
cd frontend
pnpm install
pnpm dev
```

Visit http://localhost:5173 to get started

---

## 📁 Project Structure

```
StructTrace-Engine/
├── backend/                    # Go Backend
│   ├── main.go                # Entry + Gin Routes
│   ├── handlers/              # HTTP Handlers
│   │   ├── operation.go       # Algorithm Operation API
│   │   └── benchmark.go       # Benchmark API + SSE
│   ├── datastructures/        # Data Structure Implementations
│   │   ├── rbtree.go         # Red-Black Tree
│   │   ├── avltree.go        # AVL Tree
│   │   ├── graph.go          # Graph + Dijkstra
│   │   └── snapshot.go       # Snapshot Structure Definitions
│   └── benchmark/             # Benchmark Service
│       └── runner.go         # Goroutine Concurrent Testing
│
├── frontend/                   # React Frontend
│   └── src/
│       ├── components/        # UI Components
│       │   ├── TreeCanvas     # Tree Visualization
│       │   ├── AnimationPlayer# Animation Player
│       │   ├── OperationPanel # Operation Panel
│       │   └── BenchmarkChart # Benchmark Charts
│       ├── pages/             # Pages
│       │   ├── Dashboard      # Home Page
│       │   ├── VisualizerPage # Visualization Page
│       │   └── BenchmarkPage  # Benchmark Page
│       ├── services/          # API Services
│       │   ├── api.ts        # REST Client
│       │   └── sse.ts        # SSE Client
│       └── types/             # TypeScript Types
│
└── README.md
```

---

## 🔧 API Documentation

### Algorithm Operations

```http
POST /api/v1/operations
Content-Type: application/json

{
  "structure": "rbtree",     // rbtree | avltree | graph
  "operation": "insert",      // insert | delete | search | shortest_path
  "params": { "value": 42 }
}
```

**Response:**
```json
{
  "success": true,
  "steps": [
    {
      "type": "insert",
      "description": "Create new node 42 (red)",
      "nodeId": 0,
      "treeState": [...]
    }
  ],
  "finalTree": [...]
}
```

### Benchmarking

```http
POST /api/v1/benchmark/start
Content-Type: application/json

{
  "dataSize": 10000,
  "structures": ["hashmap", "rbtree", "avltree"],
  "operation": "insert"
}
```

**Response (SSE Stream):**
```
data: {"structure":"hashmap","progress":50,"duration":123.5,...}
data: {"structure":"hashmap","progress":100,"completed":true,...}
```

---

## 🛠️ Tech Stack

### Backend
- **Go** - High-performance backend language
- **Gin** - Lightweight HTTP framework
- **Goroutines** - Concurrent benchmarking

### Frontend
- **React 19** - UI Framework
- **TypeScript** - Type Safety
- **Vite** - Build Tool
- **Recharts** - Data Visualization Charts
- **Lucide React** - Icon Library
- **React Router** - Routing Management

---

## 📊 Supported Data Structures

| Data Structure | Insert | Delete | Search | Visualization | Benchmarking |
|----------------|--------|--------|--------|---------------|--------------|
| Red-Black Tree | ✅     | 🚧     | ✅     | ✅            | ✅           |
| AVL Tree       | ✅     | 🚧     | ✅     | ✅            | ✅           |
| Graph          | ✅     | ❌     | ❌     | ✅            | ❌           |
| HashMap        | ✅     | ❌     | ✅     | ❌            | ✅           |
| B-Tree         | 🚧     | ❌     | 🚧     | ❌            | ✅           |

✅ Implemented | 🚧 In Progress | ❌ Planned

---

## 🤝 Contributing

Issues and Pull Requests are welcome!

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Recharts](https://recharts.org/)
- [Lucide Icons](https://lucide.dev/)

---

<p align="center">
  Made with ❤️ for algorithm learners
</p>
