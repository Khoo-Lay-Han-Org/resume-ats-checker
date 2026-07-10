from transformers import AutoTokenizer, AutoModelForSeq2SeqLM

tokenizer = AutoTokenizer.from_pretrained("facebook/nllb-200-1.3B")

model = AutoModelForSeq2SeqLM.from_pretrained("facebook/nllb-200-1.3B")

LANG_MAP = {
    # European - Romance
    "fr": "fra_Latn",  # French
    "es": "spa_Latn",  # Spanish
    "pt": "por_Latn",  # Portuguese
    "it": "ita_Latn",  # Italian
    "ro": "ron_Latn",  # Romanian
    "ca": "cat_Latn",  # Catalan
    "gl": "glg_Latn",  # Galician
    "oc": "oci_Latn",  # Occitan
    # European - Germanic
    "en": "eng_Latn",  # English
    "de": "deu_Latn",  # German
    "nl": "nld_Latn",  # Dutch
    "sv": "swe_Latn",  # Swedish
    "da": "dan_Latn",  # Danish
    "no": "nob_Latn",  # Norwegian
    "af": "afr_Latn",  # Afrikaans
    "is": "isl_Latn",  # Icelandic
    "lb": "ltz_Latn",  # Luxembourgish
    "yi": "ydd_Hebr",  # Yiddish
    # European - Slavic
    "ru": "rus_Cyrl",  # Russian
    "uk": "ukr_Cyrl",  # Ukrainian
    "pl": "pol_Latn",  # Polish
    "cs": "ces_Latn",  # Czech
    "sk": "slk_Latn",  # Slovak
    "bg": "bul_Cyrl",  # Bulgarian
    "hr": "hrv_Latn",  # Croatian
    "sr": "srp_Cyrl",  # Serbian
    "sl": "slv_Latn",  # Slovenian
    "mk": "mkd_Cyrl",  # Macedonian
    "bs": "bos_Latn",  # Bosnian
    "be": "bel_Cyrl",  # Belarusian
    # European - Baltic
    "lt": "lit_Latn",  # Lithuanian
    "lv": "lvs_Latn",  # Latvian
    # European - Finno-Ugric
    "fi": "fin_Latn",  # Finnish
    "et": "est_Latn",  # Estonian
    "hu": "hun_Latn",  # Hungarian
    # European - Celtic
    "cy": "cym_Latn",  # Welsh
    "ga": "gle_Latn",  # Irish
    "gd": "gla_Latn",  # Scottish Gaelic
    # European - Other
    "el": "ell_Grek",  # Greek
    "sq": "als_Latn",  # Albanian
    "hy": "hye_Armn",  # Armenian
    "ka": "kat_Geor",  # Georgian
    "eu": "eus_Latn",  # Basque
    "mt": "mlt_Latn",  # Maltese
    # Middle Eastern & Central Asian
    "ar": "arb_Arab",  # Arabic
    "fa": "pes_Arab",  # Persian/Farsi
    "ur": "urd_Arab",  # Urdu
    "he": "heb_Hebr",  # Hebrew
    "tr": "tur_Latn",  # Turkish
    "az": "azj_Latn",  # Azerbaijani
    "kk": "kaz_Cyrl",  # Kazakh
    "ky": "kir_Cyrl",  # Kyrgyz
    "uz": "uzn_Latn",  # Uzbek
    "tk": "tuk_Latn",  # Turkmen
    "tg": "tgk_Cyrl",  # Tajik
    "ps": "pbt_Arab",  # Pashto
    "ku": "kmr_Latn",  # Kurdish (Kurmanji)
    # South Asian
    "hi": "hin_Deva",  # Hindi
    "bn": "ben_Beng",  # Bengali
    "pa": "pan_Guru",  # Punjabi
    "gu": "guj_Gujr",  # Gujarati
    "mr": "mar_Deva",  # Marathi
    "ne": "npi_Deva",  # Nepali
    "si": "sin_Sinh",  # Sinhala
    "ta": "tam_Taml",  # Tamil
    "te": "tel_Telu",  # Telugu
    "kn": "kan_Knda",  # Kannada
    "ml": "mal_Mlym",  # Malayalam
    "or": "ory_Orya",  # Odia
    "as": "asm_Beng",  # Assamese
    "sd": "snd_Arab",  # Sindhi
    "dz": "dzo_Tibt",  # Dzongkha
    # Southeast Asian
    "ms": "zsm_Latn",  # Malay
    "id": "ind_Latn",  # Indonesian
    "tl": "tgl_Latn",  # Filipino/Tagalog
    "vi": "vie_Latn",  # Vietnamese
    "th": "tha_Thai",  # Thai
    "my": "mya_Mymr",  # Burmese/Myanmar
    "km": "khm_Khmr",  # Khmer
    "lo": "lao_Laoo",  # Lao
    "jv": "jav_Latn",  # Javanese
    "su": "sun_Latn",  # Sundanese
    "ceb": "ceb_Latn",  # Cebuano
    "mg": "plt_Latn",  # Malagasy
    # East Asian
    "zh-cn": "zho_Hans",  # Chinese Simplified
    "zh-tw": "zho_Hant",  # Chinese Traditional
    "ja": "jpn_Jpan",  # Japanese
    "ko": "kor_Hang",  # Korean
    "mn": "khk_Cyrl",  # Mongolian
    # African
    "sw": "swh_Latn",  # Swahili
    "yo": "yor_Latn",  # Yoruba
    "ig": "ibo_Latn",  # Igbo
    "ha": "hau_Latn",  # Hausa
    "am": "amh_Ethi",  # Amharic
    "so": "som_Latn",  # Somali
    "om": "gaz_Latn",  # Oromo
    "rw": "kin_Latn",  # Kinyarwanda
    "sn": "sna_Latn",  # Shona
    "st": "sot_Latn",  # Sotho
    "zu": "zul_Latn",  # Zulu
    "xh": "xho_Latn",  # Xhosa
    "tn": "tsn_Latn",  # Tswana
    "ln": "lin_Latn",  # Lingala
    "lg": "lug_Latn",  # Luganda
    "wo": "wol_Latn",  # Wolof
    "ff": "fuv_Latn",  # Fula
    "ak": "aka_Latn",  # Akan
    "ti": "tir_Ethi",  # Tigrinya
    "sg": "sag_Latn",  # Sango
    "lua": "lua_Latn",  # Luba-Kasai
    # Pacific
    "sm": "smo_Latn",  # Samoan
    "to": "ton_Latn",  # Tongan
    "mi": "mri_Latn",  # Maori
    "fj": "fij_Latn",  # Fijian
    # Constructed/Other
    "eo": "epo_Latn",  # Esperanto
}
