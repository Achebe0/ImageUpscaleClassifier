import React from 'react';
import './ImageCard.css';

const ImageCard = ({ image }) => {
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getStatusBadge = () => {
    if (image.folder === 'good_quality') {
      return <span className="badge badge-success">Good Quality</span>;
    } else if (image.folder === 'upscaled') {
      return <span className="badge badge-success">Upscaled {image.upscale_scale}x</span>;
    } else {
      return <span className="badge badge-error">Failed</span>;
    }
  };

  const getQualityColor = (score) => {
    if (score >= 0.8) return '#10b981';
    if (score >= 0.5) return '#f59e0b';
    return '#ef4444';
  };

  return (
    <div className="image-card">
      <div className="card-header">
        <div className="file-info">
          <div className="file-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M14 2H6C5.46957 2 4.96086 2.21071 4.58579 2.58579C4.21071 2.96086 4 3.46957 4 4V20C4 20.5304 4.21071 21.0391 4.58579 21.4142C4.96086 21.7893 5.46957 22 6 22H18C18.5304 22 19.0391 21.7893 19.4142 21.4142C19.7893 21.0391 20 20.5304 20 20V8L14 2Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              <path d="M14 2V8H20" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              <path d="M12 18V12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              <path d="M9 15L12 12L15 15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </div>
          <div className="file-details">
            <h4 className="file-name" title={image.original_key}>
              {image.original_key}
            </h4>
            <span className="file-date">{formatDate(image.processed_at)}</span>
          </div>
        </div>
        {getStatusBadge()}
      </div>
      
      <div className="card-body">
        <div className="quality-meter">
          <div className="quality-header">
            <span className="quality-label">Quality Score</span>
            <span 
              className="quality-score"
              style={{ color: getQualityColor(image.quality_score) }}
            >
              {(image.quality_score * 100).toFixed(1)}%
            </span>
          </div>
          <div className="quality-bar">
            <div 
              className="quality-fill"
              style={{ 
                width: `${image.quality_score * 100}%`,
                backgroundColor: getQualityColor(image.quality_score)
              }}
            />
          </div>
        </div>
        
        {image.upscale_scale > 0 && (
          <div className="upscale-info">
            <div className="info-row">
              <span className="info-label">Upscale Factor:</span>
              <span className="info-value">{image.upscale_scale}x</span>
            </div>
          </div>
        )}
      </div>
      
      <div className="card-footer">
        {image.s3_url && image.s3_url !== '#' ? (
          <a 
            href={image.s3_url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="btn btn-primary btn-sm"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style={{width: '16px', height: '16px'}}>
              <path d="M18 13V19C18 19.5304 17.7893 20.0391 17.4142 20.4142C17.0391 20.7893 16.5304 21 16 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V8C3 7.46957 3.21071 6.96086 3.58579 6.58579C3.96086 6.21071 4.46957 6 5 6H11" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              <path d="M15 3H21V9" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              <path d="M10 14L21 3" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
            View in S3
          </a>
        ) : (
          <button className="btn btn-secondary btn-sm" disabled>
            Not Available
          </button>
        )}
      </div>
    </div>
  );
};

export default ImageCard;
