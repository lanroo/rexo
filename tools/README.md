# Documentation tools

The REXO Architectural Constitution is generated from source, not hand-edited.

| File | Purpose |
| --- | --- |
| `rexo_pdf_theme.py` | Shared visual system: palette, typography, tables, code blocks, diagrams and page templates. Both editions import it, so they cannot drift apart visually. |
| `generate_rexo_architecture_pdf_pt.py` | Portuguese (PT-BR) edition. Owns the prose only. |
| `generate_rexo_architecture_pdf_en.py` | English edition. Owns the prose only. |

To rebuild both editions:

```shell
python -m pip install -r requirements-docs.txt
python tools/generate_rexo_architecture_pdf_pt.py
python tools/generate_rexo_architecture_pdf_en.py
```

The palette is inherited from the v2.0 document and is treated as brand
identity: do not change the colours in `rexo_pdf_theme.py` without an ADR.
Short strings baked into the page chrome and the diagram artwork live in that
module under `CHROME` and `DIAGRAM_LABELS`; a new language means adding one
entry to each and one generator script.

Generated files under `output/` and visual QA renders under `tmp/` are not
committed. The reviewed PDF snapshots are stored in `docs/architecture/`. The
historical `AIOS_Constituicao_Arquitetural_v2.0.pdf` predates the REXO rename
and is kept for traceability only; its generator was superseded by the scripts
above.
