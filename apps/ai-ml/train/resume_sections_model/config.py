from transformers import AutoTokenizer
import torch

try:
    model = torch.load("model/model.pth")
except:
    model = None

try:
    tokeniser = AutoTokenizer.from_pretrained("google-bert/bert-base-cased")
except:
    tokeniser = None


idx_to_class = {
    0: "personal_info",
    1: "summary",
    2: "skills",
    3: "experience",
    4: "education",
    5: "certificates",
    6: "objective",
}
