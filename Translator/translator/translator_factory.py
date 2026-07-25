from translator.i_translator import ITranslator
from translator.translator import Translator


class TranslatorFactory:
    @staticmethod
    def new_translator() -> ITranslator:
        return Translator()
