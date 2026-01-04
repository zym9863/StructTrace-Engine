# 🌳 StructTrace Engine

**数据结构演变追踪引擎** - 实时可视化数据结构操作过程，深入理解算法原理

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?logo=typescript)

---

## ✨ 特性

- 🎬 **逐帧动画可视化** - 观察每一步算法执行过程
- 🌲 **多种数据结构** - 红黑树、AVL树、图结构、HashMap
- ⚡ **并发基准测试** - 实时比较不同数据结构的性能
- 📊 **交互式图表** - 直观展示执行时间、内存使用和操作速度
- 🎨 **现代化 UI** - 精美的暗色主题，流畅的动画效果

## 🖼️ 预览

### 算法可视化
实时观察红黑树插入、旋转、变色等操作过程

### 性能基准测试
并发测试多种数据结构的插入/搜索性能

---

## 🚀 快速开始

### 前置要求

- **Go** 1.21+
- **Node.js** 18+
- **pnpm** 8+

### 安装与运行

```bash
# 克隆项目
git clone https://github.com/zym9863/StructTrace-Engine.git
cd StructTrace-Engine

# 启动后端
cd backend
go run main.go

# 新终端，启动前端
cd frontend
pnpm install
pnpm dev
```

访问 http://localhost:5173 开始使用

---

## 📁 项目结构

```
StructTrace-Engine/
├── backend/                    # Go 后端
│   ├── main.go                # 入口 + Gin 路由
│   ├── handlers/              # HTTP 处理器
│   │   ├── operation.go       # 算法操作 API
│   │   └── benchmark.go       # 基准测试 API + SSE
│   ├── datastructures/        # 数据结构实现
│   │   ├── rbtree.go         # 红黑树
│   │   ├── avltree.go        # AVL树
│   │   ├── graph.go          # 图 + Dijkstra
│   │   └── snapshot.go       # 快照结构定义
│   └── benchmark/             # 基准测试服务
│       └── runner.go         # Goroutine 并发测试
│
├── frontend/                   # React 前端
│   └── src/
│       ├── components/        # UI 组件
│       │   ├── TreeCanvas     # 树可视化
│       │   ├── AnimationPlayer# 动画播放器
│       │   ├── OperationPanel # 操作面板
│       │   └── BenchmarkChart # 基准图表
│       ├── pages/             # 页面
│       │   ├── Dashboard      # 主页
│       │   ├── VisualizerPage # 可视化页面
│       │   └── BenchmarkPage  # 基准测试页面
│       ├── services/          # API 服务
│       │   ├── api.ts        # REST 客户端
│       │   └── sse.ts        # SSE 客户端
│       └── types/             # TypeScript 类型
│
└── README.md
```

---

## 🔧 API 文档

### 算法操作

```http
POST /api/v1/operations
Content-Type: application/json

{
  "structure": "rbtree",     // rbtree | avltree | graph
  "operation": "insert",      // insert | delete | search | shortest_path
  "params": { "value": 42 }
}
```

**响应：**
```json
{
  "success": true,
  "steps": [
    {
      "type": "insert",
      "description": "创建新节点 42 (红色)",
      "nodeId": 0,
      "treeState": [...]
    }
  ],
  "finalTree": [...]
}
```

### 基准测试

```http
POST /api/v1/benchmark/start
Content-Type: application/json

{
  "dataSize": 10000,
  "structures": ["hashmap", "rbtree", "avltree"],
  "operation": "insert"
}
```

**响应 (SSE 流)：**
```
data: {"structure":"hashmap","progress":50,"duration":123.5,...}
data: {"structure":"hashmap","progress":100,"completed":true,...}
```

---

## 🛠️ 技术栈

### 后端
- **Go** - 高性能后端语言
- **Gin** - 轻量级 HTTP 框架
- **Goroutines** - 并发基准测试

### 前端
- **React 19** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Recharts** - 数据可视化图表
- **Lucide React** - 图标库
- **React Router** - 路由管理

---

## 📊 支持的数据结构

| 数据结构 | 插入 | 删除 | 搜索 | 可视化 | 基准测试 |
|---------|------|------|------|--------|---------|
| 红黑树   | ✅   | 🚧   | ✅   | ✅     | ✅      |
| AVL树    | ✅   | 🚧   | ✅   | ✅     | ✅      |
| 图       | ✅   | ❌   | ❌   | ✅     | ❌      |
| HashMap  | ✅   | ❌   | ✅   | ❌     | ✅      |
| B-Tree   | 🚧   | ❌   | 🚧   | ❌     | ✅      |

✅ 已实现 | 🚧 开发中 | ❌ 计划中

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

## 🙏 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Recharts](https://recharts.org/)
- [Lucide Icons](https://lucide.dev/)

---

<p align="center">
  Made with ❤️ for algorithm learners
</p>
