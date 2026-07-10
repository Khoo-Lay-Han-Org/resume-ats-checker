import logging

from pymilvus import FieldSchema, MilvusClient, CollectionSchema, DataType

logging.getLogger("pymilvus").setLevel(logging.CRITICAL)

client = MilvusClient("./database.db")
id_field = FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True)
vector_field = FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=50)
text_field = FieldSchema(name="text", dtype=DataType.VARCHAR, max_length=65535)
schema = CollectionSchema(fields=[id_field, vector_field, text_field])

try:
    client.create_collection(collection_name="data", schema=schema)
except:
    client.load_collection("data")
