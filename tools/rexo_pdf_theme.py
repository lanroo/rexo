"""Shared visual system for the REXO Architectural Constitution PDFs.

The palette and page geometry are inherited from the v2.0 document and must be
kept stable across languages and revisions. Content lives in the per-language
generator scripts; everything that decides how the document *looks* lives here
so the Portuguese and English editions can never drift apart visually.
"""

from pathlib import Path
from xml.sax.saxutils import escape

from reportlab.graphics.shapes import Circle, Drawing, Line, Polygon, Rect, String
from reportlab.lib import colors
from reportlab.lib.colors import HexColor
from reportlab.lib.enums import TA_CENTER, TA_LEFT, TA_RIGHT
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle, getSampleStyleSheet
from reportlab.lib.units import mm
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import (
    BaseDocTemplate,
    Frame,
    HRFlowable,
    KeepTogether,
    LongTable,
    NextPageTemplate,
    PageBreak,
    PageTemplate,
    Paragraph,
    Preformatted,
    Spacer,
    Table,
    TableStyle,
)
from reportlab.platypus.tableofcontents import TableOfContents


PAGE_W, PAGE_H = A4

# --- v2.0 palette (do not change without an ADR) -----------------------------
NAVY = HexColor("#0B132B")
NAVY_2 = HexColor("#16213E")
BLUE = HexColor("#246BFD")
CYAN = HexColor("#42D3FF")
PURPLE = HexColor("#7758FF")
GREEN = HexColor("#17B890")
ORANGE = HexColor("#FF9F43")
RED = HexColor("#E65353")
INK = HexColor("#172033")
MUTED = HexColor("#5E6B82")
LIGHT = HexColor("#F4F7FB")
LINE = HexColor("#D9E1EC")
WHITE = colors.white

ACCENTS = [BLUE, CYAN, PURPLE, GREEN, ORANGE, RED]

MARGIN_LR = 18 * mm
MARGIN_TOP = 20 * mm
MARGIN_BOTTOM = 17 * mm
CONTENT_W = PAGE_W - 2 * MARGIN_LR


def register_fonts():
    font_dir = Path("C:/Windows/Fonts")
    candidates = {
        "UI": font_dir / "segoeui.ttf",
        "UI-Bold": font_dir / "segoeuib.ttf",
        "UI-Semibold": font_dir / "seguisb.ttf",
        "UI-Light": font_dir / "segoeuil.ttf",
        "Mono": font_dir / "consola.ttf",
        "Mono-Bold": font_dir / "consolab.ttf",
    }
    for name, path in candidates.items():
        if path.exists():
            try:
                pdfmetrics.registerFont(TTFont(name, str(path)))
            except Exception:
                pass
    have = pdfmetrics.getRegisteredFontNames()
    return {
        "regular": "UI" if "UI" in have else "Helvetica",
        "bold": "UI-Bold" if "UI-Bold" in have else "Helvetica-Bold",
        "semibold": "UI-Semibold" if "UI-Semibold" in have else "Helvetica-Bold",
        "light": "UI-Light" if "UI-Light" in have else "Helvetica",
        "mono": "Mono" if "Mono" in have else "Courier",
        "mono_bold": "Mono-Bold" if "Mono-Bold" in have else "Courier-Bold",
    }


FONTS = register_fonts()


# --- localisation ------------------------------------------------------------
# Short strings baked into the chrome (headers, footers) and into the diagram
# artwork. Prose stays in the per-language generator scripts.

DOC_STRINGS = {}
DIAG = {}

CHROME = {
    "pt": {
        "header_left": "REXO · CONSTITUIÇÃO ARQUITETURAL",
        "header_right": "v3.0 — Julho de 2026",
        "footer_left": "Documento de direção arquitetural",
        "page_word": "Página",
        "toc_title": "Sumário",
        "part_word": "Parte",
    },
    "en": {
        "header_left": "REXO · ARCHITECTURAL CONSTITUTION",
        "header_right": "v3.0 — July 2026",
        "footer_left": "Architectural direction document",
        "page_word": "Page",
        "toc_title": "Contents",
        "part_word": "Part",
    },
}

