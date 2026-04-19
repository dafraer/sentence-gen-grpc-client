<br />
<div align="center">

<h3 align="center">Sengen</h3>

  <p align="center">
    Sengen is an Anki add-on that automatically generates flashcards with example sentences, translations, and definitions. Just enter a word, pick your languages, and send it straight to your Anki deck.
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li>
      <a href="#tutorial">Tutorial</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Sengen connects to a backend service via gRPC to generate high-quality Anki cards on demand. Enter a word, choose the word language and translation language, and Sengen will generate an example sentence with translation, a standalone translation, or a monolingual definition — then add it directly to your chosen Anki deck. Optionally attach text-to-speech audio to any card.




<!-- TUTORIAL -->
## Tutorial

### Installation

1. Open Anki.
2. Go to **Tools → Add-ons → Get Add-ons…**
3. Enter the code **`289659486`** and click **OK**.
4. Restart Anki when prompted.
5. Open Anki and click **Tools → Sengen** to launch the add-on. On first use, required dependencies are downloaded automatically — Anki may appear frozen for ~30 seconds while this happens.

> **Note:** An internet connection is required both for the initial dependency download and for generating cards.



<!-- USAGE -->
## Usage

#### 1. Open Sengen

With Anki running, click **Tools → Sengen** from the menu bar.

#### 2. Choose a tab

Sengen has three card generation modes accessible from the tab bar at the top:

- **Generate Sentence** — generates an example sentence for a word, paired with its translation.
- **Translate** — generates a direct translation card for a word.
- **Generate Definition** — generates a monolingual definition card for a word.

#### 3. Fill in the form

Each tab has a short form:

| Field | Description |
|-------|-------------|
| Word | The word you want to study |
| Word Language | The language the word is in |
| Translation Language | The language to translate into (Sentence and Translate tabs) |
| Translation / Definition Hint | Optional hint to guide the AI output |
| Deck | Select the Anki deck to add the card to |
| Audio | Toggle to attach TTS audio to the card |
| Voice | Choose Male or Female TTS voice (enabled when Audio is checked) |

#### 4. Submit

Click **Generate**. Sengen will call the backend, generate the card content, and add it directly to your selected Anki deck. A system notification will confirm success or report an error.

#### 5. Review in Anki

Start a review session for your deck — your new card will be there, ready to study.



<!-- CONTACT -->
## Contact

Kamil Nuriev — [telegram](https://t.me/dafraer) — kdnuriev@gmail.com
