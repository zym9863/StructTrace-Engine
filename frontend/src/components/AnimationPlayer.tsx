import { useState, useEffect, useCallback, useRef } from 'react';
import { Play, Pause, SkipBack, SkipForward, RotateCcw, ChevronLeft, ChevronRight } from 'lucide-react';
import type { Step } from '../types';
import './AnimationPlayer.css';

interface AnimationPlayerProps {
    steps: Step[];
    currentStep: number;
    onStepChange: (step: number) => void;
    onReset?: () => void;
}

export function AnimationPlayer({
    steps,
    currentStep,
    onStepChange,
    onReset,
}: AnimationPlayerProps) {
    const [isPlaying, setIsPlaying] = useState(false);
    const [speed, setSpeed] = useState(1);
    const intervalRef = useRef<number | null>(null);

    const totalSteps = steps.length;
    const currentStepData = steps[currentStep];

    // Auto-play logic
    useEffect(() => {
        if (isPlaying && currentStep < totalSteps - 1) {
            intervalRef.current = window.setInterval(() => {
                onStepChange(currentStep + 1);
            }, 1000 / speed);
        } else if (currentStep >= totalSteps - 1) {
            setIsPlaying(false);
        }

        return () => {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }
        };
    }, [isPlaying, currentStep, totalSteps, speed, onStepChange]);

    const handlePlayPause = useCallback(() => {
        if (currentStep >= totalSteps - 1) {
            onStepChange(0);
        }
        setIsPlaying(!isPlaying);
    }, [isPlaying, currentStep, totalSteps, onStepChange]);

    const handlePrevious = useCallback(() => {
        if (currentStep > 0) {
            onStepChange(currentStep - 1);
        }
    }, [currentStep, onStepChange]);

    const handleNext = useCallback(() => {
        if (currentStep < totalSteps - 1) {
            onStepChange(currentStep + 1);
        }
    }, [currentStep, totalSteps, onStepChange]);

    const handleFirst = useCallback(() => {
        onStepChange(0);
        setIsPlaying(false);
    }, [onStepChange]);

    const handleLast = useCallback(() => {
        onStepChange(totalSteps - 1);
        setIsPlaying(false);
    }, [totalSteps, onStepChange]);

    const handleReset = useCallback(() => {
        setIsPlaying(false);
        onStepChange(0);
        onReset?.();
    }, [onStepChange, onReset]);

    const handleSliderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(e.target.value, 10);
        onStepChange(value);
        setIsPlaying(false);
    };

    const getStepTypeIcon = (type: string) => {
        const icons: Record<string, string> = {
            insert: '➕',
            delete: '➖',
            compare: '🔍',
            rotate_left: '↩️',
            rotate_right: '↪️',
            color_change: '🎨',
            rebalance: '⚖️',
            found: '✅',
            not_found: '❌',
            complete: '🎉',
            visit: '👁️',
            select_node: '🎯',
            mark_visited: '✓',
            update_distance: '📏',
        };
        return icons[type] || '•';
    };

    if (totalSteps === 0) {
        return (
            <div className="animation-player empty">
                <p>执行操作以查看算法步骤</p>
            </div>
        );
    }

    return (
        <div className="animation-player">
            {/* Step Description */}
            <div className="step-info">
                <span className="step-icon">{getStepTypeIcon(currentStepData?.type || '')}</span>
                <div className="step-details">
                    <span className="step-type">{currentStepData?.type?.replace('_', ' ')}</span>
                    <span className="step-description">{currentStepData?.description}</span>
                </div>
                <span className="step-counter">{currentStep + 1} / {totalSteps}</span>
            </div>

            {/* Progress Bar */}
            <div className="progress-container">
                <input
                    type="range"
                    min={0}
                    max={totalSteps - 1}
                    value={currentStep}
                    onChange={handleSliderChange}
                    className="progress-slider"
                />
                <div
                    className="progress-fill"
                    style={{ width: `${(currentStep / Math.max(totalSteps - 1, 1)) * 100}%` }}
                />
            </div>

            {/* Controls */}
            <div className="controls">
                <div className="control-group">
                    <button onClick={handleReset} className="control-btn" title="重置">
                        <RotateCcw size={18} />
                    </button>
                </div>

                <div className="control-group main-controls">
                    <button onClick={handleFirst} className="control-btn" disabled={currentStep === 0} title="跳到开始">
                        <SkipBack size={18} />
                    </button>
                    <button onClick={handlePrevious} className="control-btn" disabled={currentStep === 0} title="上一步">
                        <ChevronLeft size={20} />
                    </button>
                    <button onClick={handlePlayPause} className="control-btn play-btn" title={isPlaying ? '暂停' : '播放'}>
                        {isPlaying ? <Pause size={24} /> : <Play size={24} />}
                    </button>
                    <button onClick={handleNext} className="control-btn" disabled={currentStep >= totalSteps - 1} title="下一步">
                        <ChevronRight size={20} />
                    </button>
                    <button onClick={handleLast} className="control-btn" disabled={currentStep >= totalSteps - 1} title="跳到结束">
                        <SkipForward size={18} />
                    </button>
                </div>

                <div className="control-group speed-control">
                    <label>速度</label>
                    <select value={speed} onChange={(e) => setSpeed(parseFloat(e.target.value))}>
                        <option value={0.5}>0.5x</option>
                        <option value={1}>1x</option>
                        <option value={2}>2x</option>
                        <option value={4}>4x</option>
                    </select>
                </div>
            </div>
        </div>
    );
}
