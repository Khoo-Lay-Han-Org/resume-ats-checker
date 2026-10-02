# AI-ML Best Coding Practice

## Please understand that this is for apps/ai-ml/ ONLY

1) Ask the user for more clarity and avoid guessing
2) Look at the existing related codes for convention examples
3) Do a study on best practices before coding
4) Always devise a clear plan and ask user for approval of the plan first

## Python / ML Conventions

1) Type hints are required for function signatures
2) Use ruff for linting and black for formatting
3) Use pyright for static type checking
4) Notebooks (.ipynb) are for training and experimentation only
5) Production code should be in .py modules, not notebooks
6) Model configs load from config.py with graceful fallbacks
7) Training artifacts belong in train/<model_name>/
