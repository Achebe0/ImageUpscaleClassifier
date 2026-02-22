import api from './api';

export const imageService = {
  // Check backend health
  async checkHealth() {
    try {
      const response = await api.get('/api/health');
      return response.data;
    } catch (error) {
      throw new Error('Backend is not available');
    }
  },

  // Upload single image
  async uploadImage(file, onProgress) {
    const formData = new FormData();
    formData.append('image', file);

    try {
      const response = await api.post('/api/images/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          if (onProgress && progressEvent.total) {
            const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
            onProgress(progress);
          }
        },
      });
      return response.data;
    } catch (error) {
      if (error.response?.data?.error) {
        throw new Error(error.response.data.error);
      }
      throw new Error('Failed to upload image');
    }
  },

  // Get image info (placeholder - backend not fully implemented)
  async getImage(folder, filename) {
    try {
      const response = await api.get(`/api/images/${folder}/${filename}`);
      return response.data;
    } catch (error) {
      throw new Error('Failed to get image');
    }
  },

  // List processed images (placeholder - backend not fully implemented)
  async listImages(folder) {
    try {
      const response = await api.get(`/api/images/list/${folder}`);
      return response.data;
    } catch (error) {
      throw new Error('Failed to list images');
    }
  },
};

export default imageService;
