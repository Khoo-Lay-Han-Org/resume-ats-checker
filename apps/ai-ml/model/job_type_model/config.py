import torch
from transformers import AutoTokenizer
import pickle as pk

try:
    model = torch.load("model/model.pth")
except:
    model = None

try:
    tokeniser = AutoTokenizer.from_pretrained("google-bert/bert-base-cased")
except:
    tokeniser = None

job_types = {}

try:
    with open("dataset/convert_prediction.pkl", "rb") as file:
        job_types = pk.load(file)
except:
    job_types = None
