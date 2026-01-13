import { useState, useCallback } from 'react';
import { TreeCanvas } from '../components/TreeCanvas';
import { GraphCanvas } from '../components/GraphCanvas';
import { AnimationPlayer } from '../components/AnimationPlayer';
import { OperationPanel } from '../components/OperationPanel';
import { executeOperation, resetStructure } from '../services/api';
import type { Step, TreeNodeSnapshot, GraphNodeSnapshot, GraphEdgeSnapshot, OperationRequest } from '../types';
import './VisualizerPage.css';

export function VisualizerPage() {
    const [structure, setStructure] = useState<string>('rbtree');
    const [steps, setSteps] = useState<Step[]>([]);
    const [currentStep, setCurrentStep] = useState<number>(0);
    const [treeNodes, setTreeNodes] = useState<TreeNodeSnapshot[]>([]);
    const [graphNodes, setGraphNodes] = useState<GraphNodeSnapshot[]>([]);
    const [graphEdges, setGraphEdges] = useState<GraphEdgeSnapshot[]>([]);
    const [highlightedNodes, setHighlightedNodes] = useState<number[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Check if we're in graph mode
    const isGraphMode = structure === 'graph';

    const handleExecute = useCallback(
        async (structureType: string, operation: string, params: Record<string, unknown>) => {
            setIsLoading(true);
            setError(null);

            try {
                const request: OperationRequest = {
                    structure: structureType as OperationRequest['structure'],
                    operation: operation as OperationRequest['operation'],
                    params: params as OperationRequest['params'],
                };

                const result = await executeOperation(request);

                if (result.success) {
                    setSteps(result.steps);
                    setCurrentStep(0);

                    // Handle tree structures
                    if (structureType !== 'graph') {
                        // Set initial tree state if available
                        if (result.steps.length > 0 && result.steps[0].treeState) {
                            setTreeNodes(result.steps[0].treeState);
                        }

                        // Update final tree state
                        if (result.finalTree) {
                            setTreeNodes(result.finalTree);
                        }
                    } else {
                        // Handle graph structure
                        // Set initial graph state if available
                        if (result.steps.length > 0) {
                            const firstStep = result.steps[0];
                            if (firstStep.graphNodes) {
                                setGraphNodes(firstStep.graphNodes);
                            }
                            if (firstStep.graphEdges) {
                                setGraphEdges(firstStep.graphEdges);
                            }
                        }

                        // Update final graph state
                        if (result.finalGraph) {
                            setGraphNodes(result.finalGraph.nodes);
                            setGraphEdges(result.finalGraph.edges);
                        }
                    }
                } else {
                    setError(result.message || '操作失败');
                }
            } catch (err) {
                setError(err instanceof Error ? err.message : '连接服务器失败');
            } finally {
                setIsLoading(false);
            }
        },
        []
    );

    const handleReset = useCallback(async () => {
        try {
            await resetStructure();
            setSteps([]);
            setCurrentStep(0);
            setTreeNodes([]);
            setGraphNodes([]);
            setGraphEdges([]);
            setHighlightedNodes([]);
            setError(null);
        } catch (err) {
            setError('重置失败');
        }
    }, []);

    const handleStepChange = useCallback(
        (stepIndex: number) => {
            setCurrentStep(stepIndex);

            const step = steps[stepIndex];
            if (step) {
                // Update tree visualization based on current step
                if (step.treeState) {
                    setTreeNodes(step.treeState);
                }

                // Update graph visualization based on current step
                if (step.graphNodes) {
                    setGraphNodes(step.graphNodes);
                }
                if (step.graphEdges) {
                    setGraphEdges(step.graphEdges);
                }

                // Update highlighted nodes
                if (step.highlight) {
                    setHighlightedNodes(step.highlight);
                } else if (step.nodeId !== undefined) {
                    setHighlightedNodes([step.nodeId]);
                } else {
                    setHighlightedNodes([]);
                }
            }
        },
        [steps]
    );

    return (
        <div className="visualizer-page">
            <header className="page-header">
                <div className="header-content">
                    <h1>🌳 数据结构可视化</h1>
                    <p>实时观察数据结构的操作过程，理解算法的每一步</p>
                </div>
            </header>

            <main className="page-content">
                <div className="visualizer-layout">
                    {/* Left Panel - Controls */}
                    <aside className="controls-panel">
                        <OperationPanel
                            onExecute={handleExecute}
                            onReset={handleReset}
                            isLoading={isLoading}
                            selectedStructure={structure}
                            onStructureChange={setStructure}
                        />
                    </aside>

                    {/* Main Area - Visualization */}
                    <div className="visualization-area">
                        {error && (
                            <div className="error-banner">
                                <span>⚠️ {error}</span>
                                <button onClick={() => setError(null)}>✕</button>
                            </div>
                        )}

                        <div className="canvas-container">
                            {isGraphMode ? (
                                <GraphCanvas
                                    nodes={graphNodes}
                                    edges={graphEdges}
                                    animated={true}
                                />
                            ) : (
                                <TreeCanvas
                                    nodes={treeNodes}
                                    highlightedNodes={highlightedNodes}
                                    showLabels={true}
                                    animated={true}
                                />
                            )}
                        </div>

                        <div className="player-container">
                            <AnimationPlayer
                                steps={steps}
                                currentStep={currentStep}
                                onStepChange={handleStepChange}
                                onReset={handleReset}
                            />
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}
