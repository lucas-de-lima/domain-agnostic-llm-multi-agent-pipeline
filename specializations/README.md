# Specializations (Use-case presets)

A **specialization** is a preset of prompts for the same 4-stage pipeline (Context, Extract, Synthesize, Audit). It tailors the pipeline to a specific use case (e.g. recipes, programming lessons, morning news) without changing the core agents.

## Adding a new specialization via PR

1. Create a new folder under `specializations/<slug>/` (e.g. `specializations/morning-news/`).
2. Add a **manifest** and four prompt files as described below.
3. Open a PR; once merged, users can run the pipeline with `--specialization=<slug>`.

## Folder structure

```
specializations/
  <slug>/           # e.g. recipes, programming-lessons, morning-news
    manifest.json   # required: id, name, description
    context.txt     # required: prompt for Agent 0 (domain/role classification)
    extract.txt     # required: prompt for Agent 1 (structured extraction)
    synthesize.txt  # required: prompt for Agent 2 (document synthesis)
    audit.txt       # required: prompt for Agent 3 (validation vs original)
```

## manifest.json

```json
{
  "id": "your-slug",
  "name": "Human-readable name",
  "description": "Short description of what this specialization does."
}
```

- `id`: must match the folder name (slug). Used in `--specialization=your-slug`.
- `name`: display name.
- `description`: one or two sentences for docs and UI.

## Prompt files

- **context.txt**: Used by Agent 0. No placeholders; instruct the model to return a JSON with `main_subject`, `complexity_level`, `expert_role_1`, `expert_role_2`, `expert_role_3`, `target_audience`.
- **extract.txt**, **synthesize.txt**, **audit.txt**: Can use Go template placeholders filled from the context JSON after Agent 0 runs:
  - `{{.MainSubject}}`
  - `{{.ComplexityLevel}}`
  - `{{.ExpertRole1}}`, `{{.ExpertRole2}}`, `{{.ExpertRole3}}`
  - `{{.TargetAudience}}`

For **audit.txt**, the pipeline will provide the draft and the original source in the user message; your prompt should instruct the model to compare them and fix hallucinations or inconsistencies.

## Example: recipes

See `specializations/recipes/` for a full example that turns recipe videos or text into formatted written recipes.
