from transformers import AutoTokenizer, AutoModelForTokenClassification, pipeline

model_name = "Nucha/Nucha_SkillNER_BERT"
tokenizer = AutoTokenizer.from_pretrained(model_name)
model = AutoModelForTokenClassification.from_pretrained(model_name)

ner_pipeline = pipeline(
    task="ner",  # type: ignore
    model=model,
    tokenizer=tokenizer,
    aggregation_strategy="simple",
)
