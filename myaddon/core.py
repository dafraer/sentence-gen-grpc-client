import os
import re
import time
import tempfile

from .proto import sentence_gen_pb2

_MALE = "Male"
_AUDIO_FORMAT = "{}_{}.wav"


def _to_grpc_gender(gender):
    if gender == _MALE:
        return sentence_gen_pb2.GENDER_MALE
    return sentence_gen_pb2.GENDER_FEMALE


def _safe_filename(word):
    return re.sub(r"[^\w\-]", "_", word)


class Core:
    def __init__(self, grpc_client):
        self._client = grpc_client

    # ── Main-thread Anki calls ────────────────────────────────────────────

    def get_deck_names(self):
        from aqt import mw
        return sorted(d.name for d in mw.col.decks.all_names_and_ids())

    # ── Background-safe gRPC calls ────────────────────────────────────────

    def rpc_generate_sentence(self, word, word_language, translation_language,
                              translation_hint, include_audio, audio_gender):
        resp = self._client.generate_sentence(
            word=word,
            word_language=word_language,
            translation_language=translation_language,
            translation_hint=translation_hint,
            include_audio=include_audio,
            voice_gender=_to_grpc_gender(audio_gender),
        )
        audio_data = bytes(resp.audio.data) if include_audio and resp.audio.data else None
        return {
            "original_sentence": resp.original_sentence,
            "translated_sentence": resp.translated_sentence,
            "audio_data": audio_data,
        }

    def rpc_translate(self, word, from_language, to_language,
                      translation_hint, include_audio, audio_gender):
        resp = self._client.translate(
            word=word,
            from_language=from_language,
            to_language=to_language,
            translation_hint=translation_hint,
            include_audio=include_audio,
            voice_gender=_to_grpc_gender(audio_gender),
        )
        audio_data = bytes(resp.audio.data) if include_audio and resp.audio.data else None
        return {
            "translation": resp.translation,
            "audio_data": audio_data,
        }

    def rpc_generate_definition(self, word, language, definition_hint,
                                include_audio, audio_gender):
        resp = self._client.generate_definition(
            word=word,
            language=language,
            definition_hint=definition_hint,
            include_audio=include_audio,
            voice_gender=_to_grpc_gender(audio_gender),
        )
        audio_data = bytes(resp.audio.data) if include_audio and resp.audio.data else None
        return {
            "definition": resp.definition,
            "audio_data": audio_data,
        }

    # ── Main-thread Anki card operations ─────────────────────────────────

    def anki_add_sentence_card(self, word, deck_name, result):
        audio_fn = None
        if result["audio_data"]:
            audio_fn = _save_audio(word, result["audio_data"])
        _add_card(
            model_name="Basic (and reversed card)",
            deck_name=deck_name,
            front=result["original_sentence"],
            back=result["translated_sentence"],
            audio_filename=audio_fn,
        )

    def anki_add_translation_card(self, word, deck_name, result):
        audio_fn = None
        if result["audio_data"]:
            audio_fn = _save_audio(word, result["audio_data"])
        _add_card(
            model_name="Basic (and reversed card)",
            deck_name=deck_name,
            front=word,
            back=result["translation"],
            audio_filename=audio_fn,
        )

    def anki_add_definition_card(self, word, deck_name, result):
        audio_fn = None
        if result["audio_data"]:
            audio_fn = _save_audio(word, result["audio_data"])
        _add_card(
            model_name="Basic",
            deck_name=deck_name,
            front=word,
            back=result["definition"],
            audio_filename=audio_fn,
        )


def _save_audio(word, data):
    from aqt import mw
    filename = _AUDIO_FORMAT.format(_safe_filename(word), int(time.time()))
    with tempfile.NamedTemporaryFile(suffix=".wav", delete=False, prefix="sengen_") as f:
        f.write(data)
        temp_path = f.name
    try:
        stored_name = mw.col.media.add_file(temp_path)
    finally:
        try:
            os.unlink(temp_path)
        except OSError:
            pass
    return stored_name


def _add_card(model_name, deck_name, front, back, audio_filename=None):
    from aqt import mw
    from anki.notes import Note

    if audio_filename:
        front = front + f"[sound:{audio_filename}]"

    model = mw.col.models.by_name(model_name)
    if model is None:
        raise RuntimeError(f"Note type '{model_name}' not found. Please ensure it exists in Anki.")

    note = Note(mw.col, model)
    note["Front"] = front
    note["Back"] = back

    deck = mw.col.decks.by_name(deck_name)
    if deck is None:
        raise RuntimeError(f"Deck '{deck_name}' not found.")

    mw.col.add_note(note, deck["id"])
    mw.reset()
