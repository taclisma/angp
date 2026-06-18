# ANGp — Plano de Teste

## 1. Descrição do Sistema

### 1.1 Visão Geral
ANGp é uma aplicação desktop multiplataforma desenvolvida em Go que permite ao usuário desenhar com caneta e borracha, escolher tamanhos de pincel e persistir desenhos em arquivos .png.

### 1.2 Funcionalidades Principais
- Desenhar pontos/traços livres com ferramenta caneta (clique ou clique e arrasto)
- Apagar traços com ferramenta borracha
- Selecionar tamanho do pincel (S/M/L)
- Limpar a tela
- Salvar desenho em arquivo (formato .png)
- Carregar desenho de arquivo

### 1.3 Funcionalidades Planejadas (não implementadas)
- Desenhar linhas (clique e arrasto)
- Desenhar retângulos (clique e arrasto)
- Desenhar círculos (clique e arrasto)
- Selecionar cor de uma paleta pré-definida

### 1.4 Plataformas Suportadas
- Linux
- Windows
- macOS

---

## 2. Tecnologias Utilizadas

| Componente | Tecnologia | Versão | Justificativa |
|---|---|---|---|
| **Linguagem** | Go | 1.26+ | Requisito do projeto; compilação multiplataforma nativa |
| **Framework GUI** | Fyne | v2.7+ | Suporta multiplataforma; simples; aplicação desktop nativa |
| **Persistência** | PNG | - | Salva imagem gerada |
| **Análise Estática** | golangci-lint | v2.11+ | Detector padrão para Go; identifica erros comuns, estilo e segurança |
| **Segurança** | gosec | - | Detecta padrões inseguros no código (via golangci-lint) |
| **Vulnerabilidades** | govulncheck | latest | Verifica CVEs em dependências usando a base oficial do Go |
| **Dependabot** | GitHub | - | Alerta e PRs automáticos para dependências vulneráveis |
| **Testes Unitários** | Go testing (stdlib) + testify | 1.26+ | Nativo + assertions legíveis |
| **Testes de GUI** | fyne.io/fyne/v2/test | v2.7+ | Simulação de interações em memória; roda em CI |
| **CI/CD** | GitHub Actions | - | Gratuito; integrado ao repositório; suporta múltiplos SOs |
| **Release** | GoReleaser | latest | Cross-compilation, checksums, GitHub Releases |

---

## 3. Ferramentas de Teste

> **Convenção de nomenclatura de testes:**
> - Testes unitários: `TestUnit_<Funcionalidade>`
> - Testes de integração: `TestIntegration_<Fluxo>`
> - Testes de sistema: `TestSystem_<Cenario>`
>
> Filtrar por nível: `go test ./... -run TestUnit`, `go test ./... -run TestIntegration`, `go test ./... -run TestSystem`.

### 3.1 Verificação Estática
```bash
golangci-lint run ./...    # lint + gosec
govulncheck ./...          # vulnerabilidades em dependências
```

### 3.2 Testes Unitários
- **Ferramenta**: Pacote `testing` + `github.com/stretchr/testify`
- **Estrutura**: Arquivos `*_test.go` no mesmo diretório do código
- **Comando**: `go test ./... -v`
- **Cobertura**: `go test ./... -cover`

### 3.3 Testes de Integração
- **Ferramenta:** Testes Go com múltiplos pacotes
- **Foco:** Fluxo completo (desenhar → salvar → carregar → validar pixels)
- **Comando:** `go test ./... -v -run Integration`

### 3.4 Testes de Sistema (GUI automatizado)
- **Ferramenta:** Pacote `fyne.io/fyne/v2/test`
- **Foco:** Fluxo de usuário completo simulado em memória
- **Capacidades:**
  - `test.NewWindow(content)` — janela virtual em memória
  - `test.Tap(obj)` / `test.TapAt(obj, pos)` — simula clique
  - `test.Drag(c, pos, deltaX, deltaY)` — simula arrasto
  - `test.AssertRendersToImage(t, "master.png", c)` — golden file comparison
