from ..service.embedding_model import model


def embed_scraped_content(queries):
    all_content = []
    for query in queries:
        result = model.encode(query, truncate_dim=50, normalize_embeddings=True)

        all_content.append({"text": query, "embedding": result})

    return all_content
