# and-plugin-konu_ac — New Topic Plugin

**English** · [Türkçe](README.tr.md)

Create new forum topics for [AND](https://github.com/AND-community/AND), save drafts and move
between categories. ("Konu aç" is Turkish for "open a topic".)

---

## Overview

`and-plugin-konu_ac` handles writing to the forum. AND's main forum view is read-only; every topic
is created through this plugin.

This plugin is **not shown** in the menu (its `Label` field is empty). When you press `n` in the
forum, AND starts it in the same terminal with the selected category (`AND_CATEGORY`); when the
plugin closes, the forum list is refreshed. If the plugin isn't installed, the `n` key does nothing
and no new topics can be created in the forum.

---

## Installation

Konu Aç lives in its own repository; it is cloned into the `Eklentiler/` folder of the
[AND](https://github.com/AND-community/AND) repository and built from AND's root directory:

```bash
git clone https://github.com/AND-community/konu_ac.git Eklentiler/konu_ac
go build -o Eklentiler/konu_ac/and-plugin-konu_ac ./Eklentiler/konu_ac

# Windows
go build -o Eklentiler/konu_ac/and-plugin-konu_ac.exe ./Eklentiler/konu_ac
```

When it starts, AND finds `Eklentiler/konu_ac/plugin.json` next to its own binary together with
this binary. The binary's name comes from the `name` field of plugin.json (`and-plugin-konu_ac`).
The repository also contains a prebuilt Windows binary, `and-plugin-konu_ac.exe`.

---

## How to open it

Open the forum and press `n`.  
AND passes the selected category automatically; the category selection screen is skipped.

---

## Usage

### Topic form

On the category row, `◀` / `▶` change the category.

| Key | Action |
|-----|--------|
| `tab` | Switch between Title ↔ Body |
| `enter` (Title) | Move to the Body field |
| `enter` (Body) | New line |
| `ctrl+k` (Body) | Insert a code block template at the cursor |
| `ctrl+s` | Send the topic |
| `ctrl+d` | Save as a draft (you stay in the form) |
| `ctrl+t` | Open the draft list (if there are drafts) |
| `◀` / `▶` | Change the category (only while in the Title field) |
| `esc` | Save the content as a draft and return to the forum |
| `ctrl+c` | Quit without saving |

### Draft list

| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Select a draft |
| `enter` or `e` | Load the selected draft into the form (to edit it) |
| `p` | Send the selected draft directly, without opening the form |
| `x` | Delete the selected draft |
| `esc` | Back to the form |

The plugin opens in AND's interface language (Turkish, English, German, Spanish, French or
Russian).

---

## Drafts

When you press `esc` and the title or body is not empty, a draft is saved automatically.  
When the same category is opened again, drafts can be brought back with `ctrl+t`.

Drafts are stored in `AND_DATA_DIR/taslaklar_<category>.json` files.  
These files are local, are never sent to the network and are listed in `.gitignore`.

---

## Code blocks

Text written between `[code]` and `[/code]` in the body is shown as a code box with line numbers
and syntax highlighting when the topic is viewed, as on other software forums. The tags are not
case-sensitive.

The `[code=language]` form sets the language (`python`, `go`, `js`/`javascript`,
`c`/`cpp`/`java`, `rust` and their common aliases are recognised). Without a language, a generic
highlighting based on the keywords of common languages is applied.

Pressing `ctrl+k` in the body inserts a ready-made `[code]` / `[/code]` template at the cursor and
leaves the cursor between the tags — no need to type them by hand. If there is an unclosed `[code]`
tag when you send (`ctrl+s`) or save a draft (`ctrl+d`), the action stops and a warning is shown.

```
There's a problem in this function:

[code=go]
func topla(a, b int) int {
    return a - b
}
[/code]

Why does it give the wrong result?
```

---

## Length limits

| Field | Maximum |
|------|---------|
| Title | 100 characters |
| Body | 8192 characters |

---

## Categories

Category values are the same on the network in every language (they travel in Turkish); AND and
this plugin show them in the interface language. The English names are:

| | | | |
|--|--|--|--|
| Python | C / C++ | Rust | Go |
| JavaScript | Java / Kotlin | C# / .NET | PHP |
| Swift / Objective-C | Software | Web | Mobile |
| Artificial Intelligence | Databases | DevOps | Cloud Computing |
| Linux | IT | Network & System Administration | Cybersecurity |
| Hardware | Embedded Systems / IoT | Robotics | Blockchain / Web3 |
| Algorithms & Data Structures | Testing & Quality | UI / UX Design | Game Development |
| Open Source | Career | Freelance | Learning Resources |
| Events & Announcements | General | | |

---

## A note on moderation

The forum has no topic approval: a topic spreads over the network as soon as it is sent and,
whoever wrote it, is deleted automatically **5 days** after it was created. The founder or a
moderator can delete a topic early from the Admin Panel when needed (see
[Eklentiler/admin/README.md](https://github.com/AND-community/AND/blob/main/Eklentiler/admin/README.md),
in Turkish).

---

## Manifest

```json
{
  "name":        "konu_ac",
  "label":       "",
  "version":     "2.6.1",
  "description": "Yeni konu oluşturma ve taslak yönetimi (menüde gizli — forumdan n ile açılır)",
  "author":      "AND",
  "permissions": ["forum"],
  "requires":    []
}
```

Because `label` is empty, this plugin doesn't appear in the main menu. It asks for the `forum`
permission the first time it is opened. (The description in the manifest is in Turkish; AND shows
it translated in the other interface languages.)

---

## Source

Source code: [main.go](main.go)

---

## Version history

| Version | Change |
|-------|------------|
| 2.6.1 | Drafts are written atomically: a crash or power cut while saving can't corrupt the draft file |
| 2.6.0 | Interface in AND's language (6 languages); categories come from AND's shared `protokol` package; asks for the `forum` permission on first launch |
| 2.5.0 | Topic approval removed; the `ctrl+p` permanence-request shortcut and its form field were deleted — every topic now lasts 5 days, whoever wrote it |
| 2.4.0 | Live preview removed (it caused screen flicker/leftovers); replaced by a one-key code block template with `ctrl+k` and a warning about unclosed `[code]` tags before sending; cursor movement with the left/right arrows in the Body field fixed |
| 2.3.0 | Language tag and badge with `[code=language]`; line numbers; basic syntax highlighting; live preview while typing `[code]` in the body |
| 2.2.0 | Content written with `[code]...[/code]` is shown as a code box when the topic is viewed; an info line was added to the form |
| 2.1.0 | Ran inline in AND's main process (no separate binary); Tab to switch fields; category ◀/▶; Enter for a new line in Body |
| 2.0.0 | Standalone binary; talks to AND over HTTP IPC; draft system |
| 1.0.0 | First version |
