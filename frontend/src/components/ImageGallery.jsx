import React, { useState } from 'react';
import ImageCard from './ImageCard';
import './ImageGallery.css';

// Mock data for demonstration (will be replaced with real API data)
const mockImages = {
  good_quality: [
    {
      id: 1,
      original_key: 'sample-high-res.jpg',
      status: 'success',
      folder: 'good_quality',
      quality_score: 0.92,
      processed_at: new Date().toISOString(),
      s3_url: '#'
    }
  ],
  upscaled: [
    {
      id: 2,
      original_key: 'sample-low-res.jpg',
      status: 'success',
      folder: 'upscaled',
      quality_score: 0.45,
      upscale_scale: 2,
      processed_at: new Date().toISOString(),
      s3_url: '#'
    }
  ],
  couldnt_upscale: []
};

const ImageGallery = () => {
  const [activeTab, setActiveTab] = useState('good_quality');
  const [images] = useState(mockImages);

  const tabs = [
    { id: 'good_quality', label: 'Good Quality', count: images.good_quality.length },
    { id: 'upscaled', label: 'Upscaled', count: images.upscaled.length },
    { id: 'couldnt_upscale', label: 'Couldn\'t Upscale', count: images.couldnt_upscale.length }
  ];

  const currentImages = images[activeTab] || [];

  return (
    <section className="image-gallery">
      <h2 className="section-title">Processed Images</h2>
      
      <div className="gallery-tabs">
        {tabs.map(tab => (
          <button
            key={tab.id}
            className={`tab ${activeTab === tab.id ? 'active' : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            <span className="tab-label">{tab.label}</span>
            <span className="tab-count">{tab.count}</span>
          </button>
        ))}
      </div>
      
      <div className="gallery-content">
        {currentImages.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="3" y="3" width="18" height="18" rx="2" stroke="currentColor" strokeWidth="2"/>
                <circle cx="8.5" cy="8.5" r="1.5" fill="currentColor"/>
                <path d="M21 15L16 10L5 21" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </div>
            <p className="empty-text">No images in this category yet</p>
            <p className="empty-hint">Upload an image to see it here</p>
          </div>
        ) : (
          <div className="image-grid">
            {currentImages.map(image => (
              <ImageCard key={image.id} image={image} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
};

export default ImageGallery;
