from transformers import AutoTokenizer

JobTypeModel = None
JobTypeTokeniser = AutoTokenizer.from_pretrained("google-bert/bert-base-cased")
JobTypeClassConverter = None