- **Comando:** `go test ./... -v -run System`
- **Status:** Não implementado

---

## 4. Procedimentos

### 4.1 Estrutura de Commits
Padrão Conventional Commits:
- `feat:` nova funcionalidade
- `fix:` correção de bug
- `test:` testes
- `refactor:` reorganização
- `docs:` documentação

### 4.2 Nomenclatura de Branches
- `feat/nome-feature`
- `fix/descricao-bug`
- `test/descricao-teste`
- `refactor/descricao`
- `docs/descricao`

### 4.3 Fluxo de Pull Request
1. Criar branch a partir de `main`
2. Commits com mensagens claras
3. Rodar localmente: `make test && make lint`
4. Push + abrir PR
5. CI/CD passa
6. Merge

---

## 5. Requisitos e Configuração

### 5.1 Ambiente
- **Go:** 1.26+
- **SO:** Linux, Windows ou macOS

### 5.2 Setup Local
```bash
git clone <repo>
go mod download
make test    # roda testes
make lint    # roda linter + gosec
make vuln    # verifica vulnerabilidades
make dev     # roda o app
```

### 5.3 Makefile
| Target | Comando |
|---|---|
| `make dev` | `go run ./cmd/angp` |
| `make test` | `go test ./... -v` |
| `make test-unit` | `go test ./... -v -run TestUnit` |
| `make test-integration` | `go test ./... -v -run TestIntegration` |
| `make cover` | `go test ./... -cover` |
| `make lint` | `golangci-lint run ./...` |
| `make vuln` | `govulncheck ./...` |
| `make build` | `go build -o angp ./cmd/angp` |
| `make release-dry` | `goreleaser release --snapshot --clean` |

### 5.4 Cobertura Mínima
- Meta: **70%**
- Atual: **95.5%**

---

## 6. Matriz de Funcionalidades vs Testes

> **Pirâmide de testes:** Maioria dos cenários coberta por testes unitários (rápidos e isolados), integração verifica interação entre componentes, sistema valida fluxo completo do usuário via GUI.

### 6.1 Funcionalidades (caminho feliz)

| # | Funcionalidade | Teste Unitário | Teste Integração | Teste Sistema | Status |
|---|---|---|---|---|---|
| 1 | Desenhar linha | Validar pontos início/fim; rejeitar pontos coincidentes | Desenhar → encode PNG → verificar pixels | `Drag` + `AssertRendersToImage` | ⚠️ Feature não implementada |
| 2 | Desenhar retângulo | Calcular largura/altura; validar bounds | Desenhar → encode PNG → verificar contorno | `Drag` + `AssertRendersToImage` | ⚠️ Feature não implementada |
| 3 | Desenhar círculo | Calcular centro e raio; validar raio > 0 | Desenhar → encode PNG → verificar circunferência | `Drag` + `AssertRendersToImage` | ⚠️ Feature não implementada |
| 4 | Selecionar cor | Verificar cor padrão aplicada; verificar cor selecionada | Selecionar cor → desenhar → verificar pixels | `Tap` paleta + `Drag` + `AssertRendersToImage` | ⚠️ Feature não implementada (color picker) |
| 5 | Limpar tela | Resetar estado interno (lista de strokes vazia) | Desenhar → limpar → encode PNG → verificar branco | `Tap` limpar + `AssertRendersToImage` | ✅ `TestUnit_Clear` + `TestIntegration_ClearThenSave` |
| 6 | Salvar desenho | Encode PNG sem erro; verificar dimensões | Desenhar → salvar → ler arquivo → comparar pixels | `Tap` salvar + verificar arquivo | ✅ `TestUnit_Save` + `TestIntegration_DrawSaveLoad` |
| 7 | Carregar desenho | Decodificar PNG válido; validar dimensões | Salvar → carregar → comparar estado | `Tap` carregar + `AssertRendersToImage` | ✅ `TestUnit_Load` + `TestIntegration_DrawSaveLoad` |

### 6.2 Cenários de erro e borda

