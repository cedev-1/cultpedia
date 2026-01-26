# Contributing to Cultpedia

Thank you for your interest in contributing to Cultpedia! Your contributions help improve the quality and diversity of the content.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Contributing Questions](#contributing-questions)
- [Contributing Code](#contributing-code)

## Code of Conduct

By participating, you agree to maintain a respectful and inclusive environment. Harassment or discriminatory behavior will not be tolerated.

## Getting Started

You have to **Fork the repository** on GitHub and clone your fork into your local machine.

```bash
git clone https://github.com/YOUR_USERNAME/cultpedia.git
cd cultpedia
```

Replace `YOUR_USERNAME` with your GitHub username.

### Build the cultpedia tool

You have two ways to get the **cultpedia tool**, you can get the pre-built binary or build it from source.

- [Use Release](#use-release)
- [Manual build](#manual-build)

### Use Release

If you want to use the pre-built binary releases, go to the [Releases page](https://github.com/Culturae-org/cultpedia/releases). Download the ZIP file for your platform, extract it, and run the `cultpedia` binary (or `cultpedia.exe` on Windows).

> [!WARNING]
> You have to move the binary into the cloned repository folder to use it with datasets (see this step [Fork](#getting-started)).

### Manual Build

**Linux/macOS:**
```bash
./build.sh
```

**Windows:**
```bash
build.bat
```

**Manual build:**
```bash
go build -o cultpedia ./cmd
```

> [!TIP]
> If you use **Nix**, you can set up the development environment (binary will be built automatically)

**Familiarize yourself** with the [data format](FORMAT.md).

## Contributing Questions

> [!IMPORTANT]
> Do not edit the questions.ndjson file directly. Use the interactive TUI tool described below!

### Available Themes

You can use one of these existing themes for your question:

| Theme | Description |
|-------|-------------|
| `science` | Physics, chemistry, biology, astronomy, inventions |
| `history` | Ancient history, modern history, wars, emperors, revolutions |
| `geography` | Countries, capitals, landmarks, continents |
| `sports` | Football, motorsport, Formula 1, Olympics |
| `gaming` | Video games, studios, consoles |

You can also create new **theme**, new **subthemes** and **tags** as needed. They will be automatically added to the dataset when your PR is merged.

### Question Content Guidelines

#### i18n Fields (title, stem, explanation)

Each question requires translations in **French (fr)**, **English (en)**, and **Spanish (es)**:

| Field | Description | Example |
|-------|-------------|---------|
| `title` | Short descriptive title | "First woman Nobel Prize" |
| `stem` | The actual question (min. 10 characters) | "Who was the first woman to win a Nobel Prize in 1903?" |
| `explanation` | Educational explanation of the answer (min. 20 characters) | "Marie Curie won the Nobel Prize in Physics in 1903, shared with her husband Pierre Curie and Henri Becquerel, for their research on radioactivity." |

#### Difficulty Levels

- beginner
- intermediate
- advanced
- pro

#### Validation Rules

Your question must respect these constraints:

| Field | Constraint |
|-------|------------|
| `slug` | Lowercase letters, numbers, and hyphens only. No leading/trailing hyphens. Must be unique. |
| `points` | Between 0.5 and 5.0 |
| `estimated_seconds` | Between 5 and 30 |
| `answers` | Exactly 4 for `single_choice`, exactly 2 for `true_false` |
| `sources` | At least one URL required |

#### Sources

Provide at least one reliable source URL (Wikipedia, official websites, academic sources, reputable news outlets). This allows reviewers to verify the accuracy of your question.

### Step-by-Step Guide

1. **Prerequisites**: Make sure you have completed Fork, Clone, and Build steps.

2. **Choose your question type and edit the template file:**
   - **Standard questions (4 choices):** Edit [`datasets/new-question.json`](../datasets/new-question.json)
   - **True/False questions (2 choices):** Edit [`datasets/new-question-true-false.json`](../datasets/new-question-true-false.json)

3. **Validate locally** (optional but recommended):
   ```bash
   ./cultpedia validate
   ```

4. **Run the TUI tool** to add your question:
   ```bash
   ./cultpedia
   ```
   Follow **Validate new question** and **Add question to dataset** steps.

   ![Interactive Tool](./cultpedia.gif)

   ![add_question](./add_question.png)

> [!WARNING]
> Don't edit **.ndjson** and **manifest.json** files directly.

5. **Create a branch and push to your fork:**

   ```bash
   # Create a new branch
   git checkout -b add-question-{slug}

   # Add only the questions file
   git add datasets/general-knowledge/questions.ndjson

   # Commit your changes
   git commit -m "feat: add {slug}"

   # Push to your fork
   git push origin add-question-{slug}
   ```

   Replace `{slug}` with your question's slug (e.g., `add-question-science-physics-marie-curie-nobel`).

6. **Create a Pull Request** on GitHub from your branch to the main repository.

7. **CI will automatically**:
   - ✓ Validate all questions
   - ✓ Check for duplicates
   - ✓ Verify all translations
   - ✓ Reject if unwanted files were modified

8. Reviewers will check your PR, may request changes, and finally merge it.

### Troubleshooting

| Error | Solution |
|-------|----------|
| "slug must be lowercase with hyphens only" | Use only `a-z`, `0-9`, and `-`. No spaces, underscores, or uppercase. |
| "must have exactly 4 answers" | Add or remove answers to have exactly 4 (or 2 for true/false). |
| "missing X translation" | Add the missing language (fr, en, or es) to i18n fields. |
| "stem too short" | Write a more detailed question (minimum 10 characters). |
| "explanation too short" | Provide a more detailed explanation (minimum 20 characters). |
| "at least one source URL is required" | Add a source URL to verify your question's accuracy. |

### Important Rules

- **Only modify** `datasets/general-knowledge/questions.ndjson`
- **Do not commit** manifest.json, themes.ndjson, subthemes.ndjson, or tags.ndjson
- These files are updated automatically by the CI

## Contributing Code

### Setup

1. **Fork and clone** the repository (see [Getting Started](#getting-started))

2. **Install Go** (version 1.24 or later)

3. **Build the project:**
   ```bash
   go build -o cultpedia ./cmd
   ```

4. **Run tests:**
   ```bash
   go test ./... -v
   ```

5. **Run linter:**
   ```bash
   golangci-lint run --timeout=5m
   ```

### Project Structure

```
cmd/main.go           # CLI entry point
internal/
├── actions/          # Domain logic (question management, API server)
├── checks/           # Validation pipeline
├── models/           # Data structures
├── ui/               # Bubble Tea TUI components
└── utils/            # File I/O helpers
```

### Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Ensure tests pass: `go test ./... -v`
4. Ensure linter passes: `golangci-lint run`
5. Submit a PR with a clear description of changes

## Thank You!

We appreciate your contributions to Cultpedia! Your efforts help create a richer educational resource for everyone.

## Additional Resources

- [Data Format Guide](FORMAT.md)
- [JSON Schema](../schemas/question.schema.json)
- [Question Template](../datasets/new-question.json)
