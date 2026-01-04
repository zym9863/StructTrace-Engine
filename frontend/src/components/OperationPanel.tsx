import { useState } from 'react';
import { Send, RotateCcw, TreeDeciduous, GitBranch, Binary } from 'lucide-react';
import './OperationPanel.css';

interface OperationPanelProps {
    onExecute: (structure: string, operation: string, params: Record<string, unknown>) => void;
    onReset: () => void;
    isLoading?: boolean;
    selectedStructure?: string;
    onStructureChange?: (structure: string) => void;
}

export function OperationPanel({
    onExecute,
    onReset,
    isLoading = false,
    selectedStructure = 'rbtree',
    onStructureChange,
}: OperationPanelProps) {
    const [value, setValue] = useState('');
    const [operation, setOperation] = useState('insert');

    const structures = [
        { id: 'rbtree', name: '红黑树', icon: TreeDeciduous },
        { id: 'avltree', name: 'AVL树', icon: GitBranch },
        { id: 'graph', name: '图', icon: Binary },
    ];

    const operations: Record<string, { id: string; name: string }[]> = {
        rbtree: [
            { id: 'insert', name: '插入' },
            { id: 'search', name: '搜索' },
            { id: 'delete', name: '删除' },
        ],
        avltree: [
            { id: 'insert', name: '插入' },
            { id: 'search', name: '搜索' },
            { id: 'delete', name: '删除' },
        ],
        graph: [
            { id: 'insert', name: '添加节点' },
            { id: 'shortest_path', name: '最短路径' },
        ],
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!value.trim()) return;

        const numValue = parseInt(value, 10);
        if (isNaN(numValue)) return;

        onExecute(selectedStructure, operation, { value: numValue });
        setValue('');
    };

    const handleQuickInsert = () => {
        const randomValue = Math.floor(Math.random() * 100) + 1;
        onExecute(selectedStructure, 'insert', { value: randomValue });
    };

    return (
        <div className="operation-panel">
            <div className="panel-header">
                <h3>操作面板</h3>
                <button className="reset-btn" onClick={onReset} title="重置数据结构">
                    <RotateCcw size={16} />
                    重置
                </button>
            </div>

            {/* Structure Selection */}
            <div className="section">
                <label className="section-label">数据结构</label>
                <div className="structure-tabs">
                    {structures.map((s) => (
                        <button
                            key={s.id}
                            className={`structure-tab ${selectedStructure === s.id ? 'active' : ''}`}
                            onClick={() => onStructureChange?.(s.id)}
                        >
                            <s.icon size={18} />
                            {s.name}
                        </button>
                    ))}
                </div>
            </div>

            {/* Operation Selection */}
            <div className="section">
                <label className="section-label">操作类型</label>
                <div className="operation-buttons">
                    {operations[selectedStructure]?.map((op) => (
                        <button
                            key={op.id}
                            className={`operation-btn ${operation === op.id ? 'active' : ''}`}
                            onClick={() => setOperation(op.id)}
                        >
                            {op.name}
                        </button>
                    ))}
                </div>
            </div>

            {/* Value Input */}
            <form onSubmit={handleSubmit} className="section">
                <label className="section-label">输入值</label>
                <div className="input-group">
                    <input
                        type="number"
                        value={value}
                        onChange={(e) => setValue(e.target.value)}
                        placeholder="输入数字..."
                        className="value-input"
                        disabled={isLoading}
                    />
                    <button
                        type="submit"
                        className="submit-btn"
                        disabled={isLoading || !value.trim()}
                    >
                        {isLoading ? (
                            <span className="loading-spinner" />
                        ) : (
                            <Send size={18} />
                        )}
                    </button>
                </div>
            </form>

            {/* Quick Actions */}
            <div className="section">
                <label className="section-label">快捷操作</label>
                <div className="quick-actions">
                    <button
                        className="quick-btn"
                        onClick={handleQuickInsert}
                        disabled={isLoading}
                    >
                        🎲 随机插入
                    </button>
                    <button
                        className="quick-btn"
                        onClick={() => {
                            [10, 20, 30, 40, 50].forEach((v, i) => {
                                setTimeout(() => onExecute(selectedStructure, 'insert', { value: v }), i * 500);
                            });
                        }}
                        disabled={isLoading}
                    >
                        📚 批量插入
                    </button>
                </div>
            </div>
        </div>
    );
}
