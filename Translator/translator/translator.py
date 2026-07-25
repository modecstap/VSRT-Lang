import argostranslate.package
import argostranslate.translate
from argostranslate.translate import Hypothesis
from argostranslate.translate import Language
from nltk.corpus import wordnet

from translator.i_translator import ITranslator, TranslatorError
from translator.models.tlanslation import Translation
from translator.models.word import Word


class TranslationPackageError(TranslatorError):
    """Translation package is unavailable."""


class TranslationExecutionError(TranslatorError):
    """Translation execution failed."""


class Translator(ITranslator):
    """English to Russian translator based on NLTK and Argos Translate."""

    def __init__(self) -> None:
        """Initialize translator."""
        self._translator = self._load_translator()
        self.translate("initial")

    def translate(self, phrase: str, count: int = 4) -> Translation:
        """
        Translate phrase from English to Russian.

        :param phrase: source phrase
        :param count: count of Translation

        :raises TranslationExecutionError:

        :return: translation result
        """
        translated = self._translate_text(phrase, count)
        return Translation(
            original=phrase,
            translations=translated,
        )

    def translate_bulk(
        self,
        phrase: list[str],
    ) -> list[Translation]:
        """
        Translate multiple phrases.

        :param phrase: source phrases

        :raises TranslationExecutionError:

        :return: translation results
        """
        return [self.translate(item) for item in phrase]

    def take_synonyms(self, target: Word) -> list[Word]:
        """
        Take English synonyms.

        :param target: source word

        :return: synonym list
        """
        return self._collect_related_words(
            target.content,
            antonyms=False,
        )

    def take_antonyms(self, target: Word) -> list[Word]:
        """
        Take English antonyms.

        :param target: source word

        :return: antonym list
        """
        return self._collect_related_words(
            target.content,
            antonyms=True,
        )

    def take_context(self, target: Word) -> list[Translation]:
        """
        Take usage examples translated to Russian.

        :param target: source word

        :raises TranslationExecutionError:

        :return: translated examples
        """
        examples: list[Translation] = []
        for example in self._collect_examples(target.content):
            examples.append(self.translate(example, 1))
        return examples

    def take_base_form(self, target: Word) -> Word:
        base_form = wordnet.morphy(target.content)
        return Word(content=base_form)

    def _load_translator(self):
        """Load installed Argos translator."""
        languages = argostranslate.translate.get_installed_languages()
        source = self._find_language(languages, "en")
        target = self._find_language(languages, "ru")
        if source is None or target is None:
            raise TranslationPackageError(
                "English to Russian package is not installed."
            )
        translator = source.get_translation(target)
        if translator is None:
            raise TranslationPackageError(
                "English to Russian translator is unavailable."
            )
        return translator

    def _translate_text(self, phrase: str, count: int) -> list[str]:
        """Translate text using Argos."""
        hypotheses = self.get_hypotheses(phrase, count)
        return list(map(lambda h: h.value, hypotheses))

    def get_hypotheses(self, phrase: str, count: int) -> list[Hypothesis]:
        try:
            return self._translator.hypotheses(phrase, count)
        except Exception as error:
            raise TranslationExecutionError(
                f"Failed to translate '{phrase}'."
            ) from error

    @staticmethod
    def _find_language(
            languages: list,
            code: str,
    ) -> Language:
        """Find installed language."""
        for language in languages:
            if language.code == code:
                return language
        return None

    @staticmethod
    def _collect_related_words(
            word: str,
        antonyms: bool,
    ) -> list[Word]:
        """Collect synonyms or antonyms."""
        values: set[str] = set()
        for synset in wordnet.synsets(word):
            for lemma in synset.lemmas():
                if antonyms:
                    values.update(
                        item.name().replace("_", " ")
                        for item in lemma.antonyms()
                    )
                else:
                    values.add(lemma.name().replace("_", " "))
        values.discard(word)
        return [Word(content=item) for item in sorted(values)]

    @staticmethod
    def _collect_examples(word: str, max_count: int = 3) -> list[str]:
        """Collect usage examples."""
        examples: set[str] = set()
        for synset in wordnet.synsets(word):
            if len(examples) >= max_count:
                break
            examples.update(synset.examples())
        return sorted(examples)