DIAGRAM_LABELS = {
    "pt": {
        "arch_title": "REXO — arquitetura por planos",
        "arch_apps": ("Apps e Distributions", "Course Builder, SaaS Factory, Creative Studio, APIs e UI"),
        "arch_control": ("Control Plane", "Intent Gateway, Goal Compiler, Planner, Orchestrator, Scheduler e Policy Engine"),
        "arch_cap": ("Capability Fabric", "Registry + Resolver"),
        "arch_tool": ("Tool Fabric", "Gateway + providers"),
        "arch_mem": ("Memory Fabric", "Context + knowledge"),
        "arch_exec": ("Execution Plane", "Workers efêmeros e sandboxes"),
        "arch_qual": ("Quality e Governance", "Evals, security, audit e approvals"),
        "cap_req": ("Capability request", "media.generate-image@2"),
        "cap_res": ("Resolver", "constraints + policy"),
        "cap_sel": ("Plano selecionado", "provider + evaluator"),
        "cap_providers": [
            ("ComfyUI + Flux", "privado / local"),
            ("ComfyUI + SDXL", "baixo custo"),
            ("API de imagem cloud", "rápido"),
            ("Provider especialista", "alta fidelidade"),
        ],
        "life_steps": ["Objetivo", "Plan DAG", "Resolver", "Executar", "Avaliar", "Entregar"],
        "life_caption": "Falhas geram tarefas de correção específicas; runs permanecem retomáveis, auditáveis e sensíveis ao budget.",
        "mem_layers": [
            ("Task Working Memory", "mínima, efêmera, limitada por tokens"),
            ("Run Memory", "plano, estado, evidências e decisões"),
            ("Project Memory", "arquitetura, convenções e histórico"),
            ("Organizational Knowledge", "padrões validados entre projetos"),
            ("Provenance e Artifact Store", "fontes, hashes, versões e lineage"),
        ],
        "evo_steps": ["Observar", "Candidato", "Sandbox", "Evals", "Canário", "Promover"],
        "evo_rollback": "rollback em regressão",
        "dev_title": "Developer & Creation Plane",
        "dev_surfaces": [
            ("SDK", "programático"),
            ("CLI", "automação"),
            ("Studio", "IDE visual"),
            ("Canvas", "views de grafo"),
            ("Creator", "fábrica agentic"),
        ],
        "dev_contracts": ("Contratos e manifestos canônicos", "um formato para toda superfície de criação"),
        "dev_validators": ("Validators, Conformance e Registry", "validar, testar, empacotar, assinar e publicar"),
        "dev_note": "Nenhum formato proprietário de Studio ou Creator",
        "creator_steps": ["Platform Intent", "Domain Design", "Gerar", "Validar", "Sandbox", "Publicar"],
        "creator_note1": "Saídas: Domain Pack, policies, conhecimento, Skills, workflows, evaluators, Apps e Distribution.",
        "creator_note2": "O Creator produz candidatos. Publicação estável ainda exige evals, revisão de segurança e aprovação.",
        "mode_boxes": [
            ("Architect Mode", "criticar e decidir"),
            ("Engineering Mode", "construir fase aprovada"),
            ("Operator Mode", "executar workflow governado"),
            ("Creator Mode", "desenhar uma distribuição"),
            ("Governance Mode", "revisar, auditar, liberar"),
        ],
        "mode_note": "Transições são explícitas e controladas; o escopo de aprovação nunca se expande implicitamente.",
        "econ_title": "Economy Engine — sequência de decisão",
        "econ_steps": [
            ("Reuso", "artifact válido?"),
            ("Incremental", "só o que mudou"),
            ("Contexto", "mínimo necessário"),
            ("Provider", "menor adequado"),
            ("Budget", "tokens, tempo, custo"),
            ("Execução", "com piso de qualidade"),
        ],
        "econ_note": "Reutilizar antes de gerar. Gerar o mínimo necessário. Nunca abaixo do piso de qualidade.",
    },
    "en": {
        "arch_title": "REXO — architecture by planes",
        "arch_apps": ("Apps and Distributions", "Course Builder, SaaS Factory, Creative Studio, APIs and UI"),
        "arch_control": ("Control Plane", "Intent Gateway, Goal Compiler, Planner, Orchestrator, Scheduler and Policy Engine"),
        "arch_cap": ("Capability Fabric", "Registry + Resolver"),
        "arch_tool": ("Tool Fabric", "Gateway + providers"),
        "arch_mem": ("Memory Fabric", "Context + knowledge"),
        "arch_exec": ("Execution Plane", "Ephemeral workers and sandboxes"),
        "arch_qual": ("Quality and Governance", "Evals, security, audit and approvals"),
        "cap_req": ("Capability request", "media.generate-image@2"),
        "cap_res": ("Resolver", "constraints + policy"),
        "cap_sel": ("Selected plan", "provider + evaluator"),
        "cap_providers": [
            ("ComfyUI + Flux", "private / local"),
            ("ComfyUI + SDXL", "low cost"),
            ("Cloud image API", "fast"),
            ("Specialist provider", "high fidelity"),
        ],
        "life_steps": ["Goal", "Plan DAG", "Resolve", "Execute", "Evaluate", "Deliver"],
        "life_caption": "Failures create targeted correction tasks; runs remain resumable, auditable and budget-aware.",
        "mem_layers": [
            ("Task Working Memory", "minimal, ephemeral, token-budgeted"),
            ("Run Memory", "plan, state, evidence and decisions"),
            ("Project Memory", "architecture, conventions and history"),
            ("Organizational Knowledge", "validated cross-project patterns"),
            ("Provenance and Artifact Store", "sources, hashes, versions and lineage"),
        ],
        "evo_steps": ["Observe", "Candidate", "Sandbox", "Evals", "Canary", "Promote"],
        "evo_rollback": "rollback on regression",
        "dev_title": "Developer & Creation Plane",
        "dev_surfaces": [
            ("SDK", "programmatic"),
            ("CLI", "automation"),
            ("Studio", "visual IDE"),
            ("Canvas", "graph views"),
            ("Creator", "agentic factory"),
        ],
        "dev_contracts": ("Canonical Contracts and Manifests", "one format for every creation surface"),
        "dev_validators": ("Validators, Conformance and Registry", "validate, test, package, sign and publish"),
        "dev_note": "No proprietary Studio or Creator format",
        "creator_steps": ["Platform Intent", "Domain Design", "Generate", "Validate", "Sandbox", "Publish"],
        "creator_note1": "Outputs: Domain Pack, policies, knowledge, Skills, workflows, evaluators, Apps and Distribution.",
        "creator_note2": "Creator produces candidates. Stable publication still requires evals, security review and approval.",
        "mode_boxes": [
            ("Architect Mode", "critique and decide"),
            ("Engineering Mode", "build approved phase"),
            ("Operator Mode", "run governed workflow"),
            ("Creator Mode", "design a distribution"),
            ("Governance Mode", "review, audit, release"),
        ],
        "mode_note": "Transitions are explicit and gated; approval scope never expands implicitly.",
        "econ_title": "Economy Engine — decision sequence",
        "econ_steps": [
            ("Reuse", "valid artifact?"),
            ("Incremental", "only what changed"),
            ("Context", "minimum required"),
            ("Provider", "smallest adequate"),
            ("Budget", "tokens, time, cost"),
            ("Execute", "above quality floor"),
        ],
        "econ_note": "Reuse before generating. Generate only what is needed. Never below the quality floor.",
    },
}


def use_language(lang):
    """Select the language used by the page chrome and the diagram artwork."""
    DOC_STRINGS.clear()
    DOC_STRINGS.update(CHROME[lang])
    DIAG.clear()
    DIAG.update(DIAGRAM_LABELS[lang])


