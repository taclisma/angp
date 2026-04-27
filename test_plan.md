# Plano de Teste - ANGp

## 1. Descrição do Sistema

### 1.1 Visão Geral
ANGp é uma aplicação desktop multiplataforma desenvolvida em Go que permite ao usuário desenhar formas geométricas, escolher cores e persistir desenhos em arquivos em .png.

### 1.2 Funcionalidades Principais
- Desenhar linhas (clique e arrasto)
- Desenhar retângulos (clique e arrasto)
- Desenhar círculos (clique e arrasto)
- Desenhar pontos/traços livres com ferramenta caneta (clique ou clique e arrasto)
- Selecionar cor de uma paleta pré-definida
- Limpar a tela
- Salvar desenho em arquivo (formato .png)
- Carregar desenho de arquivo

### 1.3 Plataformas Suportadas
- Linux
- Windows
- macOS

---

## 2. Tecnologias Utilizadas

| Componente | Tecnologia | Versão | Justificativa |
|---|---|---|---|
| **Linguagem** | Go | 1.21+ | Requisito do projeto; compilação multiplataforma nativa |
| **Framework GUI** | Fyne | v2.4+ | Suporta multiplataforma; simples; aplicação desktop nativa |
| **Persistência** | PNG | - | Salva imagem gerada, falta de features que justifiquem salvar com um formato para projeto neste ponto |
| **Análise Estática** | golangci-lint | 1.54+ | Detector padrão para Go; identifica erros comuns e estilo |
| **Testes Unitários** | Go testing (stdlib) | 1.21+ | Nativo, sem dependências externas |
| **Testes de GUI** | fyne.io/fyne/v2/test | v2.4+ | Simulação de interações (tap, drag, type) em memória; golden file assertions; roda em CI |
| **CI/CD** | GitHub Actions | - | Gratuito; integrado ao repositório; suporta múltiplos SOs |

---

## 3. Ferramentas de Teste

> **Convenção de nomenclatura de testes:** Para permitir a execução seletiva por nível, todos os testes devem seguir o padrão de prefixo:
> - Testes unitários: `TestUnit_<Funcionalidade>` (ex.: `TestUnit_DesenharLinha`)
> - Testes de integração: `TestIntegration_<Fluxo>` (ex.: `TestIntegration_DesenharSalvarCarregar`)
> - Testes de sistema: `TestSystem_<Cenario>` (ex.: `TestSystem_FluxoDesenhoCompleto`)
>
> Isso permite filtrar por nível com `-run`: `go test ./... -run TestUnit`, `go test ./... -run TestIntegration`, `go test ./... -run TestSystem`.

### 3.1 Verificação Estática
- **golangci-lint:** Análise de código (linting, segurança, estilo)
- **Comandos:**
  ```bash
  golangci-lint run ./...
  ```

### 3.2 Testes Unitários
- **Ferramenta**: Pacote `testing` padrão do Go
- **Estrutura**: Arquivos `*_test.go` no mesmo diretório do código
- **Comando**: `go test ./... -v`
- **Cobertura**: `go test ./... -cover`

### 3.3 Testes de Integração
- **Ferramenta:** Testes Go com múltiplos pacotes
- **Foco:** Fluxo completo (desenhar → salvar → carregar → validar)
- **Comando:** `go test ./... -v -run Integration`

### 3.4 Testes de Sistema (GUI automatizado)
- **Ferramenta:** Pacote `fyne.io/fyne/v2/test`
- **Foco:** Fluxo de usuário completo simulado em memória (sem necessidade de display)
- **Capacidades utilizadas:**
  - `test.NewWindow(content)` — cria janela virtual em memória para teste
  - `test.Tap(obj)` / `test.TapAt(obj, pos)` — simula clique em botões e widgets
  - `test.TapCanvas(c, pos)` — simula clique em posição absoluta no canvas
  - `test.Drag(c, pos, deltaX, deltaY)` — simula arrasto (desenho de formas)
  - `test.MoveMouse(c, pos)` — simula movimentação do mouse
  - `test.Type(obj, chars)` — simula entrada de texto
  - `test.AssertRendersToImage(t, "master.png", c)` — comparação com golden file (imagem de referência)
