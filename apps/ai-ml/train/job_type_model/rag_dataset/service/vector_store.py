from pymilvus import MilvusClient, CollectionSchema, FieldSchema, DataType
from pymilvus.orm.constants import IS_PRIMARY
from ..storage.route_config import VECTOR_DB_STORE

client = MilvusClient(VECTOR_DB_STORE)
id_field = FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True)
vector_field = FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=50)
text_field = FieldSchema(name="text", dtype=DataType.STRING)
schema = CollectionSchema(fields=[id_field, vector_field, text_field])

if client.has_collection("data"):
    client.load_collection("data")
else:
    client.create_collection(collection_name="data", schema=schema)
