import threading

from aqt.qt import *
from aqt import mw
from aqt.utils import showWarning, tooltip

# PyQt5 / PyQt6 enum compatibility
try:
    _FORM_EXPAND = QFormLayout.FieldGrowthPolicy.ExpandingFieldsGrow
    _SIZE_EXPAND = QSizePolicy.Policy.Expanding
    _SIZE_FIXED = QSizePolicy.Policy.Fixed
except AttributeError:
    _FORM_EXPAND = QFormLayout.ExpandingFieldsGrow  # type: ignore[attr-defined]
    _SIZE_EXPAND = QSizePolicy.Expanding  # type: ignore[attr-defined]
    _SIZE_FIXED = QSizePolicy.Fixed  # type: ignore[attr-defined]

from . import languages as lang_module
from .grpc_client import (
    GrpcError,
    InvalidArgumentError,
    DeadlineExceededError,
    ResourceExhaustedError,
    InternalServerError,
    UnavailableError,
    UnknownError,
)

MAX_WORD_LEN = 100
MAX_HINT_LEN = 100
PROJECT_URL = "https://github.com/dafraer/sentence-gen-grpc-client"


class SengenDialog(QDialog):
    def __init__(self, core, parent=None):
        super().__init__(parent)
        self._core = core
        self.setWindowTitle("Sengen")
        self.resize(720, 480)

        try:
            deck_names = core.get_deck_names()
        except Exception:
            deck_names = []

        layout = QVBoxLayout(self)
        layout.setContentsMargins(0, 0, 0, 0)

        tabs = QTabWidget()
        layout.addWidget(tabs)

        tabs.addTab(
            self._make_translation_form_tab(
                title="Generate sentences",
                deck_names=deck_names,
                has_translation_lang=True,
                on_rpc=self._rpc_generate_sentence,
                on_anki=self._anki_add_sentence,
                success_msg=lambda word: f"Sentence with the word '{word}' has been added successfully",
            ),
            "Generate sentence",
        )
        tabs.addTab(
            self._make_translation_form_tab(
                title="Generate translation",
                deck_names=deck_names,
                has_translation_lang=True,
                on_rpc=self._rpc_translate,
                on_anki=self._anki_add_translation,
                success_msg=lambda word: f"Translation with the word '{word}' has been added successfully",
            ),
            "Translate",
        )
        tabs.addTab(
            self._make_definition_tab(deck_names),
            "Generate definition",
        )
        tabs.addTab(self._make_tutorial_tab(), "Tutorial")

    # ── Per-tab RPC/Anki delegates ────────────────────────────────────────

    def _rpc_generate_sentence(self, word, word_lang, trans_lang, hint, include_audio, gender):
        return self._core.rpc_generate_sentence(
            word=word,
            word_language=word_lang,
            translation_language=trans_lang,
            translation_hint=hint,
            include_audio=include_audio,
            audio_gender=gender,
        )

    def _anki_add_sentence(self, word, deck, result):
        self._core.anki_add_sentence_card(word, deck, result)

    def _rpc_translate(self, word, word_lang, trans_lang, hint, include_audio, gender):
        return self._core.rpc_translate(
            word=word,
            from_language=word_lang,
            to_language=trans_lang,
            translation_hint=hint,
            include_audio=include_audio,
            audio_gender=gender,
        )

    def _anki_add_translation(self, word, deck, result):
        self._core.anki_add_translation_card(word, deck, result)

    # ── Tab builders ──────────────────────────────────────────────────────

    def _make_translation_form_tab(self, title, deck_names, has_translation_lang,
                                   on_rpc, on_anki, success_msg):
        widget = QWidget()
        vbox = QVBoxLayout(widget)
        vbox.setContentsMargins(12, 12, 12, 12)

        header = QLabel(f"<h2>{title}</h2>")
        vbox.addWidget(header)

        form_layout = QFormLayout()
        form_layout.setFieldGrowthPolicy(_FORM_EXPAND)

        word_edit = QLineEdit()
        hint_edit = QLineEdit()

        word_lang_combo = _make_lang_combo()
        trans_lang_combo = _make_lang_combo() if has_translation_lang else None

        audio_check = QCheckBox()
        voice_combo = _make_voice_combo()

        deck_combo = _make_deck_combo(deck_names)

        form_layout.addRow("Word:", word_edit)
        form_layout.addRow("Translation Hint:", hint_edit)
        form_layout.addRow("Pick word language:", word_lang_combo)
        if trans_lang_combo is not None:
            form_layout.addRow("Pick translation language:", trans_lang_combo)
        form_layout.addRow("Audio:", audio_check)
        form_layout.addRow("Voice:", voice_combo)
        form_layout.addRow("Deck:", deck_combo)

        submit_btn = QPushButton("Generate")
        submit_btn.setEnabled(False)

        btn_row = QHBoxLayout()
        btn_row.addStretch()
        btn_row.addWidget(submit_btn)

        vbox.addLayout(form_layout)
        vbox.addLayout(btn_row)
        vbox.addStretch()

        audio_check.stateChanged.connect(
            lambda state: voice_combo.setEnabled(bool(state))
        )

        def validate(_=None):
            word = word_edit.text()
            hint = hint_edit.text()
            lang1 = _combo_value(word_lang_combo)
            lang2 = _combo_value(trans_lang_combo) if trans_lang_combo is not None else "ok"
            deck = _combo_value(deck_combo)
            if trans_lang_combo is not None:
                err = self._validate_form(word, hint, lang1, lang2, deck)
            else:
                err = self._validate_definition_form(word, hint, lang1, deck)
            submit_btn.setEnabled(err is None)

        word_edit.textChanged.connect(validate)
        hint_edit.textChanged.connect(validate)
        word_lang_combo.currentIndexChanged.connect(validate)
        if trans_lang_combo is not None:
            trans_lang_combo.currentIndexChanged.connect(validate)
        deck_combo.currentIndexChanged.connect(validate)

        def on_submit():
            word = word_edit.text()
            hint = hint_edit.text()
            lang1 = _combo_value(word_lang_combo)
            lang2 = _combo_value(trans_lang_combo) if trans_lang_combo is not None else ""
            deck = _combo_value(deck_combo)

            if trans_lang_combo is not None:
                err = self._validate_form(word, hint, lang1, lang2, deck)
            else:
                err = self._validate_definition_form(word, hint, lang1, deck)
            if err:
                showWarning(err)
                submit_btn.setEnabled(False)
                return

            word_lang_code = lang_module.get_language_code(lang1)
            if not word_lang_code:
                showWarning("Unknown language")
                return

            trans_lang_code = None
            if trans_lang_combo is not None:
                trans_lang_code = lang_module.get_language_code(lang2)
                if not trans_lang_code:
                    showWarning("Unknown language")
                    return

            include_audio = audio_check.isChecked()
            gender = voice_combo.currentText()

            # Clear fields and disable button before starting background task
            word_edit.clear()
            hint_edit.clear()
            submit_btn.setEnabled(False)

            # Snapshot the word for use in callbacks (word_edit is already cleared)
            captured_word = word

            def bg():
                try:
                    result = on_rpc(
                        captured_word,
                        word_lang_code,
                        trans_lang_code,
                        hint,
                        include_audio,
                        gender,
                    )

                    def main_ok():
                        try:
                            on_anki(captured_word, deck, result)
                            tooltip(success_msg(captured_word))
                        except Exception as e:
                            showWarning(f"Error adding word '{captured_word}': {e}")
                        finally:
                            validate()

                    mw.taskman.run_on_main(main_ok)

                except Exception as e:
                    err_msg = _map_error_msg(e, captured_word)

                    def main_err():
                        showWarning(err_msg)
                        validate()

                    mw.taskman.run_on_main(main_err)

            threading.Thread(target=bg, daemon=True).start()

        submit_btn.clicked.connect(on_submit)
        word_edit.returnPressed.connect(lambda: submit_btn.click() if submit_btn.isEnabled() else None)
        hint_edit.returnPressed.connect(lambda: submit_btn.click() if submit_btn.isEnabled() else None)

        return widget

    def _make_definition_tab(self, deck_names):
        widget = QWidget()
        vbox = QVBoxLayout(widget)
        vbox.setContentsMargins(12, 12, 12, 12)

        header = QLabel("<h2>Generate definition</h2>")
        vbox.addWidget(header)

        form_layout = QFormLayout()
        form_layout.setFieldGrowthPolicy(_FORM_EXPAND)

        word_edit = QLineEdit()
        hint_edit = QLineEdit()
        word_lang_combo = _make_lang_combo()
        audio_check = QCheckBox()
        voice_combo = _make_voice_combo()
        deck_combo = _make_deck_combo(deck_names)

        form_layout.addRow("Word:", word_edit)
        form_layout.addRow("Definition Hint:", hint_edit)
        form_layout.addRow("Pick word language:", word_lang_combo)
        form_layout.addRow("Audio:", audio_check)
        form_layout.addRow("Voice:", voice_combo)
        form_layout.addRow("Deck:", deck_combo)

        submit_btn = QPushButton("Generate")
        submit_btn.setEnabled(False)

        btn_row = QHBoxLayout()
        btn_row.addStretch()
        btn_row.addWidget(submit_btn)

        vbox.addLayout(form_layout)
        vbox.addLayout(btn_row)
        vbox.addStretch()

        audio_check.stateChanged.connect(
            lambda state: voice_combo.setEnabled(bool(state))
        )

        def validate(_=None):
            word = word_edit.text()
            hint = hint_edit.text()
            lang = _combo_value(word_lang_combo)
            deck = _combo_value(deck_combo)
            err = self._validate_definition_form(word, hint, lang, deck)
            submit_btn.setEnabled(err is None)

        word_edit.textChanged.connect(validate)
        hint_edit.textChanged.connect(validate)
        word_lang_combo.currentIndexChanged.connect(validate)
        deck_combo.currentIndexChanged.connect(validate)

        def on_submit():
            word = word_edit.text()
            hint = hint_edit.text()
            lang = _combo_value(word_lang_combo)
            deck = _combo_value(deck_combo)

            err = self._validate_definition_form(word, hint, lang, deck)
            if err:
                showWarning(err)
                submit_btn.setEnabled(False)
                return

            lang_code = lang_module.get_language_code(lang)
            if not lang_code:
                showWarning("Unknown language")
                return

            include_audio = audio_check.isChecked()
            gender = voice_combo.currentText()

            word_edit.clear()
            hint_edit.clear()
            submit_btn.setEnabled(False)

            captured_word = word

            def bg():
                try:
                    result = self._core.rpc_generate_definition(
                        word=captured_word,
                        language=lang_code,
                        definition_hint=hint,
                        include_audio=include_audio,
                        audio_gender=gender,
                    )

                    def main_ok():
                        try:
                            self._core.anki_add_definition_card(captured_word, deck, result)
                            tooltip(f"Definition of the word '{captured_word}' has been added successfully")
                        except Exception as e:
                            showWarning(f"Error adding word '{captured_word}': {e}")
                        finally:
                            validate()

                    mw.taskman.run_on_main(main_ok)

                except Exception as e:
                    err_msg = _map_error_msg(e, captured_word)

                    def main_err():
                        showWarning(err_msg)
                        validate()

                    mw.taskman.run_on_main(main_err)

            threading.Thread(target=bg, daemon=True).start()

        submit_btn.clicked.connect(on_submit)
        word_edit.returnPressed.connect(lambda: submit_btn.click() if submit_btn.isEnabled() else None)
        hint_edit.returnPressed.connect(lambda: submit_btn.click() if submit_btn.isEnabled() else None)

        return widget

    def _make_tutorial_tab(self):
        widget = QWidget()
        vbox = QVBoxLayout(widget)
        vbox.setContentsMargins(12, 12, 12, 12)

        header = QLabel("<h2>Tutorial</h2>")
        vbox.addWidget(header)

        body = QLabel(
            f"You can check out the tutorial on the project's "
            f'<a href="{PROJECT_URL}">github page</a>'
        )
        body.setOpenExternalLinks(True)
        body.setWordWrap(True)
        vbox.addWidget(body)

        vbox.addStretch()
        return widget

    # ── Validation (mirrors validators.go) ───────────────────────────────

    def _validate_word(self, w):
        if not w.strip():
            return "word required"
        if len(w) > MAX_WORD_LEN:
            return "word too long"
        return None

    def _validate_languages(self, lang1, lang2):
        if not lang1 or not lang2:
            return "pick languages"
        if lang1 == lang2:
            return "languages cannot be the same"
        return None

    def _validate_hint(self, h):
        if len(h) > MAX_HINT_LEN:
            return "hint too long"
        return None

    def _validate_form(self, word, hint, lang1, lang2, deck):
        err = self._validate_word(word)
        if err:
            return err
        err = self._validate_languages(lang1, lang2)
        if err:
            return err
        if not deck:
            return "pick a deck"
        return self._validate_hint(hint)

    def _validate_definition_form(self, word, hint, lang, deck):
        err = self._validate_word(word)
        if err:
            return err
        if not lang:
            return "pick a language"
        if not deck:
            return "pick a deck"
        return self._validate_hint(hint)


