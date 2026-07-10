valkey:
	valkey-cli

tone-detection-model:
	cd ../ai-ml/tone_detection_model/local/ && uv run ./inference.py

translation-model:
	cd ../ai-ml && uv run python -m translation_model.local.inference

run:
	go run main.go
