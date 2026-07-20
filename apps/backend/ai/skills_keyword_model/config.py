from transformers import AutoTokenizer, pipeline

SkillsKeywordPipeline = pipeline(
    task="ner",
    model="Nucha/Nucha_SkillNER_BERT",
    tokenizer=AutoTokenizer.from_pretrained("Nucha/Nucha_SkillNER_BERT"),
    aggregation_strategy="simple",
)
