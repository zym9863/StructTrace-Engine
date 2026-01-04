import type { OperationRequest, OperationResult, BenchmarkConfig } from '../types';

const API_BASE = 'http://localhost:8080/api/v1';

// Generic fetch wrapper with error handling
async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options?.headers,
        },
    });

    if (!response.ok) {
        const error = await response.text();
        throw new Error(error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
}

// Execute a data structure operation
export async function executeOperation(request: OperationRequest): Promise<OperationResult> {
    return fetchAPI<OperationResult>('/operations', {
        method: 'POST',
        body: JSON.stringify(request),
    });
}

// Reset data structure state
export async function resetStructure(): Promise<{ message: string }> {
    return fetchAPI<{ message: string }>('/reset', {
        method: 'POST',
    });
}

// Start benchmark with SSE
export function startBenchmark(config: BenchmarkConfig): EventSource {
    const eventSource = new EventSource(
        `${API_BASE}/benchmark/start?dataSize=${config.dataSize}&structures=${config.structures.join(',')}&operation=${config.operation}`
    );
    return eventSource;
}

// Stop running benchmark
export async function stopBenchmark(): Promise<{ message: string }> {
    return fetchAPI<{ message: string }>('/benchmark/stop', {
        method: 'POST',
    });
}

// Get benchmark status
export async function getBenchmarkStatus(): Promise<{ running: boolean }> {
    return fetchAPI<{ running: boolean }>('/benchmark/status', {
        method: 'GET',
    });
}

// Health check
export async function healthCheck(): Promise<{ status: string; service: string }> {
    return fetchAPI<{ status: string; service: string }>('/health', {
        method: 'GET',
    });
}
