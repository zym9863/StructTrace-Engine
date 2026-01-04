import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
    LineChart,
    Line,
} from 'recharts';
import type { BenchmarkResult } from '../types';
import './BenchmarkChart.css';

interface BenchmarkChartProps {
    results: BenchmarkResult[];
    progressResults?: BenchmarkResult[];
    chartType?: 'bar' | 'line';
}

const STRUCTURE_COLORS: Record<string, string> = {
    hashmap: '#22c55e',
    btree: '#eab308',
    rbtree: '#ef4444',
    avltree: '#3b82f6',
};

const STRUCTURE_NAMES: Record<string, string> = {
    hashmap: 'HashMap',
    btree: 'B-Tree',
    rbtree: '红黑树',
    avltree: 'AVL树',
};

export function BenchmarkChart({
    results,
    progressResults = [],
    chartType = 'bar',
}: BenchmarkChartProps) {
    // Combine completed and progress results
    const displayResults = results.length > 0 ? results : progressResults;

    if (displayResults.length === 0) {
        return (
            <div className="benchmark-chart empty">
                <div className="empty-icon">📊</div>
                <p>运行基准测试以查看结果</p>
            </div>
        );
    }

    // Format data for charts
    const timeData = displayResults.map((r) => ({
        name: STRUCTURE_NAMES[r.structure] || r.structure,
        time: Math.round(r.duration * 100) / 100,
        fill: STRUCTURE_COLORS[r.structure] || '#6c63ff',
    }));

    const memoryData = displayResults.map((r) => ({
        name: STRUCTURE_NAMES[r.structure] || r.structure,
        memory: Math.round(r.memoryUsed / 1024),
        fill: STRUCTURE_COLORS[r.structure] || '#6c63ff',
    }));

    const opsData = displayResults.map((r) => ({
        name: STRUCTURE_NAMES[r.structure] || r.structure,
        ops: Math.round(r.opsPerSec),
        fill: STRUCTURE_COLORS[r.structure] || '#6c63ff',
    }));

    const CustomTooltip = ({ active, payload, label }: any) => {
        if (active && payload && payload.length) {
            return (
                <div className="custom-tooltip">
                    <p className="tooltip-label">{label}</p>
                    {payload.map((entry: any, index: number) => (
                        <p key={index} style={{ color: entry.color }}>
                            {entry.name}: {entry.value.toLocaleString()}
                            {entry.name === 'time' ? ' ms' : entry.name === 'memory' ? ' KB' : ' ops/s'}
                        </p>
                    ))}
                </div>
            );
        }
        return null;
    };

    return (
        <div className="benchmark-chart">
            <div className="chart-grid">
                {/* Execution Time Chart */}
                <div className="chart-container">
                    <h4 className="chart-title">⏱️ 执行时间 (ms)</h4>
                    <ResponsiveContainer width="100%" height={200}>
                        <BarChart data={timeData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                            <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                            <XAxis dataKey="name" stroke="#8888aa" fontSize={12} />
                            <YAxis stroke="#8888aa" fontSize={12} />
                            <Tooltip content={<CustomTooltip />} />
                            <Bar dataKey="time" radius={[4, 4, 0, 0]}>
                                {timeData.map((entry, index) => (
                                    <Bar key={index} dataKey="time" fill={entry.fill} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </div>

                {/* Memory Usage Chart */}
                <div className="chart-container">
                    <h4 className="chart-title">💾 内存使用 (KB)</h4>
                    <ResponsiveContainer width="100%" height={200}>
                        <BarChart data={memoryData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                            <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                            <XAxis dataKey="name" stroke="#8888aa" fontSize={12} />
                            <YAxis stroke="#8888aa" fontSize={12} />
                            <Tooltip content={<CustomTooltip />} />
                            <Bar dataKey="memory" radius={[4, 4, 0, 0]}>
                                {memoryData.map((entry, index) => (
                                    <Bar key={index} dataKey="memory" fill={entry.fill} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </div>

                {/* Ops/Second Chart */}
                <div className="chart-container full-width">
                    <h4 className="chart-title">🚀 操作速度 (ops/s)</h4>
                    <ResponsiveContainer width="100%" height={200}>
                        <BarChart data={opsData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                            <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                            <XAxis dataKey="name" stroke="#8888aa" fontSize={12} />
                            <YAxis stroke="#8888aa" fontSize={12} />
                            <Tooltip content={<CustomTooltip />} />
                            <Legend />
                            <Bar dataKey="ops" name="操作/秒" radius={[4, 4, 0, 0]}>
                                {opsData.map((entry, index) => (
                                    <Bar key={index} dataKey="ops" fill={entry.fill} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>

            {/* Results Summary */}
            <div className="results-summary">
                <h4>📋 测试结果摘要</h4>
                <div className="summary-cards">
                    {displayResults.map((result) => (
                        <div
                            key={result.structure}
                            className="summary-card"
                            style={{ borderColor: STRUCTURE_COLORS[result.structure] }}
                        >
                            <div className="card-header">
                                <span
                                    className="structure-badge"
                                    style={{ background: STRUCTURE_COLORS[result.structure] }}
                                >
                                    {STRUCTURE_NAMES[result.structure] || result.structure}
                                </span>
                                {result.completed ? (
                                    <span className="status-badge completed">✓ 完成</span>
                                ) : (
                                    <span className="status-badge running">{result.progress}%</span>
                                )}
                            </div>
                            <div className="card-stats">
                                <div className="stat">
                                    <span className="stat-value">{result.duration.toFixed(2)}</span>
                                    <span className="stat-label">ms</span>
                                </div>
                                <div className="stat">
                                    <span className="stat-value">
                                        {(result.memoryUsed / 1024).toFixed(1)}
                                    </span>
                                    <span className="stat-label">KB</span>
                                </div>
                                <div className="stat">
                                    <span className="stat-value">
                                        {result.opsPerSec.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                                    </span>
                                    <span className="stat-label">ops/s</span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