use_language("pt")


# --- styles ------------------------------------------------------------------
def make_styles():
    base = getSampleStyleSheet()
    s = {}
    s["CoverBrand"] = ParagraphStyle(
        "CoverBrand", fontName=FONTS["semibold"], fontSize=9, leading=12,
        textColor=CYAN, alignment=TA_LEFT,
    )
    s["CoverTitle"] = ParagraphStyle(
        "CoverTitle", parent=base["Title"], fontName=FONTS["bold"], fontSize=44,
        leading=48, textColor=WHITE, alignment=TA_LEFT, spaceAfter=0,
    )
    s["CoverTitle2"] = ParagraphStyle(
        "CoverTitle2", parent=base["Title"], fontName=FONTS["light"], fontSize=27,
        leading=32, textColor=HexColor("#C9DBFF"), alignment=TA_LEFT, spaceAfter=0,
    )
    s["CoverSubtitle"] = ParagraphStyle(
        "CoverSubtitle", fontName=FONTS["regular"], fontSize=11.5, leading=18,
        textColor=HexColor("#D9E8FF"), alignment=TA_LEFT,
    )
    s["CoverAuthor"] = ParagraphStyle(
        "CoverAuthor", fontName=FONTS["semibold"], fontSize=13, leading=17,
        textColor=WHITE,
    )
    s["CoverAuthorRole"] = ParagraphStyle(
        "CoverAuthorRole", fontName=FONTS["regular"], fontSize=8.4, leading=12,
        textColor=HexColor("#9FB6DA"),
    )
    s["CoverMetaKey"] = ParagraphStyle(
        "CoverMetaKey", fontName=FONTS["semibold"], fontSize=6.8, leading=10,
        textColor=HexColor("#6E86AE"),
    )
    s["CoverMetaVal"] = ParagraphStyle(
        "CoverMetaVal", fontName=FONTS["regular"], fontSize=8.4, leading=12,
        textColor=HexColor("#D5E3FA"),
    )
    s["PartKicker"] = ParagraphStyle(
        "PartKicker", fontName=FONTS["semibold"], fontSize=9, leading=12,
        textColor=CYAN,
    )
    s["PartNumber"] = ParagraphStyle(
        "PartNumber", fontName=FONTS["bold"], fontSize=96, leading=100,
        textColor=HexColor("#31558F"), alignment=TA_RIGHT,
    )
    s["PartTitle"] = ParagraphStyle(
        "PartTitle", fontName=FONTS["bold"], fontSize=28, leading=34, textColor=WHITE,
    )
    s["PartSubtitle"] = ParagraphStyle(
        "PartSubtitle", fontName=FONTS["regular"], fontSize=10.5, leading=16,
        textColor=HexColor("#AFC5E8"),
    )
    s["Eyebrow"] = ParagraphStyle(
        "Eyebrow", fontName=FONTS["bold"], fontSize=7.6, leading=10, textColor=BLUE,
        spaceAfter=1.6 * mm,
    )
    s["H1"] = ParagraphStyle(
        "H1", fontName=FONTS["bold"], fontSize=21, leading=26, textColor=NAVY,
        spaceBefore=1 * mm, spaceAfter=2.5 * mm, keepWithNext=True,
    )
    s["H2"] = ParagraphStyle(
        "H2", fontName=FONTS["semibold"], fontSize=13, leading=17, textColor=NAVY_2,
        spaceBefore=5 * mm, spaceAfter=2 * mm, keepWithNext=True,
    )
    s["H3"] = ParagraphStyle(
        "H3", fontName=FONTS["semibold"], fontSize=10.4, leading=13.5, textColor=BLUE,
        spaceBefore=3.5 * mm, spaceAfter=1.4 * mm, keepWithNext=True,
    )
    s["Body"] = ParagraphStyle(
        "Body", fontName=FONTS["regular"], fontSize=9.5, leading=14.6, textColor=INK,
        spaceAfter=2.4 * mm,
    )
    s["Lead"] = ParagraphStyle(
        "Lead", fontName=FONTS["regular"], fontSize=10.6, leading=16.4,
        textColor=HexColor("#2A3550"), spaceAfter=3 * mm,
    )
    s["Small"] = ParagraphStyle(
        "Small", fontName=FONTS["regular"], fontSize=7.8, leading=11, textColor=MUTED,
    )
    s["Caption"] = ParagraphStyle(
        "Caption", fontName=FONTS["regular"], fontSize=7.4, leading=10.5,
        textColor=MUTED, spaceBefore=1.4 * mm, spaceAfter=3 * mm,
    )
    s["Bullet"] = ParagraphStyle(
        "Bullet", parent=s["Body"], leftIndent=5.4 * mm, firstLineIndent=-3.2 * mm,
        bulletIndent=1.5 * mm, spaceAfter=1.2 * mm,
    )
    s["Number"] = ParagraphStyle(
        "Number", parent=s["Body"], leftIndent=6.4 * mm, firstLineIndent=-4.4 * mm,
        bulletIndent=1 * mm, spaceAfter=1.5 * mm,
    )
    s["CalloutText"] = ParagraphStyle(
        "CalloutText", fontName=FONTS["regular"], fontSize=9.6, leading=15, textColor=NAVY,
    )
    s["QuoteText"] = ParagraphStyle(
        "QuoteText", fontName=FONTS["semibold"], fontSize=12.5, leading=18.5,
        textColor=WHITE, alignment=TA_LEFT,
    )
    # Paragraph-flavoured variants kept so prose can stay declarative.
    # backColor + borderPadding paints outside the text box, so the indents
    # must match the padding or the panel bleeds past the page margin.
    s["Quote"] = ParagraphStyle(
        "Quote", fontName=FONTS["semibold"], fontSize=12.5, leading=18.5,
        textColor=WHITE, backColor=PURPLE, borderColor=PURPLE, borderWidth=0,
        borderPadding=13, leftIndent=13, rightIndent=13,
        spaceBefore=5 * mm, spaceAfter=7 * mm,
    )
    s["Callout"] = ParagraphStyle(
        "Callout", fontName=FONTS["regular"], fontSize=9.6, leading=15,
        textColor=NAVY, backColor=HexColor("#EEF4FF"), borderColor=HexColor("#BFD3FF"),
        borderWidth=0.9, borderPadding=11, leftIndent=11, rightIndent=11,
        spaceBefore=4 * mm, spaceAfter=6 * mm,
    )
    # Used *inside* the code() panel, which paints its own background.
    s["CodeInner"] = ParagraphStyle(
        "CodeInner", fontName=FONTS["mono"], fontSize=7.6, leading=10.6,
        textColor=HexColor("#DCE9FF"),
    )
    # Standalone dark code panel. Carries its own background so that a bare
    # P(text, "Code") is never light-on-white.
    s["Code"] = ParagraphStyle(
        "Code", fontName=FONTS["mono"], fontSize=7.6, leading=11.4,
        textColor=HexColor("#DCE9FF"), backColor=NAVY, borderColor=NAVY,
        borderPadding=11, leftIndent=11, rightIndent=11,
        spaceBefore=3 * mm, spaceAfter=5 * mm,
    )
    # Light panel for illustrative snippets that carry inline markup (<b>),
    # which a Preformatted block cannot render.
    s["Sample"] = ParagraphStyle(
        "Sample", fontName=FONTS["mono"], fontSize=7.8, leading=12,
        textColor=NAVY_2, backColor=HexColor("#EEF2F8"), borderColor=HexColor("#C9D6E8"),
        borderWidth=0.8, borderPadding=11, leftIndent=11, rightIndent=11,
        spaceBefore=3 * mm, spaceAfter=5 * mm,
    )
    s["Pre"] = ParagraphStyle(
        "Pre", fontName=FONTS["mono"], fontSize=7.0, leading=9.6, textColor=NAVY_2,
    )
    s["TableHead"] = ParagraphStyle(
        "TableHead", fontName=FONTS["semibold"], fontSize=8.1, leading=10.5,
        textColor=WHITE,
    )
    s["TableBody"] = ParagraphStyle(
        "TableBody", fontName=FONTS["regular"], fontSize=7.8, leading=10.6,
        textColor=INK,
    )
    s["TOCHeading"] = ParagraphStyle(
        "TOCHeading", fontName=FONTS["bold"], fontSize=26, leading=30, textColor=NAVY,
        spaceAfter=7 * mm,
    )
    return s


