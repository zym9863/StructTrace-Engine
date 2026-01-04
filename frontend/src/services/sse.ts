import type { BenchmarkResult, BenchmarkConfig } from '../types';

const API_BASE = 'http://localhost:8080/api/v1';

export interface SSECallbacks {
    onProgress: (result: BenchmarkResult) => void;
    onComplete: (results: BenchmarkResult[]) => void;
    onError: (error: string) => void;
}

export class BenchmarkSSEClient {
    private eventSource: EventSource | null = null;
    private results: BenchmarkResult[] = [];

    start(config: BenchmarkConfig, callbacks: SSECallbacks): void {
        this.stop(); // Stop any existing connection
        this.results = [];

        // Build query params
        const params = new URLSearchParams({
            dataSize: config.dataSize.toString(),
            structures: config.structures.join(','),
            operation: config.operation,
        });

        // Use POST for SSE by creating a hidden form or using fetch with streaming
        // For simplicity, we'll use a workaround with GET-like EventSource
        this.eventSource = new EventSource(`${API_BASE}/benchmark/stream?${params}`);

        this.eventSource.onmessage = (event) => {
            try {
                const result: BenchmarkResult = JSON.parse(event.data);

                if (result.completed) {
                    // Update final result for this structure
                    const existingIndex = this.results.findIndex(r => r.structure === result.structure);
                    if (existingIndex >= 0) {
                        this.results[existingIndex] = result;
                    } else {
                        this.results.push(result);
                    }

                    // Check if all structures are complete
                    const allComplete = config.structures.every(
                        struct => this.results.some(r => r.structure === struct && r.completed)
                    );

                    if (allComplete) {
                        callbacks.onComplete(this.results);
                        this.stop();
                    }
                }

                callbacks.onProgress(result);
            } catch (e) {
                console.error('Failed to parse SSE message:', e);
            }
        };

        this.eventSource.onerror = () => {
            callbacks.onError('Connection to benchmark server lost');
            this.stop();
        };
    }

    stop(): void {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }
    }

    isConnected(): boolean {
        return this.eventSource !== null && this.eventSource.readyState === EventSource.OPEN;
    }
}

// Singleton instance
export const benchmarkSSE = new BenchmarkSSEClient();
