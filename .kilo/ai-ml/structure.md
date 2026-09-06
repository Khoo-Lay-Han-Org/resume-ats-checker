# AI-ML Structure Convention

## Please understand that this is for apps/ai-ml/ ONLY

### Root Files and Folders Organisation

1) pyproject.toml defines project metadata and dependencies
2) ruff.toml configures linting rules
3) black.toml configures formatting rules
4) pyrightconfig.json configures type checking
5) .python-version pins the Python version
6) .venv/ is local only and not committed
7) .ruff_cache/ is local only and not committed

### Training Directory Organisation

1) train/ contains all model training workflows
2) Each model has its own subdirectory under train/
3) Model subdirectories contain: config.py, training notebooks, readme.md, dataset/, model/
4) Dataset references should use HuggingFace when possible
5) Model artifacts are stored in model/ subdirectories
6) Training notebooks (.ipynb) are for experimentation and documentation