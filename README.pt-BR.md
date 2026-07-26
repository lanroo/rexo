# REXO

**Runtime for Execution & eXchange Orchestration**

O REXO é um runtime experimental, orientado por capacidades, para construir
sistemas de IA reutilizáveis. Modelos, ferramentas, memória, workflows,
orçamentos e critérios de qualidade são tratados como infraestrutura governada,
e não escondidos dentro de prompts gigantes.

> Estado atual: Fase 1 (walking skeleton). O CLI já executa workflows
> determinísticos de ponta a ponta (`rexo run`); seleção de provider/LLM ainda
> não foi implementada.

[English](README.md) ·
[Instalação](INSTALL.md) ·
[Constituição arquitetural](docs/architecture/constitution.md) ·
[Roadmap Core v1](docs/roadmap/core-v1.md) ·
[Primeiros passos](docs/getting-started/README.md)

## Instalação

| Plataforma | Comando |
|---|---|
| macOS / Linux | `brew install lanroo/tap/rexo` |
| Windows | `scoop bucket add rexo https://github.com/lanroo/rexo` e depois `scoop install rexo` |
| Qualquer (devs Go) | `go install github.com/lanroo/rexo/cmd/rexo@latest` |

Veja o [INSTALL.md](INSTALL.md) para downloads manuais e solução de problemas.

## O que funciona agora

- `rexo version`: mostra versão, sistema e arquitetura.
- `rexo doctor`: verifica a compatibilidade básica do computador.
- `rexo init <diretório>`: cria um projeto portátil com manifesto, memória em
  camadas, orçamento, política de qualidade e um workflow de exemplo executável.
- `rexo run <workflow.json>`: executa um workflow determinístico de ponta a
  ponta, gravando artefatos endereçados por conteúdo e um trace de execução;
  `--replay` verifica se a execução é reproduzível.
- A mesma base é testada em Windows, macOS e Linux.
- Os releases são executáveis únicos, sem exigir Python, Node, Docker ou uma
  conta de LLM.

## Por que ele é diferente

O REXO não começa por agentes com nomes e prompts permanentes. Ele começa por
**capabilities**. Um workflow pede um resultado; o runtime escolhe o provider
adequado considerando política, qualidade, custo, latência, disponibilidade e
cache.

## Teste local

Com Go 1.24 ou superior:

```shell
go test ./...
go build -o rexo ./cmd/rexo
./rexo doctor
./rexo init meu-primeiro-projeto
cd meu-primeiro-projeto
./rexo run workflow.json
```

No Windows, o executável gerado será `rexo.exe`.

## Princípios imutáveis do Core v1

1. Capability-first, não agent-first.
2. Contexto mínimo necessário para cada tarefa.
3. Workers descartáveis e artefatos duráveis.
4. Reutilizar e consultar cache antes de gerar.
5. Usar o menor modelo capaz de atingir a qualidade.
6. Orçamentos explícitos de tokens, chamadas, tempo e custo.
7. Qualidade mínima, tentativas limitadas e fallback.
8. Packs estendem a plataforma sem modificar o Core.
9. Execução observável, versionada e reproduzível.
10. A arquitetura só cresce para resolver problemas reais e recorrentes.

## O que ainda não existe

Providers de LLM e seleção de provider, MCP, Economy Engine completo, memória
semântica, Education Pack, marketplace, SDKs, Studio, Canvas e Creator fazem
parte das fases seguintes. O kernel atual é **determinístico** (sem IA/rede). O
repositório não fingirá que uma ideia no roadmap já é funcional.

## Licença e nome

Código sob Apache License 2.0. O nome REXO passou por uma triagem preliminar,
mas ainda precisará de busca jurídica formal antes de uso comercial.
