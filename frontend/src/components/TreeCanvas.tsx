import { useEffect, useRef, useState, useCallback } from 'react';
import type { TreeNodeSnapshot, NodeColor } from '../types';
import './TreeCanvas.css';

interface TreeCanvasProps {
    nodes: TreeNodeSnapshot[];
    highlightedNodes?: number[];
    showLabels?: boolean;
    animated?: boolean;
}

interface NodePosition {
    id: number;
    value: number;
    color: NodeColor;
    x: number;
    y: number;
    leftId?: number;
    rightId?: number;
}

export function TreeCanvas({
    nodes,
    highlightedNodes = [],
    showLabels = true,
    animated = true,
}: TreeCanvasProps) {
    const svgRef = useRef<SVGSVGElement>(null);
    const [dimensions, setDimensions] = useState({ width: 800, height: 500 });
    const [nodePositions, setNodePositions] = useState<Map<number, NodePosition>>(new Map());

    // Calculate node positions with proper tree layout
    const calculateLayout = useCallback(() => {
        if (nodes.length === 0) return new Map<number, NodePosition>();

        const positions = new Map<number, NodePosition>();
        const nodeMap = new Map(nodes.map(n => [n.id, n]));

        // Find root (node without parent)
        const root = nodes.find(n => n.parentId === undefined || n.parentId === null);
        if (!root) {
            // Fallback: use provided x, y positions
            nodes.forEach(node => {
                positions.set(node.id, {
                    id: node.id,
                    value: node.value,
                    color: node.color || 'black',
                    x: node.x || 400,
                    y: node.y || 50,
                    leftId: node.leftId,
                    rightId: node.rightId,
                });
            });
            return positions;
        }

        // BFS to calculate positions
        const horizontalSpacing = dimensions.width / (Math.pow(2, Math.ceil(Math.log2(nodes.length + 1)) + 1));
        const verticalSpacing = 80;
        const padding = 60;

        function layoutNode(
            nodeId: number,
            depth: number,
            xMin: number,
            xMax: number
        ) {
            const node = nodeMap.get(nodeId);
            if (!node) return;

            const x = (xMin + xMax) / 2;
            const y = depth * verticalSpacing + padding;

            positions.set(nodeId, {
                id: node.id,
                value: node.value,
                color: node.color || 'black',
                x,
                y,
                leftId: node.leftId,
                rightId: node.rightId,
            });

            if (node.leftId !== undefined) {
                layoutNode(node.leftId, depth + 1, xMin, x);
            }
            if (node.rightId !== undefined) {
                layoutNode(node.rightId, depth + 1, x, xMax);
            }
        }

        layoutNode(root.id, 0, 0, dimensions.width);
        return positions;
    }, [nodes, dimensions.width]);

    useEffect(() => {
        setNodePositions(calculateLayout());
    }, [calculateLayout]);

    // Resize handler
    useEffect(() => {
        const handleResize = () => {
            if (svgRef.current?.parentElement) {
                const rect = svgRef.current.parentElement.getBoundingClientRect();
                setDimensions({ width: rect.width, height: rect.height });
            }
        };

        handleResize();
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    // Render edges
    const edges: JSX.Element[] = [];
    nodePositions.forEach((node) => {
        if (node.leftId !== undefined) {
            const child = nodePositions.get(node.leftId);
            if (child) {
                edges.push(
                    <line
                        key={`edge-${node.id}-${child.id}`}
                        className={`tree-edge ${animated ? 'animated' : ''}`}
                        x1={node.x}
                        y1={node.y}
                        x2={child.x}
                        y2={child.y}
                    />
                );
            }
        }
        if (node.rightId !== undefined) {
            const child = nodePositions.get(node.rightId);
            if (child) {
                edges.push(
                    <line
                        key={`edge-${node.id}-${child.id}`}
                        className={`tree-edge ${animated ? 'animated' : ''}`}
                        x1={node.x}
                        y1={node.y}
                        x2={child.x}
                        y2={child.y}
                    />
                );
            }
        }
    });

    // Render nodes
    const nodeElements: JSX.Element[] = [];
    nodePositions.forEach((node) => {
        const isHighlighted = highlightedNodes.includes(node.id);
        nodeElements.push(
            <g
                key={`node-${node.id}`}
                className={`tree-node ${animated ? 'animated' : ''} ${isHighlighted ? 'highlighted' : ''}`}
                transform={`translate(${node.x}, ${node.y})`}
            >
                <circle
                    r={24}
                    className={`node-circle ${node.color}`}
                />
                {isHighlighted && (
                    <circle r={30} className="highlight-ring" />
                )}
                {showLabels && (
                    <text
                        dy="0.35em"
                        textAnchor="middle"
                        className="node-label"
                    >
                        {node.value}
                    </text>
                )}
            </g>
        );
    });

    return (
        <div className="tree-canvas-container">
            {nodes.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-icon">🌳</div>
                    <p>空树 - 执行插入操作开始构建</p>
                </div>
            ) : (
                <svg
                    ref={svgRef}
                    className="tree-canvas"
                    viewBox={`0 0 ${dimensions.width} ${dimensions.height}`}
                    preserveAspectRatio="xMidYMid meet"
                >
                    <defs>
                        <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
                            <feGaussianBlur stdDeviation="3" result="coloredBlur" />
                            <feMerge>
                                <feMergeNode in="coloredBlur" />
                                <feMergeNode in="SourceGraphic" />
                            </feMerge>
                        </filter>
                    </defs>
                    <g className="edges">{edges}</g>
                    <g className="nodes">{nodeElements}</g>
                </svg>
            )}
        </div>
    );
}