ST = make_styles()


# --- inline helpers ----------------------------------------------------------
def P(text, style="Body"):
    return Paragraph(text, ST[style])


def bullet(text):
    return Paragraph(f'<font color="#246BFD">■</font>&nbsp; {text}', ST["Bullet"])


def numbered(n, text):
    return Paragraph(
        f'<font name="{FONTS["bold"]}" color="#246BFD">{n}.</font>&nbsp; {text}',
        ST["Number"],
    )


def section(title, eyebrow=None):
    """Chapter opener: coloured eyebrow, title, and a two-tone rule."""
    out = [Spacer(1, 3 * mm)]
    if eyebrow:
        out.append(P(escape(eyebrow).upper(), "Eyebrow"))
    out.append(P(title, "H1"))
    rule = Drawing(CONTENT_W, 3)
    rule.add(Rect(0, 0.6, 26, 2, fillColor=BLUE, strokeColor=None))
    rule.add(Rect(26, 1.3, CONTENT_W - 26, 0.6, fillColor=LINE, strokeColor=None))
    out.append(rule)
    out.append(Spacer(1, 3.4 * mm))
    return out


def sub(title):
    return P(title, "H2")


def sub3(title):
    return P(title, "H3")


def caption(text):
    return P(escape(text), "Caption")


def table(headers, rows, widths=None):
    """Editorial table: solid header band, hairline rows, zebra body."""
    data = [[P(escape(str(h)), "TableHead") for h in headers]]
    for row in rows:
        data.append([P(str(cell), "TableBody") for cell in row])
    t = LongTable(data, colWidths=widths, repeatRows=1, hAlign="LEFT")
    style = [
        ("BACKGROUND", (0, 0), (-1, 0), NAVY_2),
        ("VALIGN", (0, 0), (-1, -1), "TOP"),
        ("ROWBACKGROUNDS", (0, 1), (-1, -1), [WHITE, LIGHT]),
        ("LINEBELOW", (0, 1), (-1, -1), 0.4, LINE),
        ("LINEBELOW", (0, 0), (-1, 0), 1.1, BLUE),
        ("LEFTPADDING", (0, 0), (-1, -1), 7),
        ("RIGHTPADDING", (0, 0), (-1, -1), 7),
        ("TOPPADDING", (0, 0), (-1, -1), 5.5),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 5.5),
    ]
    t.setStyle(TableStyle(style))
    return t


def code(text, dark=True):
    """Monospaced block that actually preserves its line breaks.

    The v2/v3 generators passed multi-line snippets through Paragraph, which
    collapses newlines into spaces and destroyed every YAML and tree listing.
    """
    body = Preformatted(text.strip("\n"), ST["CodeInner"] if dark else ST["Pre"])
    t = Table([[body]], colWidths=[CONTENT_W], hAlign="LEFT")
    cmds = [
        ("BACKGROUND", (0, 0), (-1, -1), NAVY if dark else LIGHT),
        ("LEFTPADDING", (0, 0), (-1, -1), 11),
        ("RIGHTPADDING", (0, 0), (-1, -1), 11),
        ("TOPPADDING", (0, 0), (-1, -1), 9),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 9),
        ("LINEBEFORE", (0, 0), (0, -1), 2.2, CYAN if dark else BLUE),
        ("ROUNDEDCORNERS", [4, 4, 4, 4]),
    ]
    try:
        t.setStyle(TableStyle(cmds))
    except Exception:  # older reportlab without ROUNDEDCORNERS
        t.setStyle(TableStyle(cmds[:-1]))
    return KeepTogether([Spacer(1, 1.5 * mm), t, Spacer(1, 3.5 * mm)])


