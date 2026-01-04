import { useState, useCallback } from 'react';
import { Play, Square, Settings } from 'lucide-react';
import { BenchmarkChart } from '../components/BenchmarkChart';
import type { BenchmarkResult, BenchmarkConfig } from '../types';
import './BenchmarkPage.css';

const API_BASE = 'http://localhost:8080/api/v1';

export function BenchmarkPage() {
    const [config, setConfig] = useState<BenchmarkConfig>({
        dataSize: 10000,
        structures: ['hashmap', 'rbtree', 'avltree'],
        operation: 'insert',
    });
    const [isRunning, setIsRunning] = useState(false);
    const [progressResults, setProgressResults] = useState<BenchmarkResult[]>([]);
    const [completedResults, setCompletedResults] = useState<BenchmarkResult[]>([]);
    const [error, setError] = useState<string | null>(null);

    const structures = [
        { id: 'hashmap', name: 'HashMap', color: '#22c55e' },
        { id: 'btree', name: 'B-Tree', color: '#eab308' },
        { id: 'rbtree', name: '红黑树', color: '#ef4444' },
        { id: 'avltree', name: 'AVL树', color: '#3b82f6' },
    ];

    const dataSizes = [
        { value: 1000, label: '1K' },
        { value: 10000, label: '10K' },
        { value: 50000, label: '50K' },
        { value: 100000, label: '100K' },
    ];

    const handleStart = useCallback(async () => {
        setIsRunning(true);
        setError(null);
        setProgressResults([]);
        setCompletedResults([]);

        try {
            const response = await fetch(`${API_BASE}/benchmark/start`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(config),
            });

            if (!response.ok) {
                throw new Error('启动基准测试失败');
            }

            // Handle SSE response
            const reader = response.body?.getReader();
            if (!reader) {
                throw new Error('无法读取响应流');
            }

            const decoder = new TextDecoder();
            const completed: BenchmarkResult[] = [];

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);
                const lines = chunk.split('\n');

                for (const line of lines) {
                    if (line.startsWith('data: ')) {
                        try {
                            const result: BenchmarkResult = JSON.parse(line.slice(6));

                            if (result.completed) {
                                completed.push(result);
                                setCompletedResults([...completed]);
                            }

                            setProgressResults((prev) => {
                                const existing = prev.findIndex((r) => r.structure === result.structure);
                                if (existing >= 0) {
                                    const updated = [...prev];
                                    updated[existing] = result;
                                    return updated;
                                }
                                return [...prev, result];
                            });
                        } catch {
                            // Skip malformed data
                        }
                    }
                }
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : '基准测试失败');
        } finally {
            setIsRunning(false);
        }
    }, [config]);

    const handleStop = useCallback(async () => {
        try {
            await fetch(`${API_BASE}/benchmark/stop`, { method: 'POST' });
        } catch {
            // Ignore stop errors
        }
        setIsRunning(false);
    }, []);

    const toggleStructure = (structureId: string) => {
        setConfig((prev) => ({
            ...prev,
            structures: prev.structures.includes(structureId)
                ? prev.structures.filter((s) => s !== structureId)
                : [...prev.structures, structureId],
        }));
    };

    return (
        <div className="benchmark-page">
            <header className="page-header">
                <div className="header-content">
                    <h1>⚡ 性能基准测试</h1>
                    <p>比较不同数据结构的性能表现，了解各自的优势和适用场景</p>
                </div>
            </header>

            <main className="page-content">
                <div className="benchmark-layout">
                    {/* Configuration Panel */}
                    <aside className="config-panel">
                        <div className="panel-section">
                            <h3>
                                <Settings size={18} />
                                测试配置
                            </h3>

                            {/* Data Size */}
                            <div className="config-group">
                                <label>数据规模</label>
                                <div className="size-options">
                                    {dataSizes.map((size) => (
                                        <button
                                            key={size.value}
                                            className={`size-btn ${config.dataSize === size.value ? 'active' : ''}`}
                                            onClick={() => setConfig((prev) => ({ ...prev, dataSize: size.value }))}
                                            disabled={isRunning}
                                        >
                                            {size.label}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            {/* Structures */}
                            <div className="config-group">
                                <label>数据结构</label>
                                <div className="structure-options">
                                    {structures.map((s) => (
                                        <label
                                            key={s.id}
                                            className={`structure-checkbox ${config.structures.includes(s.id) ? 'checked' : ''}`}
                                            style={{ '--accent-color': s.color } as React.CSSProperties}
                                        >
                                            <input
                                                type="checkbox"
                                                checked={config.structures.includes(s.id)}
                                                onChange={() => toggleStructure(s.id)}
                                                disabled={isRunning}
                                            />
                                            <span className="checkbox-indicator" />
                                            <span className="checkbox-label">{s.name}</span>
                                        </label>
                                    ))}
                                </div>
                            </div>

                            {/* Operation */}
                            <div className="config-group">
                                <label>操作类型</label>
                                <div className="operation-options">
                                    <button
                                        className={`operation-btn ${config.operation === 'insert' ? 'active' : ''}`}
                                        onClick={() => setConfig((prev) => ({ ...prev, operation: 'insert' }))}
                                        disabled={isRunning}
                                    >
                                        插入
                                    </button>
                                    <button
                                        className={`operation-btn ${config.operation === 'search' ? 'active' : ''}`}
                                        onClick={() => setConfig((prev) => ({ ...prev, operation: 'search' }))}
                                        disabled={isRunning}
                                    >
                                        搜索
                                    </button>
                                </div>
                            </div>

                            {/* Action Buttons */}
                            <div className="action-buttons">
                                {isRunning ? (
                                    <button className="action-btn stop" onClick={handleStop}>
                                        <Square size={18} />
                                        停止测试
                                    </button>
                                ) : (
                                    <button
                                        className="action-btn start"
                                        onClick={handleStart}
                                        disabled={config.structures.length === 0}
                                    >
                                        <Play size={18} />
                                        开始测试
                                    </button>
                                )}
                            </div>

                            {error && <div className="error-message">⚠️ {error}</div>}
                        </div>

                        {/* Progress Indicators */}
                        {isRunning && progressResults.length > 0 && (
                            <div className="panel-section progress-section">
                                <h3>📊 测试进度</h3>
                                <div className="progress-list">
                                    {progressResults.map((result) => (
                                        <div key={result.structure} className="progress-item">
                                            <span className="progress-label">
                                                {structures.find((s) => s.id === result.structure)?.name || result.structure}
                                            </span>
                                            <div className="progress-bar">
                                                <div
                                                    className="progress-fill"
                                                    style={{
                                                        width: `${result.progress}%`,
                                                        background: structures.find((s) => s.id === result.structure)?.color,
                                                    }}
                                                />
                                            </div>
                                            <span className="progress-value">{result.progress}%</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </aside>

                    {/* Results Area */}
                    <div className="results-area">
                        <BenchmarkChart
                            results={completedResults}
                            progressResults={progressResults}
                        />
                    </div>
                </div>
            </main>
        </div>
    );
}
