// Node colors for Red-Black Tree
export type NodeColor = 'red' | 'black';

// Step types for algorithm visualization
export type StepType =
  | 'insert'
  | 'delete'
  | 'rotate_left'
  | 'rotate_right'
  | 'color_change'
  | 'compare'
  | 'visit'
  | 'found'
  | 'not_found'
  | 'update_distance'
  | 'select_node'
  | 'mark_visited'
  | 'rebalance'
  | 'complete';

// Tree node snapshot for visualization
export interface TreeNodeSnapshot {
  id: number;
  value: number;
  color?: NodeColor;
  leftId?: number;
  rightId?: number;
  parentId?: number;
  height?: number;
  x?: number;
  y?: number;
}

// Graph node snapshot
export interface GraphNodeSnapshot {
  id: string;
  label: string;
  x: number;
  y: number;
  distance?: number;
  visited: boolean;
  inPath: boolean;
}

// Graph edge snapshot
export interface GraphEdgeSnapshot {
  from: string;
  to: string;
  weight: number;
  inPath: boolean;
  selected: boolean;
}

// Single step in algorithm execution
export interface Step {
  type: StepType;
  description: string;
  nodeId?: number;
  targetId?: number;
  value?: number;
  oldColor?: NodeColor;
  newColor?: NodeColor;
  treeState?: TreeNodeSnapshot[];
  graphNodes?: GraphNodeSnapshot[];
  graphEdges?: GraphEdgeSnapshot[];
  highlight?: number[];
}

// Operation result from API
export interface OperationResult {
  success: boolean;
  message?: string;
  steps: Step[];
  finalTree?: TreeNodeSnapshot[];
  finalGraph?: {
    nodes: GraphNodeSnapshot[];
    edges: GraphEdgeSnapshot[];
  };
}

// Operation request
export interface OperationRequest {
  structure: 'rbtree' | 'avltree' | 'graph';
  operation: 'insert' | 'delete' | 'search' | 'shortest_path';
  params: {
    value?: number;
    from?: string;
    to?: string;
  };
}

// Benchmark configuration
export interface BenchmarkConfig {
  dataSize: number;
  structures: string[];
  operation: 'insert' | 'search';
}

// Benchmark result
export interface BenchmarkResult {
  structure: string;
  operation: string;
  dataSize: number;
  duration: number;
  memoryUsed: number;
  opsPerSec: number;
  progress: number;
  completed: boolean;
}

// Animation player state
export interface AnimationState {
  isPlaying: boolean;
  currentStep: number;
  totalSteps: number;
  speed: number;
}
