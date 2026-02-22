import torch
from .model import UpscalerNet
from PIL import Image
import numpy as np
import os
import sys

class Upscaler:
    def __init__(self, model_path, device='cpu'):
        self.device = device
        self.model = UpscalerNet().to(device)
        
        if os.path.exists(model_path):
            self.model.load_state_dict(torch.load(model_path, map_location=device))
        
        self.model.eval()
    
    def upscale(self, input_path, output_path):
        """Upscale image from file path"""
        img = Image.open(input_path).convert('RGB')
        img_array = np.array(img)
        img_tensor = torch.tensor(img_array).permute(2, 0, 1).unsqueeze(0).float() / 255.0
        
        with torch.no_grad():
            upscaled = self.model(img_tensor.to(self.device))
        
        upscaled_np = (upscaled.squeeze(0).permute(1, 2, 0).cpu().numpy() * 255).astype(np.uint8)
        result = Image.fromarray(upscaled_np)
        result.save(output_path)
        return output_path

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python inference.py <input_image> <output_image>")
        sys.exit(1)
    
    upscaler = Upscaler('models/upscaler.pth')
    upscaler.upscale(sys.argv[1], sys.argv[2])
    print(f"Upscaled image saved to {sys.argv[2]}")