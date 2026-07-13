import tensorflow as tf
import pickle

model = None
try:
    model = tf.keras.models.load_model("model_package/model.keras")
except:
    model = None

vectorizer = None
try:
    vectorizer = pickle.load(open("model_package/vectorizer.pkl", "rb"))
except:
    vectorizer = None
