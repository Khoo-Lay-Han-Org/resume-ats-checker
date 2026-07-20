from ..typing import TextRequest
from .config import JobTypeModel, JobTypeTokeniser, JobTypeClassConverter


def predict(data: dict) -> str:
    request = TextRequest(**data)

    JobTypeModel.eval()

    tokenised_data = JobTypeTokeniser(
        text=request.text,
        truncation=True,
        padding="max_length",
        max_length=350,
        return_tensors="pt",
    )
    attention_mask = tokenised_data["attention_mask"]
    tokenised_description = tokenised_data["input_ids"]

    class_indices_tensor = JobTypeModel(tokenised_description, attention_mask)
    class_index = class_indices_tensor.argmax(dim=1).item()

    return JobTypeClassConverter[class_index]
