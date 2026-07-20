import json
import sys

from . import (
    job_type_model,
    resume_section_model,
    tone_detection_model,
    translation_model,
    skills_keyword_model,
    text_similarity_model,
)

MODEL_FUNCTIONS = {
    "job-type-model-predict": job_type_model.predict,
    "resume-sections-model-predict": resume_section_model.predict,
    "tone-detection-model-predict": tone_detection_model.predict,
    "translation-model-predict": translation_model.predict,
    "skills-keyword-predict": skills_keyword_model.predict,
    "text-similarity-predict": text_similarity_model.predict,
}


def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "Model name required as argument"}))
        sys.exit(1)

    model_name = sys.argv[1]
    if model_name not in MODEL_FUNCTIONS:
        print(json.dumps({"error": f"Unknown model: {model_name}"}))
        sys.exit(1)

    input_data = json.loads(sys.stdin.read())
    result = MODEL_FUNCTIONS[model_name](input_data)
    print(json.dumps(result))


if __name__ == "__main__":
    main()
