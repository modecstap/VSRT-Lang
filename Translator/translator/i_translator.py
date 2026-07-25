from abc import ABC, abstractmethod

from translator.models.word import Word
from translator.models.tlanslation import Translation


class TranslatorError(Exception):
    """Base translator error."""


class ITranslator(ABC):

    @abstractmethod
    def translate(self, phrase: str) -> Translation:
        pass

    @abstractmethod
    def translate_bulk(self, phrase: list[str]) -> list[Translation]:
        pass

    @abstractmethod
    def take_synonyms(self, target: Word) -> list[Word]:
        pass

    @abstractmethod
    def take_antonyms(self, target: Word) -> list[Word]:
        pass

    @abstractmethod
    def take_context(self, target: Word) -> list[Translation]:
        pass

    @abstractmethod
    def take_base_form(self, target: Word) -> Word:
        pass
