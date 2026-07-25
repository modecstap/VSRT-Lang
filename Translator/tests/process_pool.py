import threading
import time

import pytest

from process_pool import ProcessPool, ProcessPoolReleasedError


class Worker:
    def __init__(self) -> None:
        self.value = 0

    def increment(self, delta: int) -> int:
        self.value += delta
        return self.value

    def wait(self, delay: float) -> int:
        time.sleep(delay)
        return self.value


def test_should_forward_method_call():
    pool = ProcessPool([Worker()])

    with pool.get() as worker:
        assert worker.increment(2) == 2
        assert worker.increment(3) == 5


def test_should_raise_after_release():
    pool = ProcessPool([Worker()])

    worker = pool.get()
    worker.release()

    with pytest.raises(ProcessPoolReleasedError):
        worker.increment(1)


def test_should_release_on_context_exit():
    pool = ProcessPool([Worker()])

    with pool.get() as worker:
        worker.increment(1)

    with pool.get() as worker:
        assert worker.increment(1) == 2


def test_should_wait_until_worker_becomes_available():
    pool = ProcessPool([Worker()])
    result = []

    def first_thread():
        with pool.get() as worker:
            worker.wait(0.2)
            result.append("first")

    def second_thread():
        with pool.get() as worker:
            worker.increment(1)
            result.append("second")

    thread1 = threading.Thread(target=first_thread)
    thread2 = threading.Thread(target=second_thread)

    thread1.start()
    time.sleep(0.05)
    thread2.start()

    thread1.join()
    thread2.join()

    assert result == ["first", "second"]


def test_should_use_different_workers():
    pool = ProcessPool([Worker(), Worker()])

    with pool.get() as first:
        with pool.get() as second:
            assert first.increment(1) == 1
            assert second.increment(1) == 1


def test_should_not_allow_double_release():
    pool = ProcessPool([Worker()])

    worker = pool.get()
    worker.release()

    with pytest.raises(ProcessPoolReleasedError):
        worker.release()