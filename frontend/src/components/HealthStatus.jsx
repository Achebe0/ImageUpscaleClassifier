import React from 'react';
import { useHealthCheck } from '../hooks/useHealthCheck';
import './HealthStatus.css';

const HealthStatus = () => {
  const { isHealthy, isLoading } = useHealthCheck(5000);

  return (
    <div className="health-status">
      <div className={`status-indicator ${isLoading ? 'loading' : isHealthy ? 'healthy' : 'unhealthy'}`}>
        <span className="status-dot"></span>
        <span className="status-text">
          {isLoading ? 'Checking...' : isHealthy ? 'Backend Connected' : 'Backend Disconnected'}
        </span>
      </div>
    </div>
  );
};

export default HealthStatus;
