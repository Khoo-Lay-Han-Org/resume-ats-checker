from .config import ToneDetectionModel, ToneDetectionVectorizer


def predict(data: dict) -> dict:
    phrase = data.get("phrase", "")
    vectorised_phrase = ToneDetectionVectorizer.transform([phrase]).toarray()
    result = ToneDetectionModel.predict(vectorised_phrase)
    return {"prediction": result.tolist()}
