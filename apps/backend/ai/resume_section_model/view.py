import torch

from ..typing import TextRequest
from .config import (
    ResumeSectionModel,
    ResumeSectionTokeniser,
    ResumeSectionClassConverter,
)


def predict(data: dict) -> dict:
    request = TextRequest(**data)

    ResumeSectionModel.eval()

    raw_feature = request.text

    tokenised_data = ResumeSectionTokeniser(
        text=raw_feature,
        truncation=True,
        padding="max_length",
        max_length=400,
        return_tensors="pt",
    )

    feature = tokenised_data["input_ids"]

    entity_per_word = ResumeSectionModel(feature)
    entity_per_word = entity_per_word.squeeze(0)
    entity_per_word = torch.argmax(entity_per_word, dim=-1)

    personal_info_section = []
    summary_section = []
    skills_section = []
    experience_section = []
    education_section = []
    certificates_section = []
    objective_section = []

    splitted_raw_data = raw_feature.split()
    previous_class = ""
    for i in range(min(len(entity_per_word), len(splitted_raw_data))):
        prediction = entity_per_word[i]
        prediction_item = prediction.item()

        if prediction_item != -100:
            current_class = ResumeSectionClassConverter[int(prediction_item)]

            match current_class:
                case "personal_info":
                    if current_class == previous_class:
                        latest_data = personal_info_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        personal_info_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        personal_info_section.append(splitted_raw_data[i])
                case "summary":
                    if current_class == previous_class:
                        latest_data = summary_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        summary_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        summary_section.append(splitted_raw_data[i])
                case "skills":
                    if current_class == previous_class:
                        latest_data = skills_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        skills_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        skills_section.append(splitted_raw_data[i])
                case "experience":
                    if current_class == previous_class:
                        latest_data = experience_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        experience_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        experience_section.append(splitted_raw_data[i])
                case "education":
                    if current_class == previous_class:
                        latest_data = education_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        education_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        education_section.append(splitted_raw_data[i])
                case "certificates":
                    if current_class == previous_class:
                        latest_data = certificates_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        certificates_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        certificates_section.append(splitted_raw_data[i])
                case "objective":
                    if current_class == previous_class:
                        latest_data = objective_section[-1]
                        appended_latest_data = latest_data + splitted_raw_data[i]
                        objective_section[-1] = appended_latest_data
                    elif current_class != previous_class:
                        objective_section.append(splitted_raw_data[i])

    return {
        "personal_info": personal_info_section,
        "summary": summary_section,
        "skills": skills_section,
        "experience": experience_section,
        "education": education_section,
        "certificates": certificates_section,
        "objective": objective_section,
    }
