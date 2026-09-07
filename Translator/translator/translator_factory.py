import os
from time import sleep

from translator import translator
from translator.i_translator import ITranslator
from translator.translator import Translator


TIME_BEFORE_RECONNECT = 50 # in seconds


class TranslatorFactory:
    @staticmethod
    def new_translator() -> ITranslator:
        url = os.getenv(
            "LIBRETRANSLATE_URL",
            "localhost:5000",
        )
        url = "http://"+url+"/"
        print(url)
        key = os.getenv("LIBRETRANSLATE_API_KEY", "")

        translator = None

        while not translator:
            try:
                translator = Translator(url, key)
            except Exception as e:
                print(e)
                sleep(TIME_BEFORE_RECONNECT)

        return translator
