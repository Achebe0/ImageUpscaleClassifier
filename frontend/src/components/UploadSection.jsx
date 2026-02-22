import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { toast } from 'react-toastify';
import { useImageUpload } from '../hooks/useImageUpload';
import './UploadSection.css';

const UploadSection = () => {
  const { isUploading, progress, result, error, uploadImage, reset } = useImageUpload();

  const onDrop = useCallback(async (acceptedFiles) => {
    if (acceptedFiles.length === 0) return;

    const file = acceptedFiles[0];
    
    // Validate file type
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      toast.error('Invalid file type. Please upload JPG, PNG, or WebP images.');
      return;
    }

    // Validate file size (max 100MB)
    const maxSize = 100 * 1024 * 1024;
    if (file.size > maxSize) {
      toast.error('File too large. Maximum size is 100MB.');
      return;
    }

    try {
      const response = await uploadImage(file);
      
      if (response.success) {
        toast.success('Image processed successfully!');
      } else {
        toast.error(response.message || 'Processing failed');
      }
    } catch (err) {
      toast.error(err.message || 'Upload failed');
    }
  }, [uploadImage]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'image/*': ['.jpeg', '.jpg', '.png', '.webp']
    },
    maxFiles: 1,
    disabled: isUploading
  });

  const handleReset = () => {
    reset();
  };

  const getStatusBadge = (status, folder) => {
    if (status === 'success') {
      if (folder === 'good_quality') {
        return <span className="badge badge-success">Good Quality ✓</span>;
      }
      return <span className="badge badge-success">Upscaled ✓</span>;
    }
    return <span className="badge badge-error">Failed ✗</span>;
  };

  return (
    <section className="upload-section">
      <h2 className="section-title">Upload Image</h2>
      
      {!result ? (
        <div
          {...getRootProps()}
          className={`dropzone ${isDragActive ? 'active' : ''} ${isUploading ? 'disabled' : ''}`}
        >
          <input {...getInputProps()} />
          
          {isUploading ? (
            <div className="upload-progress">
              <div className="progress-ring">
                <svg viewBox="0 0 36 36">
                  <path
                    className="progress-ring-bg"
                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                  />
                  <path
                    className="progress-ring-fill"
                    strokeDasharray={`${progress}, 100`}
                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                  />
                </svg>
                <span className="progress-text">{progress}%</span>
              </div>
              <p className="upload-status">Processing image...</p>
            </div>
          ) : (
            <div className="dropzone-content">
              <div className="upload-icon">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M21 15V19C21 19.5304 20.7893 20.0391 20.4142 20.4142C20.0391 20.7893 19.5304 21 19 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M17 8L12 3L7 8" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M12 3V15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </div>
              <p className="dropzone-text">
                {isDragActive ? 'Drop the image here...' : 'Drag & drop an image here, or click to select'}
              </p>
              <p className="dropzone-hint">
                Supports: JPG, PNG, WebP (Max 100MB)
              </p>
            </div>
          )}
        </div>
      ) : (
        <div className="result-card animate-fade-in">
          <div className="result-header">
            <h3>Processing Result</h3>
            {getStatusBadge(result.success && result.result?.status === 'success' ? 'success' : 'error', result.result?.folder)}
          </div>
          
          {result.result && (
            <div className="result-details">
              <div className="detail-row">
                <span className="detail-label">Original File:</span>
                <span className="detail-value">{result.result.original_key}</span>
              </div>
              
              <div className="detail-row">
                <span className="detail-label">Quality Score:</span>
                <span className="detail-value">
                  {(result.result.quality_score * 100).toFixed(1)}%
                </span>
              </div>
              
              {result.result.upscale_scale > 0 && (
                <div className="detail-row">
                  <span className="detail-label">Upscale Factor:</span>
                  <span className="detail-value">{result.result.upscale_scale}x</span>
                </div>
              )}
              
              {result.result.s3_url && (
                <div className="detail-row">
                  <span className="detail-label">Storage URL:</span>
                  <a 
                    href={result.result.s3_url} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="detail-link"
                  >
                    View in S3 ↗
                  </a>
                </div>
              )}
              
              {result.result.error_message && (
                <div className="detail-row error">
                  <span className="detail-label">Error:</span>
                  <span className="detail-value">{result.result.error_message}</span>
                </div>
              )}
            </div>
          )}
          
          <button onClick={handleReset} className="btn btn-primary">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style={{width: '16px', height: '16px'}}>
              <path d="M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" strokeWidth="2"/>
              <path d="M12 8V16" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
              <path d="M8 12H16" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </svg>
            Upload Another Image
          </button>
        </div>
      )}
      
      {error && (
        <div className="error-message">
          <p>{error}</p>
        </div>
      )}
    </section>
  );
};

export default UploadSection;