# ── Widget helpers ────────────────────────────────────────────────────────────

def _make_lang_combo():
    combo = QComboBox()
    combo.setSizePolicy(_SIZE_EXPAND, _SIZE_FIXED)
    combo.addItem("")  # blank placeholder
    combo.addItems(lang_module.LANGUAGE_NAMES)
    combo.setCurrentIndex(0)
    return combo


def _make_deck_combo(deck_names):
    combo = QComboBox()
    combo.setSizePolicy(_SIZE_EXPAND, _SIZE_FIXED)
    combo.addItem("")  # blank placeholder
    combo.addItems(deck_names)
    combo.setCurrentIndex(0)
    return combo


def _make_voice_combo():
    combo = QComboBox()
    combo.addItems(["Male", "Female"])
    combo.setCurrentText("Female")
    combo.setEnabled(False)
    return combo


def _combo_value(combo):
    if combo is None:
        return ""
    text = combo.currentText()
    return text  # "" if blank placeholder selected


def _map_error_msg(e, word):
    if isinstance(e, InvalidArgumentError):
        err_text = "Invalid argument"
    elif isinstance(e, DeadlineExceededError):
        err_text = "Server deadline exceeded"
    elif isinstance(e, InternalServerError):
        err_text = "Internal server error"
    elif isinstance(e, ResourceExhaustedError):
        err_text = "Quota limit reached, try again tomorrow"
    elif isinstance(e, UnavailableError):
        err_text = "Server is unavailable"
    elif isinstance(e, GrpcError):
        err_text = "Unknown error"
    else:
        err_text = str(e)
    return f"Error adding word '{word}': {err_text}"
