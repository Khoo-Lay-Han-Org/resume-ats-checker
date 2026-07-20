from langdetect import detect

from .config import TranslationModel, TranslationTokeniser, TranslationSpecifyLanguage


def predict(data: dict) -> dict:
    phrase = data.get("text") or data.get("phrase") or ""
    tgt_lang = data.get("tgt_lang", "eng_Latn")

    detected = detect(phrase)
    src_lang = TranslationSpecifyLanguage.get(detected)

    if not src_lang:
        return {"prediction": ""}

    TranslationTokeniser.src_lang = src_lang
    inputs = TranslationTokeniser(phrase, return_tensors="pt", truncation=True)

    translated_tokens = TranslationModel.generate(
        **inputs,
        forced_bos_token_id=TranslationTokeniser.convert_tokens_to_ids(tgt_lang),
        max_length=256,
    )

    return {
        "prediction": TranslationTokeniser.batch_decode(
            translated_tokens, skip_special_tokens=True
        )[0]
    }
