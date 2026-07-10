import logging

from pymilvus import FieldSchema, MilvusClient, CollectionSchema, DataType

logging.getLogger("pymilvus").setLevel(logging.CRITICAL)

client = MilvusClient("./database.db")

vector_field = FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=1000)

text_field = FieldSchema(name="text", dtype=DataType.VARCHAR, max_length=65536)

schema = CollectionSchema(fields=[vector_field, text_field])

try:
    client.create_collection(collection_name="data", schema=schema)
except:
    client.load_collection("data")