- **Golden Files:** Imagens de referência armazenadas em `testdata/`. Na primeira execução, geram-se os masters; nas seguintes, compara-se pixel a pixel. Falhas geram arquivos em `testdata/failed/` para inspeção.
- **Manutenção de Golden Files:**
  - Para regenerar imagens de referência após mudança legítima na UI, executar com a flag de atualização: `go test ./... -run System -update` (ou deletar os arquivos master existentes e re-executar os testes).
  - Diferenças de anti-aliasing entre plataformas e versões do Go podem causar falhas falsas. Recomenda-se gerar os golden files no mesmo ambiente do CI (ex.: `ubuntu-latest`) e não committar golden files gerados localmente em macOS/Windows sem validação.
  - Ao revisar PRs que alteram golden files, sempre inspecionar visualmente os diffs de imagem.
- **Comando:** `go test ./... -v -run System`
- **Execução em CI:** Roda sem display gráfico, pois o pacote `fyne/test` renderiza tudo via software em memória

---

## 4. Procedimentos

### 4.1 Estrutura de Commits
Seguir padrão Conventional Commits:
- `feat: adiciona funcionalidade`
- `fix: corrige bug`
- `test: adiciona testes`
- `refactor: reorganiza código`
- `docs: atualiza documentação`

**Exemplo:**
```
feat: implementa desenho de círculos
fix: corrige validação de cor hexadecimal
test: adiciona testes para renderização PNG
```

### 4.2 Nomenclatura de Branches
Seguir o mesmo padrão dos commits para nomes de branch:
- `feat/nome-feature` — nova funcionalidade
- `fix/descricao-bug` — correção de bug
- `test/descricao-teste` — adição ou alteração de testes
- `refactor/descricao` — refatoração
- `docs/descricao` — documentação

### 4.3 Fluxo de Pull Request
1. Criar branch a partir de `main`: `git checkout -b feat/nome-feature`
2. Fazer commits com mensagens claras
3. Rodar testes localmente: `go test ./... -v`
4. Rodar análise estática: `golangci-lint run ./...`
5. Fazer push: `git push origin feat/nome-feature`
6. Abrir PR no GitHub com descrição clara
7. Aguardar CI/CD passar
8. Merge após aprovação

---

## 5. Requisitos, Restrições e Configurações para Teste

### 5.1 Ambiente de Desenvolvimento
- **Go:** 1.21 ou superior
- **Sistema Operacional:** Linux, Windows ou macOS
- **RAM mínima:** 2 GB
- **Espaço em disco:** 500 MB (incluindo dependências)

### 5.2 Dependências
```bash
go mod download
go mod tidy
```

### 5.3 Configuração Local
1. Clonar repositório: `git clone <repo>`
2. Instalar dependências: `go mod download`
3. Instalar golangci-lint:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```
4. Rodar verificação estática: `golangci-lint run ./...`
5. Rodar testes: `go test ./... -v -cover`

### 5.4 Cobertura Mínima
- Meta de cobertura: **70%**
- Verificação no CI via script que analisa a saída de `go test -coverprofile`:
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out
  ```
  Script de enforcement no pipeline (exemplo):
  ```bash
  COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
  if [ "$(echo "$COVERAGE < 70.0" | bc)" -eq 1 ]; then
    echo "FALHA: Cobertura total ($COVERAGE%) abaixo do mínimo de 70%"
    exit 1
  fi
  ```
  O pipeline falha se a cobertura total ficar abaixo de 70%.

---

## 6. Matriz de Funcionalidades vs Testes

> **Nota:** A distribuição dos testes segue o princípio da pirâmide de testes: a maioria dos cenários é coberta por testes unitários (rápidos e isolados), testes de integração verificam a interação entre componentes, e testes de sistema via `fyne/test` validam o fluxo completo do usuário. Cenários de erro e borda são testados predominantemente no nível unitário, subindo para integração apenas quando envolvem múltiplos componentes.

