import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader
from .model import UpscalerNet

class Trainer:
    def __init__(self, model, device='cpu', learning_rate=0.001):
        self.model = model.to(device)
        self.device = device
        self.criterion = nn.MSELoss()
        self.optimizer = optim.Adam(model.parameters(), lr=learning_rate)
    
    def train(self, train_loader, epochs=10):
        self.model.train()
        
        for epoch in range(epochs):
            total_loss = 0
            for batch_idx, (low_quality, high_quality) in enumerate(train_loader):
                low_quality = low_quality.to(self.device)
                high_quality = high_quality.to(self.device)
                
                self.optimizer.zero_grad()
                output = self.model(low_quality)
                loss = self.criterion(output, high_quality)
                loss.backward()
                self.optimizer.step()
                
                total_loss += loss.item()
            
            avg_loss = total_loss / len(train_loader)
            print(f'Epoch {epoch+1}/{epochs}, Loss: {avg_loss:.4f}')
    
    def save_model(self, path):
        torch.save(self.model.state_dict(), path)
        print(f'Model saved to {path}')