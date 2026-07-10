from fastapi.responses import JSONResponse
from langdetect import detect
import torch
from sentence_transformers import util

from config import app
from api.typing import *
from api.model_and_tool import *


@app.get("/job-type-model-predict")
def job_type_model_predict(request: TextRequest):
    JobTypeModel.eval()

    feature = request.text

    tokenised_data = JobTypeTokeniser(
        text=feature,
        truncation=True,
        padding="max_length",
        max_length=350,
        return_tensor="pt",
    )
    attention_mask = tokenised_data["attention_mask"]
    tokenised_description = tokenised_data["input_ids"]

    class_indices_tensor = JobTypeModel(tokenised_description, attention_mask)
    class_index = class_indices_tensor.argmax(dim=1).item()

    job_type = JobTypeClassConverter[class_index]

    return job_type


@app.get("/resume-sections-model-predict")
def resume_sections_model_predict(request: TextRequest):
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
    entity_per_word = entity_per_word.squeeze(1)
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
    for i, prediction in enumerate(entity_per_word):
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


@app.post("/tone-detection-model-predict")
def tone_detection_model_predict(request: dict):
    phrase = request.get("phrase", "")
    vectorised_phrase = ToneDetectionVectorizer.transform([phrase]).toarray()
    result = ToneDetectionModel.predict(vectorised_phrase)
    return JSONResponse(content={"prediction": result.tolist()}, status_code=200)


@app.post("/translation-model-predict")
def translation_model_predict(request: TextRequest, tgt_lang="eng_Latn"):
    phrase = request.text

    detected = detect(phrase)
    src_lang = TranslationSpecifyLanguage.get(detected)

    if not src_lang:
        return False

    TranslationTokeniser.src_lang = src_lang
    inputs = TranslationTokeniser(phrase, return_tensors="pt", truncation=True)

    translated_tokens = TranslationModel.generate(
        **inputs,
        forced_bos_token_id=TranslationTokeniser.convert_tokens_to_ids(tgt_lang),
        max_length=256,
    )

    return JSONResponse(
        content={
            "prediction": TranslationTokeniser.batch_decode(
                translated_tokens, skip_special_tokens=True
            )[0]
        },
        status_code=200,
    )


@app.post("/skills-keyword-predict")
def skills_keyword_predict(request: TextRequest):
    text = request.text

    skills = SkillsKeywordPipeline(text)

    return skills


@app.post("/text-similarity-predict")
def text_similarity_predict(request: DoubleTextRequest):
    text1 = request.text1
    text2 = request.text2

    emb1 = TextSimilarityModel.encode(text1, convert_to_tensor=True)
    emb2 = TextSimilarityModel.encode(text2, convert_to_tensor=True)
    score = util.pytorch_cos_sim(emb1, emb2).item()

    return score
