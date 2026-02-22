import { useState, useCallback } from 'react';
import { imageService } from '../services/imageService';

export const useImageUpload = () => {
  const [isUploading, setIsUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  const uploadImage = useCallback(async (file) => {
    setIsUploading(true);
    setProgress(0);
    setResult(null);
    setError(null);

    try {
      const response = await imageService.uploadImage(file, (progressValue) => {
        setProgress(progressValue);
      });
      
      setResult(response);
      return response;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setIsUploading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setIsUploading(false);
    setProgress(0);
    setResult(null);
    setError(null);
  }, []);

  return {
    isUploading,
    progress,
    result,
    error,
    uploadImage,
    reset,
  };
};

export default useImageUpload;
