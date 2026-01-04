import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';
import { Home, TreeDeciduous, BarChart3 } from 'lucide-react';
import { Dashboard } from './pages/Dashboard';
import { VisualizerPage } from './pages/VisualizerPage';
import { BenchmarkPage } from './pages/BenchmarkPage';
import './App.css';

function Navigation() {
  const location = useLocation();

  const navItems = [
    { path: '/', icon: Home, label: '首页' },
    { path: '/visualizer', icon: TreeDeciduous, label: '可视化' },
    { path: '/benchmark', icon: BarChart3, label: '基准测试' },
  ];

  return (
    <nav className="main-nav">
      <div className="nav-brand">
        <span className="brand-icon">🌳</span>
        <span className="brand-text">StructTrace</span>
      </div>
      <div className="nav-links">
        {navItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className={`nav-link ${location.pathname === item.path ? 'active' : ''}`}
          >
            <item.icon size={18} />
            <span>{item.label}</span>
          </Link>
        ))}
      </div>
      <div className="nav-actions">
        <a
          href="https://github.com"
          target="_blank"
          rel="noopener noreferrer"
          className="github-btn"
        >
          ⭐ GitHub
        </a>
      </div>
    </nav>
  );
}

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <Navigation />
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/visualizer" element={<VisualizerPage />} />
          <Route path="/benchmark" element={<BenchmarkPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
