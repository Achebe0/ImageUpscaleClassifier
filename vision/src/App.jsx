import { useState } from 'react';
import './App.css';

function App() {
  const [file, setFile] = useState(null);
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleFileChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];
      
      // Simple format validation
      const validTypes = ['image/jpeg', 'image/png', 'image/webp'];
      if (!validTypes.includes(selectedFile.type)) {
        setError('Invalid file format. Please upload JPG, PNG, or WebP.');
        setFile(null);
        return;
      }

      setFile(selectedFile);
      setResult(null);
      setError(null);
    }
  };

  const handleUpload = async () => {
    if (!file) return;

    setLoading(true);
    setError(null);

    const formData = new FormData();
    formData.append('image', file);

    try {
      const response = await fetch('http://localhost:8080/api/upload', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Upload failed');
      }

      const data = await response.json();
      setResult(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h1>CloudSnap IQ</h1>
      <p className="subtitle">Upload an image to upscale and evaluate quality.</p>

      <div className="upload-section">
        <input 
          type="file" 
          accept="image/png, image/jpeg, image/webp" 
          onChange={handleFileChange} 
          className="file-input"
        />
        <button 
          onClick={handleUpload} 
          disabled={!file || loading}
          className="upload-btn"
        >
          {loading ? 'Processing...' : 'Upload & Analyze'}
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {result && (
        <div className="result-card">
          <h2>Analysis Result</h2>
          
          <div className="metric-row">
            <div className="metric">
              <span className="label">Original Quality:</span>
              <span className="value">{result.original_score.toFixed(1)}</span>
            </div>
            <div className="metric">
              <span className="label">Upscaled Quality:</span>
              <span className="value">{result.upscaled_score.toFixed(1)}</span>
            </div>
          </div>

          <div className="metric">
            <span className="label">Verdict:</span>
            <span className={`value ${result.quality_label === 'HighQuality' ? 'high' : 'low'}`}>
              {result.quality_label}
            </span>
          </div>

          <div className="metric">
            <span className="label">Bucket:</span>
            <span className="value bucket-name">{result.bucket}</span>
          </div>
          
          <div className="metric">
             <span className="label">Image ID:</span>
             <span className="value id">{result.image_id}</span>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
