import torch
import torch.nn as nn

class ResidualBlock(nn.Module):
    def __init__(self, channels):
        super(ResidualBlock, self).__init__()
        self.conv1 = nn.Conv2d(channels, channels, kernel_size=3, padding=1)
        self.conv2 = nn.Conv2d(channels, channels, kernel_size=3, padding=1)
        self.relu = nn.ReLU(inplace=True)
        self.bn1 = nn.BatchNorm2d(channels)
        self.bn2 = nn.BatchNorm2d(channels)
    
    def forward(self, x):
        residual = x
        out = self.relu(self.bn1(self.conv1(x)))
        out = self.bn2(self.conv2(out))
        out += residual
        return self.relu(out)

class UpscalerCNN(nn.Module):
    def __init__(self, scale_factor=4, num_residual_blocks=16):
        super(UpscalerCNN, self).__init__()
        self.scale_factor = scale_factor
        
        # Initial feature extraction
        self.conv1 = nn.Conv2d(3, 64, kernel_size=9, padding=4)
        self.relu = nn.ReLU(inplace=True)
        
        # Residual blocks
        self.residual_blocks = nn.Sequential(
            *[ResidualBlock(64) for _ in range(num_residual_blocks)]
        )
        
        # Conv after residual blocks
        self.conv2 = nn.Conv2d(64, 64, kernel_size=3, padding=1)
        self.bn = nn.BatchNorm2d(64)
        
        # Upsampling layers
        self.upsample_layers = nn.Sequential()
        for _ in range(2):  # 4x upscaling = 2^2
            self.upsample_layers.add_module('conv', nn.Conv2d(64, 256, kernel_size=3, padding=1))
            self.upsample_layers.add_module('pixel_shuffle', nn.PixelShuffle(2))
        
        # Final reconstruction
        self.conv3 = nn.Conv2d(64, 3, kernel_size=9, padding=4)
    
    def forward(self, x):
        # Initial feature extraction
        initial = self.relu(self.conv1(x))
        
        # Residual blocks
        residual = self.residual_blocks(initial)
        residual = self.bn(self.conv2(residual))
        residual += initial
        
        # Upsampling
        upsampled = self.upsample_layers(residual)
        
        # Reconstruction
        output = self.conv3(upsampled)
        return output