def callout(text, accent=BLUE):
    """Tinted panel with an accent bar, used for interpretation rules."""
    t = Table([[P(text, "CalloutText")]], colWidths=[CONTENT_W], hAlign="LEFT")
    t.setStyle(
        TableStyle(
            [
                ("BACKGROUND", (0, 0), (-1, -1), HexColor("#EEF4FF")),
                ("LINEBEFORE", (0, 0), (0, -1), 2.6, accent),
                ("LEFTPADDING", (0, 0), (-1, -1), 11),
                ("RIGHTPADDING", (0, 0), (-1, -1), 11),
                ("TOPPADDING", (0, 0), (-1, -1), 9),
                ("BOTTOMPADDING", (0, 0), (-1, -1), 9),
            ]
        )
    )
    return KeepTogether([Spacer(1, 2 * mm), t, Spacer(1, 4 * mm)])


def quote(text):
    """Full-width statement panel in the brand purple."""
    t = Table([[P(text, "QuoteText")]], colWidths=[CONTENT_W], hAlign="LEFT")
    cmds = [
        ("BACKGROUND", (0, 0), (-1, -1), PURPLE),
        ("LEFTPADDING", (0, 0), (-1, -1), 14),
        ("RIGHTPADDING", (0, 0), (-1, -1), 14),
        ("TOPPADDING", (0, 0), (-1, -1), 12),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 12),
        ("ROUNDEDCORNERS", [5, 5, 5, 5]),
    ]
    try:
        t.setStyle(TableStyle(cmds))
    except Exception:
        t.setStyle(TableStyle(cmds[:-1]))
    return KeepTogether([Spacer(1, 2.5 * mm), t, Spacer(1, 4.5 * mm)])


# --- diagram primitives ------------------------------------------------------
def box(d, x, y, w, h, title, subtitle="", fill=WHITE, stroke=LINE, title_color=NAVY):
    d.add(Rect(x, y, w, h, rx=6, ry=6, fillColor=fill, strokeColor=stroke, strokeWidth=1))
    d.add(String(x + 9, y + h - 14, title, fontName=FONTS["semibold"], fontSize=8.5, fillColor=title_color))
    if subtitle:
        d.add(String(x + 9, y + 8, subtitle, fontName=FONTS["regular"], fontSize=6.6, fillColor=MUTED))


def arrow(d, x1, y1, x2, y2, color=BLUE):
    d.add(Line(x1, y1, x2, y2, strokeColor=color, strokeWidth=1.4))
    if abs(x2 - x1) >= abs(y2 - y1):
        direction = 1 if x2 > x1 else -1
        pts = [x2, y2, x2 - 5 * direction, y2 + 3, x2 - 5 * direction, y2 - 3]
    else:
        direction = 1 if y2 > y1 else -1
        pts = [x2, y2, x2 - 3, y2 - 5 * direction, x2 + 3, y2 - 5 * direction]
    d.add(Polygon(pts, fillColor=color, strokeColor=color))


def _canvas(w, h, title=None):
    d = Drawing(w, h)
    d.add(Rect(0, 0, w, h, rx=10, ry=10, fillColor=LIGHT, strokeColor=LINE))
    if title:
        d.add(Rect(0, h - 26, 4, 22, fillColor=BLUE, strokeColor=None))
        d.add(String(14, h - 19, title, fontName=FONTS["bold"], fontSize=11, fillColor=NAVY))
    return d


def _chain(d, steps, y, xs, colors_list, box_w=58, box_h=36, font=7.4):
    for i, ((label, _sub), x) in enumerate(zip(steps, xs)):
        color = colors_list[i % len(colors_list)]
        d.add(Rect(x, y, box_w, box_h, rx=6, ry=6, fillColor=WHITE, strokeColor=color, strokeWidth=1.2))
        d.add(String(x + box_w / 2, y + box_h / 2 + 4, label, textAnchor="middle",
                     fontName=FONTS["semibold"], fontSize=font, fillColor=color))
        if _sub:
            d.add(String(x + box_w / 2, y + box_h / 2 - 7, _sub, textAnchor="middle",
                         fontName=FONTS["regular"], fontSize=5.8, fillColor=MUTED))
        if i < len(xs) - 1:
            arrow(d, x + box_w, y + box_h / 2, xs[i + 1], y + box_h / 2, LINE)


# --- diagrams ----------------------------------------------------------------
def architecture_diagram():
    d = _canvas(480, 300, DIAG["arch_title"])
    box(d, 20, 222, 440, 36, *DIAG["arch_apps"], HexColor("#EAF2FF"), HexColor("#AFC8FF"), BLUE)
    box(d, 20, 163, 440, 42, *DIAG["arch_control"], WHITE, LINE)
    for x, key, color in [(20, "arch_cap", BLUE), (170, "arch_tool", PURPLE), (320, "arch_mem", GREEN)]:
        box(d, x, 94, 140, 48, *DIAG[key], WHITE, color, color)
    box(d, 20, 30, 215, 42, *DIAG["arch_exec"], WHITE, LINE)
    box(d, 245, 30, 215, 42, *DIAG["arch_qual"], WHITE, LINE)
    arrow(d, 240, 222, 240, 205)
    arrow(d, 240, 163, 240, 143)
    arrow(d, 240, 94, 150, 72)
    arrow(d, 240, 94, 350, 72)
    return d


