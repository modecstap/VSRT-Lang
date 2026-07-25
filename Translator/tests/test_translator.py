import pytest

from translator.translator import Translator


@pytest.fixture(scope="session")
def translator():
    return Translator()


@pytest.mark.parametrize(
    "original",[
        "simple sentence",
        "difficult sentence"
    ]
)
def test_translation(translator: Translator, original: str):
    translation = translator.translate(original)

    assert translation.original == original
    assert isinstance(translation.translations, list)
    assert len(translation.translations) > 0

@pytest.mark.parametrize(
    "originals",[
        ("first sentence", "second sentence"),
        ("difficult sentence", "simple sentence")
    ]
)
def test_translation_bulk(translator: Translator, originals: list[str]):
    translations = translator.translate_bulk(originals)

    for i, translation in enumerate(translations):
        assert translation.original == originals[i]
        assert isinstance(translation.translations, list)
        assert len(translation.translations) > 0