| # | Cenário | Teste Unitário | Teste Integração | Teste Sistema | Status |
|---|---|---|---|---|---|
| 8 | Carregar PNG inválido/corrompido | Retornar erro; não causar panic | Verificar canvas mantém estado anterior | — | ✅ `TestUnit_Load_Invalid` + `TestIntegration_LoadInvalid_PreservesState` |
| 9 | Clique sem arrasto (ferramenta de forma) | Forma não é criada; lista não muda | — | — | ⚠️ Feature não implementada (formas) |
| 10 | Clique sem arrasto (caneta) | Ponto é criado na posição do clique | — | — | ✅ `TestUnit_SingleClick` |
| 11 | Desenhar fora dos limites do canvas | Coordenadas clampadas aos limites | — | — | ✅ `TestUnit_Clamp` |
| 12 | Salvar sem ter desenhado nada | Gerar PNG vazio válido | PNG com dimensões corretas e pixels brancos | — | ✅ `TestUnit_Save_EmptyCanvas` |
| 13 | Limpar tela já vazia | Estado não muda; sem erro | — | — | ✅ `TestUnit_Clear_AlreadyEmpty` |
| 14 | Salvar em caminho sem permissão | Retornar erro; não causar panic | Canvas mantém estado | — | ✅ `TestUnit_Save_BadWriter` + `TestIntegration_SaveBadWriter_PreservesState` |
| 15 | Carregar arquivo não-PNG | Retornar erro; não causar panic | Canvas mantém estado anterior | — | ✅ `TestUnit_Load_Invalid` + `TestIntegration_LoadNonPNG_PreservesState` |

### 6.3 Resumo de Cobertura

| Nível | Implementados | Pendentes |
|---|---|---|
| Unitários | 11 testes | #1-4, #9 (features não implementadas) |
| Integração | 5 testes | — |
| Sistema (GUI) | 0 | Todos (fyne/test não integrado ainda) |

---

## 7. Configuração do Ambiente CI/CD

### 7.1 GitHub Actions (`.github/workflows/go.yml`)
- **Triggers:** Push em `main`; Pull Requests

| Job | Runner | Função |
|---|---|---|
| `lint` | ubuntu-latest | golangci-lint + gosec |
| `vuln` | ubuntu-latest | govulncheck |
| `test` | ubuntu-latest | build + testes + cobertura |
| `cross-build` | ubuntu, windows, macos | build em cada SO (verifica compatibilidade) |

### 7.2 Release (`.github/workflows/release.yml`)
- **Trigger:** Manual (workflow_dispatch)
- **Input:** Versão (ex: `v1.0.0`)
- **Saída:** GitHub Release com binários para linux/windows/darwin (amd64) + checksums

### 7.3 Dependabot (`.github/dependabot.yml`)
- Verifica dependências Go semanalmente
- Abre PRs automáticos para atualizações de segurança

---

## 8. Estrutura do Projeto

```
.
├── cmd/angp/main.go              # entrypoint (GUI)
├── internal/canvas/
│   ├── canvas.go                 # domínio (pen, eraser, save, load)
│   └── canvas_test.go            # testes unitários e integração
├── .github/
│   ├── workflows/go.yml          # CI (lint, vuln, test, cross-build)
│   ├── workflows/release.yml     # Release (manual, GoReleaser)
│   └── dependabot.yml            # Alertas de segurança
├── .golangci.yml                 # Configuração do linter
├── .goreleaser.yml               # Configuração de release
├── go.mod / go.sum               # Dependências
└── makefile                      # Comandos de desenvolvimento
```

---

## 9. Referências

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Go Image/PNG Package](https://golang.org/pkg/image/png/)
- [Fyne Documentation](https://pkg.go.dev/fyne.io/fyne/v2)
- [Fyne Test Package](https://pkg.go.dev/fyne.io/fyne/v2/test)
- [testify](https://github.com/stretchr/testify)
- [golangci-lint](https://golangci-lint.run/)
- [gosec](https://github.com/securego/gosec)
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [GoReleaser](https://goreleaser.com/)
- [GitHub Actions](https://github.com/features/actions)
- [Conventional Commits](https://www.conventionalcommits.org/)
