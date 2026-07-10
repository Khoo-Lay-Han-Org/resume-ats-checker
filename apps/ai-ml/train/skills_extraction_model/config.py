import torch
from transformers import AutoTokenizer
import pickle as pk

tokeniser = AutoTokenizer.from_pretrained("google-bert/bert-base-cased")

model = torch.load("model/model.pth")

with open("skill_encoder.pkl", "rb") as file:
    mlb = pk.load(file)
