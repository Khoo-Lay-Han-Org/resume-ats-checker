from huggingface_hub import hf_hub_download
from safetensors.torch import load_file as load_file_pt
from transformers import AutoTokenizer

from .model_class import ResumeSectionsModel

ResumeSectionModel = ResumeSectionsModel()
ResumeSectionWeightPath = hf_hub_download(
    repo_id="includecctype/ResumeSectionsExtractor",
    filename="model.safetensors",
)
ResumeSectionWeight = load_file_pt(ResumeSectionWeightPath)
ResumeSectionModel.load_state_dict(ResumeSectionWeight)
ResumeSectionTokeniser = AutoTokenizer.from_pretrained("google-bert/bert-base-cased")
ResumeSectionClassConverter = {
    0: "personal_info",
    1: "summary",
    2: "skills",
    3: "experience",
    4: "education",
    5: "certificates",
    6: "objective",
}
