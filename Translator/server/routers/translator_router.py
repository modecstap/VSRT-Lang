from fastapi import APIRouter

from process_pool import ProcessPool
from server.handlers.translation_handler import TranslationHandler
from translator.translator_factory import TranslatorFactory


class TranslatorRouter:
    def __init__(self, prefix: str, tags: list[str]):

        handler = TranslationHandler()
        self.router = APIRouter(
            prefix=prefix,
            tags=tags
        )

        self.router.post("/translate")(handler.translate)
        self.router.post("/translate-bulk")(handler.translate_bulk)
        self.router.get("/synonyms")(handler.take_synonyms)
        self.router.get("/antonyms")(handler.take_antonyms)
        self.router.get("/base-form")(handler.take_base_form)
        self.router.get("/context")(handler.take_context)
