import { Link } from 'react-router-dom';
import { TreeDeciduous, BarChart3, Zap, GitBranch, Binary, Database } from 'lucide-react';
import './Dashboard.css';

export function Dashboard() {
    const features = [
        {
            icon: TreeDeciduous,
            title: '红黑树',
            description: '自平衡二叉搜索树，保证 O(log n) 操作复杂度',
            color: '#ef4444',
        },
        {
            icon: GitBranch,
            title: 'AVL树',
            description: '严格平衡的二叉搜索树，高度差不超过1',
            color: '#3b82f6',
        },
        {
            icon: Binary,
            title: '图结构',
            description: '支持 Dijkstra 最短路径算法可视化',
            color: '#22c55e',
        },
        {
            icon: Database,
            title: 'HashMap',
            description: '哈希表实现，O(1) 平均复杂度',
            color: '#eab308',
        },
    ];

    return (
        <div className="dashboard">
            {/* Hero Section */}
            <section className="hero">
                <div className="hero-content">
                    <div className="hero-badge">🚀 数据结构可视化引擎</div>
                    <h1>
                        <span className="gradient-text">StructTrace</span>
                        <br />
                        Engine
                    </h1>
                    <p className="hero-subtitle">
                        实时追踪数据结构的演变过程，理解算法的每一步执行。
                        通过动态可视化和基准测试，深入学习数据结构的原理和性能。
                    </p>
                    <div className="hero-actions">
                        <Link to="/visualizer" className="btn btn-primary">
                            <TreeDeciduous size={20} />
                            开始可视化
                        </Link>
                        <Link to="/benchmark" className="btn btn-secondary">
                            <BarChart3 size={20} />
                            性能测试
                        </Link>
                    </div>
                </div>
                <div className="hero-visual">
                    <div className="tree-animation">
                        <div className="node root">42</div>
                        <div className="node-row">
                            <div className="node left">21</div>
                            <div className="node right red">63</div>
                        </div>
                        <div className="node-row">
                            <div className="node small">10</div>
                            <div className="node small">35</div>
                            <div className="node small">55</div>
                            <div className="node small">78</div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Features Section */}
            <section className="features">
                <div className="section-header">
                    <h2>支持的数据结构</h2>
                    <p>选择一种数据结构开始探索其内部工作原理</p>
                </div>
                <div className="features-grid">
                    {features.map((feature) => (
                        <Link
                            to="/visualizer"
                            key={feature.title}
                            className="feature-card"
                            style={{ '--accent': feature.color } as React.CSSProperties}
                        >
                            <div className="feature-icon">
                                <feature.icon size={28} />
                            </div>
                            <h3>{feature.title}</h3>
                            <p>{feature.description}</p>
                            <span className="feature-link">开始探索 →</span>
                        </Link>
                    ))}
                </div>
            </section>

            {/* Stats Section */}
            <section className="stats">
                <div className="stat-card">
                    <Zap className="stat-icon" size={32} />
                    <div className="stat-content">
                        <span className="stat-value">实时</span>
                        <span className="stat-label">可视化动画</span>
                    </div>
                </div>
                <div className="stat-card">
                    <BarChart3 className="stat-icon" size={32} />
                    <div className="stat-content">
                        <span className="stat-value">并发</span>
                        <span className="stat-label">基准测试</span>
                    </div>
                </div>
                <div className="stat-card">
                    <TreeDeciduous className="stat-icon" size={32} />
                    <div className="stat-content">
                        <span className="stat-value">4+</span>
                        <span className="stat-label">数据结构</span>
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className="cta">
                <div className="cta-content">
                    <h2>准备好开始学习了吗？</h2>
                    <p>选择一个功能开始探索数据结构的奥秘</p>
                    <div className="cta-buttons">
                        <Link to="/visualizer" className="btn btn-primary btn-large">
                            开始可视化之旅
                        </Link>
                    </div>
                </div>
            </section>
        </div>
    );
}
