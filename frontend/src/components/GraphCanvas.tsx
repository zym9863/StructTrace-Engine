import React, { useEffect, useRef, useState, useCallback } from 'react';
import type { GraphNodeSnapshot, GraphEdgeSnapshot } from '../types';
import './GraphCanvas.css';

interface GraphCanvasProps {
    nodes: GraphNodeSnapshot[];
    edges: GraphEdgeSnapshot[];
    animated?: boolean;
}

interface NodePosition {
    id: string;
    label: string;
    x: number;
    y: number;
    distance?: number;
    visited: boolean;
    inPath: boolean;
}

export function GraphCanvas({
    nodes,
    edges,
    animated = true,
}: GraphCanvasProps) {
    const svgRef = useRef<SVGSVGElement>(null);
    const [dimensions, setDimensions] = useState({ width: 800, height: 500 });
    const [nodePositions, setNodePositions] = useState<Map<string, NodePosition>>(new Map());

    // Calculate node positions - use provided positions or auto-layout
    const calculateLayout = useCallback(() => {
        if (nodes.length === 0) return new Map<string, NodePosition>();

        const positions = new Map<string, NodePosition>();
        const centerX = dimensions.width / 2;
        const centerY = dimensions.height / 2;
        const radius = Math.min(dimensions.width, dimensions.height) / 3;

        nodes.forEach((node, index) => {
            // If node has valid coordinates, use them (scaled to fit viewport)
            let x = node.x;
            let y = node.y;

            // Scale positions if they're outside our viewport or use circular layout
            if (x === 0 && y === 0) {
                // Auto-layout: arrange nodes in a circle
                const angle = (2 * Math.PI * index) / nodes.length - Math.PI / 2;
                x = centerX + radius * Math.cos(angle);
                y = centerY + radius * Math.sin(angle);
            } else {
                // Scale and center existing positions
                const padding = 100;
                const scaleX = (dimensions.width - 2 * padding) / 600;
                const scaleY = (dimensions.height - 2 * padding) / 400;
                x = x * scaleX + padding;
                y = y * scaleY + padding / 2;
            }

            positions.set(node.id, {
                id: node.id,
                label: node.label,
                x,
                y,
                distance: node.distance,
                visited: node.visited,
                inPath: node.inPath,
            });
        });

        return positions;
    }, [nodes, dimensions]);

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
    const edgeElements: React.ReactElement[] = [];
    const renderedEdges = new Set<string>(); // Avoid duplicate edges for undirected graph

    edges.forEach((edge) => {
        const edgeKey = [edge.from, edge.to].sort().join('-');
        if (renderedEdges.has(edgeKey)) return;
        renderedEdges.add(edgeKey);

        const fromNode = nodePositions.get(edge.from);
        const toNode = nodePositions.get(edge.to);

        if (fromNode && toNode) {
            const midX = (fromNode.x + toNode.x) / 2;
            const midY = (fromNode.y + toNode.y) / 2;

            // Calculate angle for weight label offset
            const dx = toNode.x - fromNode.x;
            const dy = toNode.y - fromNode.y;
            const angle = Math.atan2(dy, dx);
            const labelOffset = 15;
            const labelX = midX - labelOffset * Math.sin(angle);
            const labelY = midY + labelOffset * Math.cos(angle);

            edgeElements.push(
                <g key={`edge-${edge.from}-${edge.to}`} className="graph-edge-group">
                    <line
                        className={`graph-edge ${animated ? 'animated' : ''} ${edge.inPath ? 'in-path' : ''} ${edge.selected ? 'selected' : ''}`}
                        x1={fromNode.x}
                        y1={fromNode.y}
                        x2={toNode.x}
                        y2={toNode.y}
                    />
                    <text
                        x={labelX}
                        y={labelY}
                        className="edge-weight"
                        textAnchor="middle"
                        dominantBaseline="middle"
                    >
                        {edge.weight}
                    </text>
                </g>
            );
        }
    });

    // Render nodes
    const nodeElements: React.ReactElement[] = [];
    nodePositions.forEach((node) => {
        const statusClass = node.inPath ? 'in-path' : node.visited ? 'visited' : '';

        nodeElements.push(
            <g
                key={`node-${node.id}`}
                className={`graph-node ${animated ? 'animated' : ''} ${statusClass}`}
                transform={`translate(${node.x}, ${node.y})`}
            >
                {/* Glow effect for path nodes */}
                {node.inPath && (
                    <circle r={36} className="node-glow" />
                )}

                {/* Main node circle */}
                <circle
                    r={28}
                    className="node-circle"
                />

                {/* Node label */}
                <text
                    dy="0.1em"
                    textAnchor="middle"
                    className="node-label"
                >
                    {node.label}
                </text>

                {/* Distance indicator */}
                {node.distance !== undefined && (
                    <text
                        dy="1.5em"
                        textAnchor="middle"
                        className="node-distance"
                    >
                        d={node.distance}
                    </text>
                )}
            </g>
        );
    });

    return (
        <div className="graph-canvas-container">
            {nodes.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-icon">🔗</div>
                    <p>空图 - 添加节点开始构建</p>
                </div>
            ) : (
                <svg
                    ref={svgRef}
                    className="graph-canvas"
                    viewBox={`0 0 ${dimensions.width} ${dimensions.height}`}
                    preserveAspectRatio="xMidYMid meet"
                >
                    <defs>
                        {/* Gradient for path edges */}
                        <linearGradient id="pathGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                            <stop offset="0%" stopColor="#4ff2c1" />
                            <stop offset="100%" stopColor="#6aa0ff" />
                        </linearGradient>

                        {/* Glow filter for highlighted elements */}
                        <filter id="graphGlow" x="-50%" y="-50%" width="200%" height="200%">
                            <feGaussianBlur stdDeviation="4" result="coloredBlur" />
                            <feMerge>
                                <feMergeNode in="coloredBlur" />
                                <feMergeNode in="SourceGraphic" />
                            </feMerge>
                        </filter>

                        {/* Node gradient */}
                        <radialGradient id="nodeGradient" cx="30%" cy="30%">
                            <stop offset="0%" stopColor="#6aa0ff" />
                            <stop offset="100%" stopColor="#1f3a5f" />
                        </radialGradient>

                        {/* Visited node gradient */}
                        <radialGradient id="visitedGradient" cx="30%" cy="30%">
                            <stop offset="0%" stopColor="#f7b267" />
                            <stop offset="100%" stopColor="#b96e2a" />
                        </radialGradient>

                        {/* Path node gradient */}
                        <radialGradient id="pathNodeGradient" cx="30%" cy="30%">
                            <stop offset="0%" stopColor="#7af7d2" />
                            <stop offset="100%" stopColor="#4ff2c1" />
                        </radialGradient>
                    </defs>

                    {/* Render edges first (below nodes) */}
                    <g className="edges">{edgeElements}</g>

                    {/* Render nodes on top */}
                    <g className="nodes">{nodeElements}</g>
                </svg>
            )}
        </div>
    );
}
