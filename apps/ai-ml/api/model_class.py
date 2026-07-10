import torch.nn as nn
import tensorflow as tf


class ResumeSectionsModel(nn.Module):
    def __init__(self):
        super().__init__()

        self.embed = nn.Embedding(30000, 20)
        self.lstm = nn.LSTM(20, 20, batch_first=True)
        self.linear = nn.Linear(20, 7)

    def forward(self, feature):
        x = self.embed(feature)
        x, _ = self.lstm(x)
        logits = self.linear(x)

        return logits