def capability_diagram():
    d = _canvas(480, 205)
    box(d, 15, 132, 115, 45, *DIAG["cap_req"], HexColor("#EAF2FF"), BLUE, BLUE)
    box(d, 182, 132, 115, 45, *DIAG["cap_res"], WHITE, PURPLE, PURPLE)
    box(d, 350, 132, 115, 45, *DIAG["cap_sel"], HexColor("#EFFFFA"), GREEN, GREEN)
    arrow(d, 130, 154, 182, 154)
    arrow(d, 297, 154, 350, 154)
    for x, (title, subtext) in zip([15, 132, 249, 366], DIAG["cap_providers"]):
        box(d, x, 42, 100, 43, title, subtext, WHITE, LINE)
        arrow(d, 407, 132, x + 50, 85, HexColor("#B9C6DA"))
    return d


def lifecycle_diagram():
    d = _canvas(480, 160)
    steps = [(s, "") for s in DIAG["life_steps"]]
    _chain(d, steps, 72, [12, 92, 172, 252, 332, 412], [BLUE, PURPLE, GREEN, ORANGE, RED, NAVY])
    d.add(String(14, 30, DIAG["life_caption"], fontName=FONTS["regular"], fontSize=7.3, fillColor=MUTED))
    return d


def memory_diagram():
    d = _canvas(480, 230)
    geometry = [
        (35, 174, 410, HexColor("#EAF2FF"), BLUE),
        (55, 133, 370, HexColor("#F0ECFF"), PURPLE),
        (75, 92, 330, HexColor("#EFFFFA"), GREEN),
        (95, 51, 290, HexColor("#FFF5E8"), ORANGE),
        (115, 10, 250, WHITE, NAVY),
    ]
    for (x, y, w, fill, stroke), (title, subtitle) in zip(geometry, DIAG["mem_layers"]):
        box(d, x, y, w, 32, title, subtitle, fill, stroke, stroke)
    return d


def evolution_diagram():
    d = _canvas(480, 150)
    steps = [(s, "") for s in DIAG["evo_steps"]]
    _chain(d, steps, 72, [10, 90, 170, 250, 330, 410], [BLUE, PURPLE, ORANGE, RED, GREEN, NAVY], box_w=60, box_h=34, font=7.2)
    d.add(Line(440, 72, 440, 35, strokeColor=RED, strokeWidth=1))
    d.add(Line(440, 35, 200, 35, strokeColor=RED, strokeWidth=1))
    d.add(Polygon([200, 35, 206, 38, 206, 32], fillColor=RED, strokeColor=RED))
    d.add(String(310, 21, DIAG["evo_rollback"], textAnchor="middle", fontName=FONTS["regular"], fontSize=7.2, fillColor=RED))
    return d


def developer_plane_diagram():
    d = _canvas(480, 250, DIAG["dev_title"])
    xs = [16, 108, 200, 292, 384]
    for x, (title, subtitle) in zip(xs, DIAG["dev_surfaces"]):
        box(d, x, 164, 80, 40, title, subtitle, WHITE, BLUE, BLUE)
        arrow(d, x + 40, 164, 240, 136, HexColor("#B9C6DA"))
    box(d, 95, 94, 290, 40, *DIAG["dev_contracts"], HexColor("#F0ECFF"), PURPLE, PURPLE)
    arrow(d, 240, 94, 240, 78)
    box(d, 95, 36, 290, 40, *DIAG["dev_validators"], HexColor("#EFFFFA"), GREEN, GREEN)
    d.add(String(240, 16, DIAG["dev_note"], textAnchor="middle", fontName=FONTS["semibold"], fontSize=7.4, fillColor=MUTED))
    return d


def creator_diagram():
    d = _canvas(480, 180)
    steps = [(s, "") for s in DIAG["creator_steps"]]
    _chain(d, steps, 96, [10, 90, 170, 250, 330, 410], [BLUE, PURPLE, ORANGE, RED, GREEN, NAVY], box_w=60, box_h=38, font=6.6)
    d.add(String(20, 62, DIAG["creator_note1"], fontName=FONTS["regular"], fontSize=7.1, fillColor=MUTED))
    d.add(String(20, 38, DIAG["creator_note2"], fontName=FONTS["semibold"], fontSize=7.1, fillColor=RED))
    return d


def mode_diagram():
    d = _canvas(480, 210)
    top = DIAG["mode_boxes"][:3]
    bottom = DIAG["mode_boxes"][3:]
    fills = [HexColor("#EAF2FF"), HexColor("#FFF5E8"), HexColor("#EFFFFA")]
    strokes = [BLUE, ORANGE, GREEN]
    for x, (title, subtitle), fill, stroke in zip([15, 177, 340], top, fills, strokes):
        box(d, x, 138, 125, 42, title, subtitle, fill, stroke, stroke)
    arrow(d, 140, 159, 177, 159)
    arrow(d, 302, 159, 340, 159)
    box(d, 95, 65, 125, 42, *bottom[0], HexColor("#F0ECFF"), PURPLE, PURPLE)
    box(d, 260, 65, 125, 42, *bottom[1], WHITE, RED, RED)
    arrow(d, 157, 107, 110, 138, HexColor("#B9C6DA"))
    arrow(d, 322, 107, 402, 138, HexColor("#B9C6DA"))
    d.add(String(240, 28, DIAG["mode_note"], textAnchor="middle", fontName=FONTS["semibold"], fontSize=7.3, fillColor=NAVY))
    return d


def economy_diagram():
    """New in v3: the Economy Engine decision sequence."""
    d = _canvas(480, 175, DIAG["econ_title"])
    _chain(d, DIAG["econ_steps"], 74, [10, 90, 170, 250, 330, 410],
           [BLUE, CYAN, PURPLE, GREEN, ORANGE, NAVY], box_w=60, box_h=42, font=6.8)
    d.add(Rect(10, 26, 460, 22, rx=5, ry=5, fillColor=HexColor("#EAF2FF"), strokeColor=HexColor("#BFD3FF")))
    d.add(String(240, 34, DIAG["econ_note"], textAnchor="middle", fontName=FONTS["semibold"], fontSize=7.6, fillColor=NAVY))
    return d


