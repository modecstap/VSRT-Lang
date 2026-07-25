from __future__ import annotations

from multiprocessing import Process, Queue
from threading import Condition
from typing import Any

from translator.i_translator import ITranslator
from translator.models.tlanslation import Translation
from translator.models.word import Word
from translator.translator_factory import TranslatorFactory


class ProcessPoolError(Exception):
    """Base process pool error."""


class ProcessPoolClosedError(ProcessPoolError):
    """Pool is closed."""


class ProcessPoolProxyReleasedError(ProcessPoolError):
    """Proxy has already been released."""


class ProcessPoolProxy:
    """
    Translator proxy bound to a dedicated worker process.
    """

    def __init__(
        self,
        pool: "ProcessPool",
        worker_id: int,
        request_queue: Queue,
        response_queue: Queue,
    ):
        """
        Create process proxy.

        :param pool: owner pool
        :param worker_id: worker identifier
        :param request_queue: worker request queue
        :param response_queue: worker response queue
        """
        self._pool = pool
        self._worker_id = worker_id
        self._request_queue = request_queue
        self._response_queue = response_queue
        self._released = False

    def __enter__(self) -> "ProcessPoolProxy":
        """
        Enter proxy context.

        :return: current proxy
        """
        return self

    def __exit__(self, *_: Any) -> None:
        """
        Leave proxy context.
        """
        self.release()

    def translate(self, phrase: str) -> Translation:
        """
        Translate phrase.

        :param phrase: source phrase

        :return: translation
        """
        return self._call("translate", phrase)

    def translate_bulk(self, phrase: list[str]) -> list[Translation]:
        """
        Translate phrases.

        :param phrase: source phrases

        :return: translations
        """
        return self._call("translate_bulk", phrase)

    def take_synonyms(self, target: Word) -> list[Word]:
        """
        Take synonyms.

        :param target: target word

        :return: synonyms
        """
        return self._call("take_synonyms", target)

    def take_antonyms(self, target: Word) -> list[Word]:
        """
        Take antonyms.

        :param target: target word

        :return: antonyms
        """
        return self._call("take_antonyms", target)

    def take_context(self, target: Word) -> list[Translation]:
        """
        Take translation context.

        :param target: target word

        :return: context
        """
        return self._call("take_context", target)

    def take_base_form(self, target: Word) -> Word:
        """
        Take base form.

        :param target: target word

        :return: base form
        """
        return self._call("take_base_form", target)

    def release(self) -> None:
        """
        Release occupied worker.

        :raises ProcessPoolProxyReleasedError:
        """
        if self._released:
            raise ProcessPoolProxyReleasedError()

        self._released = True
        self._pool._release(self._worker_id)

    def _call(self, method: str, argument: Any) -> Any:
        """
        Execute remote method.

        :param method: translator method
        :param argument: method argument

        :raises ProcessPoolProxyReleasedError:

        :return: method result
        """
        if self._released:
            raise ProcessPoolProxyReleasedError()

        self._request_queue.put((method, argument))
        success, payload = self._response_queue.get()

        if success:
            return payload

        raise payload


class ProcessPool:
    """
    Pool of dedicated translator processes.
    """

    def __init__(
        self,
        factory: TranslatorFactory,
        pool_size: int,
    ):
        """
        Create process pool.

        :param factory: translator factory
        :param pool_size: worker count
        """
        self._closed = False
        self._condition = Condition()
        self._busy = [False] * pool_size
        self._workers: list[tuple[Queue, Queue, Process]] = []

        for _ in range(pool_size):
            request_queue = Queue()
            response_queue = Queue()
            process = Process(
                target=self._worker,
                args=(
                    factory,
                    request_queue,
                    response_queue,
                ),
                daemon=True,
            )
            process.start()
            self._workers.append(
                (
                    request_queue,
                    response_queue,
                    process,
                )
            )

    def __enter__(self) -> "ProcessPool":
        """
        Enter pool context.

        :return: current pool
        """
        return self

    def __exit__(self, *_: Any) -> None:
        """
        Leave pool context.
        """
        self.close()

    def get(self) -> ProcessPoolProxy:
        """
        Acquire free worker.

        :raises ProcessPoolClosedError:

        :return: process proxy
        """
        with self._condition:
            while True:
                if self._closed:
                    raise ProcessPoolClosedError()

                worker_id = self._take_free_worker()

                if worker_id >= 0:
                    request, response, _ = self._workers[worker_id]
                    return ProcessPoolProxy(
                        self,
                        worker_id,
                        request,
                        response,
                    )

                self._condition.wait()

    def close(self) -> None:
        """
        Stop all worker processes.
        """
        with self._condition:
            if self._closed:
                return

            self._closed = True
            self._condition.notify_all()

        for request_queue, _, _ in self._workers:
            request_queue.put(None)

        for _, _, process in self._workers:
            process.join()

    def _release(self, worker_id: int) -> None:
        """
        Release occupied worker.

        :param worker_id: worker identifier
        """
        with self._condition:
            self._busy[worker_id] = False
            self._condition.notify()

    def _take_free_worker(self) -> int:
        """
        Reserve free worker.

        :return: worker identifier
        """
        for index, busy in enumerate(self._busy):
            if not busy:
                self._busy[index] = True
                return index

        return -1

    @staticmethod
    def _worker(
        factory: TranslatorFactory,
        request_queue: Queue,
        response_queue: Queue,
    ) -> None:
        """
        Execute worker loop.

        :param factory: translator factory
        :param request_queue: incoming requests
        :param response_queue: outgoing responses
        """
        translator: ITranslator = factory.new_translator()

        while True:
            command = request_queue.get()

            if command is None:
                return

            method, argument = command

            try:
                result = getattr(translator, method)(argument)
                response_queue.put((True, result))
            except Exception as error:
                response_queue.put((False, error))