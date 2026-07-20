import pickle as pk

from huggingface_hub import hf_hub_download
from safetensors.tensorflow import load_file as load_file_ts
import tensorflow as tf

ToneDetectionModel = tf.keras.Sequential(
    [
        tf.keras.layers.Input(shape=(1000,)),
        tf.keras.layers.Dense(128, activation="relu"),
        tf.keras.layers.Dropout(0.4),
        tf.keras.layers.Dense(64, activation="relu"),
        tf.keras.layers.Dropout(0.4),
        tf.keras.layers.Dense(32, activation="relu"),
        tf.keras.layers.Dense(1, activation="linear"),
    ]
)
ToneDetectionWeightPath = hf_hub_download(
    repo_id="includecctype/ToneDetection",
    filename="model_weights.safetensors",
)
ToneDetectionWeight = load_file_ts(ToneDetectionWeightPath)
for weight in ToneDetectionModel.weights:
    if weight.name in ToneDetectionWeight:
        value = ToneDetectionWeight[weight.name]
        if weight.shape == value.shape:
            weight.assign(value)
        else:
            print(f"Warning: Shape mismatch for {weight.name}: {weight.shape} vs {value.shape}")
    else:
        print(f"Warning: {weight.name} not found in safetensors file")
ToneDetectionVectorizerPath = hf_hub_download(
    repo_id="includecctype/ToneDetection",
    filename="vectorizer.pkl",
)
with open(ToneDetectionVectorizerPath, "rb") as f:
    ToneDetectionVectorizer = pk.load(f)