# --- page furniture ----------------------------------------------------------
def _dot_field(canvas, x0, y0, cols, rows, step, color, alpha=0.18, r=1.1):
    canvas.saveState()
    canvas.setFillColor(color)
    canvas.setFillAlpha(alpha)
    for i in range(cols):
        for j in range(rows):
            canvas.circle(x0 + i * step, y0 + j * step, r, stroke=0, fill=1)
    canvas.restoreState()


def cover_page(canvas, doc):
    canvas.saveState()
    canvas.setFillColor(NAVY)
    canvas.rect(0, 0, PAGE_W, PAGE_H, stroke=0, fill=1)

    # Layered brand geometry, echoing the v2 cover but softer and deeper.
    canvas.setFillColor(BLUE)
    canvas.setFillAlpha(0.42)
    canvas.circle(PAGE_W - 12 * mm, PAGE_H - 16 * mm, 52 * mm, stroke=0, fill=1)
    canvas.setFillColor(PURPLE)
    canvas.setFillAlpha(0.34)
    canvas.circle(PAGE_W - 44 * mm, PAGE_H - 4 * mm, 30 * mm, stroke=0, fill=1)
    canvas.setFillColor(PURPLE)
    canvas.setFillAlpha(0.40)
    canvas.circle(PAGE_W - 26 * mm, 16 * mm, 32 * mm, stroke=0, fill=1)
    canvas.setFillColor(CYAN)
    canvas.setFillAlpha(0.22)
    canvas.circle(14 * mm, 34 * mm, 20 * mm, stroke=0, fill=1)
    canvas.setFillAlpha(1)

    canvas.setStrokeColor(CYAN)
    canvas.setStrokeAlpha(0.55)
    canvas.setLineWidth(0.9)
    canvas.circle(PAGE_W - 12 * mm, PAGE_H - 16 * mm, 66 * mm, stroke=1, fill=0)
    canvas.setStrokeAlpha(1)

    _dot_field(canvas, 20 * mm, 126 * mm, 9, 3, 5.2 * mm, CYAN, alpha=0.30)

    # Accent ribbon under the title block.
    seg_w = 22 * mm
    for i, c in enumerate([BLUE, CYAN, PURPLE]):
        canvas.setFillColor(c)
        canvas.rect(20 * mm + i * seg_w, PAGE_H - 116 * mm, seg_w, 1.6 * mm, stroke=0, fill=1)

    # Baseline rule above the credit block.
    canvas.setStrokeColor(HexColor("#2C3F63"))
    canvas.setLineWidth(0.8)
    canvas.line(20 * mm, 112 * mm, PAGE_W - 20 * mm, 112 * mm)
    canvas.restoreState()


def part_page(canvas, doc):
    canvas.saveState()
    canvas.setFillColor(NAVY)
    canvas.rect(0, 0, PAGE_W, PAGE_H, stroke=0, fill=1)
    canvas.setFillColor(BLUE)
    canvas.setFillAlpha(0.30)
    canvas.circle(PAGE_W - 6 * mm, 26 * mm, 44 * mm, stroke=0, fill=1)
    canvas.setFillColor(PURPLE)
    canvas.setFillAlpha(0.26)
    canvas.circle(PAGE_W - 34 * mm, 6 * mm, 26 * mm, stroke=0, fill=1)
    canvas.setFillAlpha(1)
    _dot_field(canvas, 20 * mm, 30 * mm, 6, 4, 5 * mm, CYAN, alpha=0.22)
    canvas.setFillColor(CYAN)
    canvas.rect(20 * mm, PAGE_H - 42 * mm, 26 * mm, 1.6 * mm, stroke=0, fill=1)
    canvas.restoreState()


def normal_page(canvas, doc):
    canvas.saveState()
    # Top accent bar + header
    canvas.setFillColor(BLUE)
    canvas.rect(0, PAGE_H - 4.2 * mm, PAGE_W * 0.42, 1.6 * mm, stroke=0, fill=1)
    canvas.setFillColor(CYAN)
    canvas.rect(PAGE_W * 0.42, PAGE_H - 4.2 * mm, PAGE_W * 0.18, 1.6 * mm, stroke=0, fill=1)
    canvas.setFillColor(HexColor("#E3EAF5"))
    canvas.rect(PAGE_W * 0.60, PAGE_H - 4.2 * mm, PAGE_W * 0.40, 1.6 * mm, stroke=0, fill=1)

    canvas.setFont(FONTS["semibold"], 7.2)
    canvas.setFillColor(NAVY)
    canvas.drawString(MARGIN_LR, PAGE_H - 12 * mm, DOC_STRINGS["header_left"])
    canvas.setFont(FONTS["regular"], 7.2)
    canvas.setFillColor(MUTED)
    canvas.drawRightString(PAGE_W - MARGIN_LR, PAGE_H - 12 * mm, DOC_STRINGS["header_right"])
    canvas.setStrokeColor(LINE)
    canvas.setLineWidth(0.5)
    canvas.line(MARGIN_LR, PAGE_H - 14.6 * mm, PAGE_W - MARGIN_LR, PAGE_H - 14.6 * mm)

    # Footer
    canvas.line(MARGIN_LR, 12.6 * mm, PAGE_W - MARGIN_LR, 12.6 * mm)
    canvas.setFont(FONTS["regular"], 7.2)
    canvas.setFillColor(MUTED)
    canvas.drawString(MARGIN_LR, 8.4 * mm, DOC_STRINGS["footer_left"])
    canvas.setFillColor(NAVY)
    canvas.setFont(FONTS["semibold"], 7.6)
    canvas.drawRightString(PAGE_W - MARGIN_LR, 8.4 * mm, f"{DOC_STRINGS['page_word']} {doc.page}")
    canvas.restoreState()