> **Prioridade de implementação:** Recomenda-se implementar os testes na seguinte ordem: (1) cenários unitários de erro e borda (#8–#15), pois são rápidos de escrever e protegem contra regressões básicas; (2) cenários unitários do caminho feliz (#1–#7); (3) testes de integração; (4) testes de sistema com `fyne/test`. Dentro de cada nível, priorizar salvar/carregar (#6, #7, #8, #14, #15) por serem as funcionalidades com maior risco de perda de dados.

### 6.1 Funcionalidades (caminho feliz)

| # | Funcionalidade | Teste Unitário | Teste Integração | Teste Sistema (`fyne/test`) | Status |
|---|---|---|---|---|---|
| 1 | Desenhar linha | ✅ Validar pontos início/fim; rejeitar pontos coincidentes | ✅ Desenhar linha → encode PNG → decodificar e verificar pixels na posição esperada | ✅ `Drag` no canvas + `AssertRendersToImage` | TODO |
| 2 | Desenhar retângulo | ✅ Calcular largura/altura a partir de dois pontos; validar bounds | ✅ Desenhar retângulo → encode PNG → verificar dimensões e pixels do contorno | ✅ `Drag` no canvas + `AssertRendersToImage` | TODO |
| 3 | Desenhar círculo | ✅ Calcular centro e raio a partir de dois pontos; validar raio > 0 | ✅ Desenhar círculo → encode PNG → verificar pixels na circunferência | ✅ `Drag` no canvas + `AssertRendersToImage` | TODO |
| 4 | Selecionar cor | ✅ Verificar que cor padrão é aplicada; verificar cor selecionada dentro da paleta | ✅ Selecionar cor → desenhar forma → verificar pixels com a cor esperada | ✅ `Tap` na paleta + `Drag` + `AssertRendersToImage` | TODO |
| 5 | Limpar tela | ✅ Resetar estado interno (lista de formas vazia) | ✅ Desenhar → limpar → encode PNG → verificar imagem em branco | ✅ `Drag` + `Tap` limpar + `AssertRendersToImage` | TODO |
| 6 | Salvar desenho | ✅ Encode PNG sem erro; verificar dimensões do arquivo gerado | ✅ Desenhar → salvar → ler arquivo do disco → decodificar PNG → comparar pixels | ✅ `Drag` + `Tap` salvar + verificar arquivo PNG no disco | TODO |
| 7 | Carregar desenho | ✅ Decodificar PNG válido sem erro; validar dimensões da imagem | ✅ Salvar PNG → carregar → comparar estado do canvas com o original | ✅ `Tap` carregar + `AssertRendersToImage` | TODO |

### 6.2 Cenários de erro e borda

| # | Cenário | Teste Unitário | Teste Integração | Teste Sistema | Status |
|---|---|---|---|---|---|
| 8 | Carregar PNG inválido/corrompido | ✅ Retornar erro descritivo; não causar panic | ✅ Verificar que canvas mantém estado anterior | — | TODO |
| 9 | Clique sem arrasto (ferramenta de forma) | ✅ Forma não é criada; lista de formas não muda | — | — | TODO |
| 10 | Clique sem arrasto (ferramenta de ponto/caneta) | ✅ Ponto é criado na posição do clique | — | — | TODO |
| 11 | Desenhar fora dos limites do canvas | ✅ Coordenadas são clampadas aos limites do canvas | — | — | TODO |
| 12 | Salvar sem ter desenhado nada | ✅ Gerar PNG vazio válido sem erro | ✅ Arquivo PNG gerado com dimensões corretas e pixels em branco | — | TODO |
| 13 | Limpar tela já vazia | ✅ Estado não muda; operação não retorna erro | — | — | TODO |
| 14 | Salvar em caminho sem permissão ou inválido | ✅ Retornar erro descritivo (ex.: permissão negada); não causar panic | ✅ Verificar que canvas mantém estado e nenhum arquivo é criado | — | TODO |
| 15 | Carregar arquivo não-PNG (ex.: .jpg, .txt, .gif) | ✅ Retornar erro descritivo; não causar panic | ✅ Verificar que canvas mantém estado anterior após tentativa | — | TODO |

---

## 7. Configuração do Ambiente CI/CD

### 7.1 GitHub Actions
- **Arquivo:** `.github/workflows/test.yml`
- **Triggers:** Push em `main`; Pull Requests
- **Matriz de SOs:** O pipeline deve executar em `ubuntu-latest`, `windows-latest` e `macos-latest` via `strategy.matrix.os` para validar o suporte multiplataforma declarado na Seção 1.3.
- **Jobs:**
  1. Lint (golangci-lint)
  2. Test (go test)
  3. Coverage Report

### 7.2 Passos do Pipeline
1. Checkout do código
2. Setup Go 1.21+
3. Download de dependências
4. Executar golangci-lint
5. Executar testes unitários
6. Gerar relatório de cobertura (`go test -coverprofile`)
7. Falhar se cobertura < 70%

---

## 8. Referências

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Go Image Package](https://golang.org/pkg/image/)
- [Go Image/PNG Package](https://golang.org/pkg/image/png/)
- [Fyne Documentation](https://pkg.go.dev/fyne.io/fyne/v2)
- [Fyne Test Package](https://pkg.go.dev/fyne.io/fyne/v2/test)
- [Fyne - Testing Graphical Apps](https://docs.fyne.io/started/testing/)
- [Fyne - Golden File Testing](https://docs.fyne.io/started/testing/#golden-files)
- [golangci-lint](https://golangci-lint.run/)
- [GitHub Actions](https://github.com/features/actions)
- [Conventional Commits](https://www.conventionalcommits.org/)
