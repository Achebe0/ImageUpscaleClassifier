import { useState, useEffect, useCallback } from 'react';
import { imageService } from '../services/imageService';

export const useHealthCheck = (interval = 5000) => {
  const [isHealthy, setIsHealthy] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const checkHealth = useCallback(async () => {
    try {
      const response = await imageService.checkHealth();
      const healthy = response.status === 'healthy';
      setIsHealthy(healthy);
      setError(null);
    } catch (err) {
      setIsHealthy(false);
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    // Initial check
    checkHealth();

    // Set up interval
    const intervalId = setInterval(checkHealth, interval);

    return () => clearInterval(intervalId);
  }, [checkHealth, interval]);

  return { isHealthy, isLoading, error, checkHealth };
};

export default useHealthCheck;