# --- structural blocks -------------------------------------------------------
def cover_story(brand, title, subtitle, description, author, author_role, meta_rows):
    """Title page. `meta_rows` is a list of (label, value) pairs."""
    meta = Table(
        [[P(escape(k).upper(), "CoverMetaKey"), P(escape(v), "CoverMetaVal")] for k, v in meta_rows],
        colWidths=[38 * mm, CONTENT_W - 38 * mm - 4 * mm],
        hAlign="LEFT",
    )
    meta.setStyle(
        TableStyle(
            [
                ("VALIGN", (0, 0), (-1, -1), "TOP"),
                ("LEFTPADDING", (0, 0), (-1, -1), 0),
                ("RIGHTPADDING", (0, 0), (-1, -1), 6),
                ("TOPPADDING", (0, 0), (-1, -1), 2.1),
                ("BOTTOMPADDING", (0, 0), (-1, -1), 2.1),
            ]
        )
    )
    return [
        Spacer(1, 24 * mm),
        P(escape(brand).upper(), "CoverBrand"),
        Spacer(1, 6 * mm),
        P(escape(title), "CoverTitle"),
        P(escape(subtitle), "CoverTitle2"),
        Spacer(1, 14 * mm),
        P(escape(description), "CoverSubtitle"),
        Spacer(1, 84 * mm),
        P(escape(author), "CoverAuthor"),
        P(escape(author_role), "CoverAuthorRole"),
        Spacer(1, 7 * mm),
        meta,
        NextPageTemplate("normal"),
        PageBreak(),
    ]


def part_divider(kicker, number, title, subtitle):
    """Full-bleed part opener: text block left, oversized numeral right."""
    left = [
        P(escape(kicker).upper(), "PartKicker"),
        Spacer(1, 4 * mm),
        P(escape(title), "PartTitle"),
        Spacer(1, 3 * mm),
        P(escape(subtitle), "PartSubtitle"),
    ]
    inner_w = PAGE_W - 40 * mm
    head = Table(
        [[left, P(escape(number), "PartNumber")]],
        colWidths=[inner_w * 0.68, inner_w * 0.32],
        hAlign="LEFT",
    )
    head.setStyle(
        TableStyle(
            [
                ("VALIGN", (0, 0), (0, 0), "MIDDLE"),
                ("VALIGN", (1, 0), (1, 0), "TOP"),
                ("LEFTPADDING", (0, 0), (-1, -1), 0),
                ("RIGHTPADDING", (0, 0), (-1, -1), 0),
                ("TOPPADDING", (0, 0), (-1, -1), 0),
                ("BOTTOMPADDING", (0, 0), (-1, -1), 0),
            ]
        )
    )
    return [
        NextPageTemplate("part"),
        PageBreak(),
        Spacer(1, 74 * mm),
        head,
        NextPageTemplate("normal"),
        PageBreak(),
    ]


def toc_block():
    """Contents page with dot leaders."""
    toc = TableOfContents()
    toc.dotsMinLevel = 0
    toc.levelStyles = [
        ParagraphStyle(
            "TOC0", fontName=FONTS["semibold"], fontSize=9, leading=15, textColor=NAVY,
            leftIndent=0, firstLineIndent=0, spaceBefore=2.2,
        ),
        ParagraphStyle(
            "TOC1", fontName=FONTS["regular"], fontSize=7.4, leading=11, textColor=MUTED,
            leftIndent=8 * mm, firstLineIndent=0,
        ),
    ]
    # No trailing PageBreak: whatever follows (a part divider or a section)
    # supplies its own break, and emitting one here leaves a blank page.
    return [P(DOC_STRINGS["toc_title"], "TOCHeading"), toc]


class ArchitectureDocTemplate(BaseDocTemplate):
    def __init__(self, filename, **kwargs):
        super().__init__(filename, **kwargs)
        self._heading_counter = 0

    def beforeDocument(self):
        self._heading_counter = 0

    def afterFlowable(self, flowable):
        if isinstance(flowable, Paragraph):
            style = flowable.style.name
            if style in ("H1", "H2", "PartTitle"):
                text = flowable.getPlainText()
                key = f"heading-{self._heading_counter}"
                self._heading_counter += 1
                self.canv.bookmarkPage(key)
                if style in ("H1", "PartTitle"):
                    self.canv.addOutlineEntry(text, key, level=0, closed=False)
                    if style == "PartTitle":
                        return  # part openers are bookmarks, not TOC rows
                self.notify("TOCEntry", (0 if style == "H1" else 1, text, self.page, key))


def build_pdf(output, story, title, author, subject, keywords=""):
    output = Path(output)
    output.parent.mkdir(parents=True, exist_ok=True)
    doc = ArchitectureDocTemplate(
        str(output),
        pagesize=A4,
        leftMargin=MARGIN_LR,
        rightMargin=MARGIN_LR,
        topMargin=MARGIN_TOP,
        bottomMargin=MARGIN_BOTTOM,
        title=title,
        author=author,
        subject=subject,
        creator="REXO documentation toolchain",
    )
    cover_frame = Frame(20 * mm, 20 * mm, PAGE_W - 40 * mm, PAGE_H - 40 * mm, id="cover", showBoundary=0)
    part_frame = Frame(20 * mm, 20 * mm, PAGE_W - 40 * mm, PAGE_H - 40 * mm, id="part", showBoundary=0)
    normal_frame = Frame(
        MARGIN_LR, MARGIN_BOTTOM,
        PAGE_W - 2 * MARGIN_LR, PAGE_H - MARGIN_TOP - MARGIN_BOTTOM,
        id="normal", showBoundary=0,
    )
    doc.addPageTemplates(
        [
            PageTemplate(id="cover", frames=[cover_frame], onPage=cover_page),
            PageTemplate(id="part", frames=[part_frame], onPage=part_page),
            PageTemplate(id="normal", frames=[normal_frame], onPage=normal_page),
        ]
    )
    doc.multiBuild(story)
    return output
