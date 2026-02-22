import React, { useState } from 'react';
import './App.css';

function App() {
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  const handleFileSelect = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile && selectedFile.type.startsWith('image/')) {
      setFile(selectedFile);
      setError(null);
    } else {
      setError('Please select a valid image file');
    }
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.currentTarget.classList.add('drag-over');
  };

  const handleDragLeave = (e) => {
    e.currentTarget.classList.remove('drag-over');
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.currentTarget.classList.remove('drag-over');
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile && droppedFile.type.startsWith('image/')) {
      setFile(droppedFile);
      setError(null);
    } else {
      setError('Please drop a valid image file');
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setError('No file selected');
      return;
    }

    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const formData = new FormData();
      formData.append('image', file);

      const response = await fetch('http://localhost:8080/api/images/upload', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Upload failed');
      }

      const data = await response.json();
      setResult(data.result);
      setFile(null);
    } catch (err) {
      setError(err.message || 'Failed to upload image');
    } finally {
      setLoading(false);
    }
  };

  const reset = () => {
    setFile(null);
    setResult(null);
    setError(null);
  };

  return (
    <div className="container">
      <header>
        <h1>VisionCloud</h1>
        <p>Image Upscaling Pipeline</p>
      </header>

      <main>
        {!result ? (
          <div className="upload-section">
            <div
              className="upload-box"
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
            >
              <input
                type="file"
                id="fileInput"
                onChange={handleFileSelect}
                accept="image/*"
                disabled={loading}
                style={{ display: 'none' }}
              />
              <label htmlFor="fileInput">
                {loading ? (
                  <div className="spinner"></div>
                ) : (
                  <>
                    <p className="icon">📷</p>
                    <p className="text">
                      {file ? file.name : 'Drag image here or click to select'}
                    </p>
                  </>
                )}
              </label>
            </div>

            {error && <div className="error">{error}</div>}

            <button
              onClick={handleUpload}
              disabled={!file || loading}
              className="upload-btn"
            >
              {loading ? 'Processing...' : 'Upload & Process'}
            </button>
          </div>
        ) : (
          <div className="result-section">
            <h2>Processing Complete ✓</h2>
            <div className="result-info">
              <div className="info-row">
                <span className="label">File:</span>
                <span>{result.original_key}</span>
              </div>
              <div className="info-row">
                <span className="label">Status:</span>
                <span className={`status ${result.status}`}>
                  {result.status.toUpperCase()}
                </span>
              </div>
              <div className="info-row">
                <span className="label">Quality Score:</span>
                <span>
                  {result.quality_score.toFixed(2)} 
                  <div className="score-bar">
                    <div
                      className="score-fill"
                      style={{ width: `${result.quality_score * 100}%` }}
                    ></div>
                  </div>
                </span>
              </div>
              <div className="info-row">
                <span className="label">Destination:</span>
                <span className={`folder ${result.folder}`}>
                  {result.folder.replace(/_/g, ' ')}
                </span>
              </div>
              {result.s3_url && (
                <div className="info-row">
                  <span className="label">S3 URL:</span>
                  <a href={result.s3_url} target="_blank" rel="noopener noreferrer">
                    View on S3 ↗
                  </a>
                </div>
              )}
              {result.error_message && (
                <div className="error">{result.error_message}</div>
              )}
            </div>

            <button onClick={reset} className="reset-btn">
              Upload Another
            </button>
          </div>
        )}
      </main>

      <footer>
        <p>VisionCloud © 2026 | Go + React + AWS S3</p>
      </footer>
    </div>
  );
}

export default App;
