import os

from process_pool import ProcessPool
from translator.models.tlanslation import Translation
from translator.models.word import Word
from translator.translator_factory import TranslatorFactory


class TranslationHandler:
    def __init__(self):
        self._pool = ProcessPool(
            TranslatorFactory(),
            os.getenv("TRANSLATOR_COUNT", 1)
        )

    def translate(self, phrase: str) -> Translation:
        with self._pool.get() as translator:
            response = translator.translate(phrase)

        return response

    def translate_bulk(self, phrases: list[str]) -> list[Translation]:
        with self._pool.get() as translator:
            response = translator.translate_bulk(phrases)

        return response

    def take_synonyms(self, target: str) -> list[str]:
        with self._pool.get() as translator:

            words = translator.take_synonyms(Word(content=target))

        return [word.content for word in words]

    def take_antonyms(self, target: str) -> list[str]:
        with self._pool.get() as translator:
            words = translator.take_antonyms(Word(content=target))

        return [word.content for word in words]

    def take_context(self, target: str) -> list[Translation]:
        with self._pool.get() as translator:
            response = translator.take_context(Word(content=target))

        return response

    def take_base_form(self, target: str) -> str:
        with self._pool.get() as translator:
            word = translator.take_base_form(Word(content=target))

        return word.content
