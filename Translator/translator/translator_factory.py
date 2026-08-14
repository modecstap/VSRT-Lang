import os

from translator.i_translator import ITranslator
from translator.translator import Translator


class TranslatorFactory:
    @staticmethod
    def new_translator() -> ITranslator:
        url = os.getenv(
            "LIBRETRANSLATE_URL",
            "localhost:5000",
        )
        url = "http://"+url
        key = os.getenv("LIBRETRANSLATE_API_KEY", "")
        return Translator(url, key)
