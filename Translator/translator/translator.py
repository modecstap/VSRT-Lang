from urllib import error, request
import json

from nltk.corpus import wordnet

from translator.i_translator import ITranslator, TranslatorError
from translator.models.tlanslation import Translation
from translator.models.word import Word


class TranslationPackageError(TranslatorError):
    """Translation package is unavailable."""


class TranslationExecutionError(TranslatorError):
    """Translation execution failed."""


class Translator(ITranslator):
    """English to Russian translator based on NLTK and LibreTranslate API."""

    def __init__(self, url: str, key: str) -> None:
        """Initialize translator."""
        self._url = url
        self._api_key = key
        self.translate("initial")

    def translate(self, phrase: str, count: int = 4) -> Translation:
        """
        Translate phrase from English to Russian.

        :param phrase: source phrase
        :param count: count of translations

        :raises TranslationExecutionError:

        :return: translation result
        """
        translated = self._translate_text(phrase, count)
        return Translation(original=phrase, translations=translated)

    def translate_bulk(self, phrase: list[str]) -> list[Translation]:
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
        return self._collect_related_words(target.content, antonyms=False)

    def take_antonyms(self, target: Word) -> list[Word]:
        """
        Take English antonyms.

        :param target: source word

        :return: antonym list
        """
        return self._collect_related_words(target.content, antonyms=True)

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
        """
        Take the base form of an English word.

        :param target: source word

        :return: base form of the word
        """
        base_form = wordnet.morphy(target.content)
        return Word(content=base_form)

    def _translate_text(self, phrase: str, count: int) -> list[str]:
        """Translate text using the LibreTranslate API."""
        if count < 1:
            raise TranslationExecutionError(
                "Translation count must be greater than zero."
            )

        payload = self._build_payload(phrase, count)
        try:
            response = self._request_translation(payload)
        except (error.URLError, TimeoutError) as exception:
            raise TranslationExecutionError(
                f"Failed to translate '{phrase}'."
            ) from exception

        return self._extract_translations(response, phrase)

    def _build_payload(self, phrase: str, count: int) -> bytes:
        """Build a LibreTranslate request payload."""
        payload = {
            "q": phrase,
            "source": "en",
            "target": "ru",
            "format": "text",
            "alternatives": max(count - 1, 0),
        }
        if self._api_key:
            payload["api_key"] = self._api_key
        return json.dumps(payload).encode("utf-8")

    def _request_translation(self, payload: bytes) -> dict:
        """Send a translation request to LibreTranslate."""
        endpoint = f"{self._url}/translate"
        request_data = request.Request(
            endpoint,
            data=payload,
            headers={"Content-Type": "application/json"},
            method="POST",
        )
        try:
            with request.urlopen(request_data, timeout=30) as response:
                return json.loads(response.read().decode("utf-8"))
        except error.HTTPError as exception:
            raise TranslationExecutionError(
                f"LibreTranslate returned HTTP {exception.code}."
            ) from exception
        except json.JSONDecodeError as exception:
            raise TranslationExecutionError(
                "LibreTranslate returned an invalid response."
            ) from exception

    @staticmethod
    def _extract_translations(
        response: dict,
        phrase: str,
    ) -> list[str]:
        """Extract translations from a LibreTranslate response."""
        translated = response.get("translatedText")
        if not isinstance(translated, str):
            raise TranslationExecutionError(
                f"LibreTranslate returned no translation for '{phrase}'."
            )

        alternatives = response.get("alternatives", [])
        if not isinstance(alternatives, list):
            alternatives = []
        return [translated, *[
            item for item in alternatives if isinstance(item, str)
        ]]

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