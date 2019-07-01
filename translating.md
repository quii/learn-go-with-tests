# Diretivas de Tradução/Revisão

Requisitos: [Git](https://git-scm.com/downloads) instalado!

## Traduções x Revisões

É importante diferenciarmos durante todo o nosso processo quando estamos traduzindo ou estamos revisando.
As traduções acontecem diretamente relacionadas ao conteúdo original, escrito em inglês. Elas passaram por duas aprovações de outras pessoas envolvidas no projeto e já foram aceitas junto ao conteúdo oficial.
As revisões acontecem quando alguma melhoria pode ser feita na tradução.

## Começando a brincadeira

Na aba Projetos você pode encontrar todas as Issues criadas de acordo com cada tópico discutido no livro.
Tendo selecionado um tópico para ser traduzido/revisado:

-   Dê assign na issue para o seu apelido

-   Clone o repositório:

```bash
git clone https://github.com/larien/learn-go-with-tests.git
```

-   Crie uma branch seguindo o padrão `traducao/<nome-do-topico>` ou `revisao/<nome-do-topico>`:

```bash
git checkout -b traducao/<nome-do-topico>
```

-   Você pode salvar a sua tradução aos poucos com `commits`:

```bash
git commit -m "Pequena descrição do que foi traduzido aqui"
```

-   E não se esqueça de subir sua tradução para o repositório remoto!

```bash
git push -u origin traducao/<nome-do-topico>
```

E que comecem as traduções! :))

## Pontos para prestar atenção

-   Como pode ter percebido, não é um livro formal. Logo, não é necessário ter uma linguagem rebuscada: tente usar termos e palavras acessíveis para quem está conhecendo a linguagem.
-   Se for um termo muito técnico, não custa nada colocar uma breve explicação do que ele é ou linkar para uma referência externa.
-   Lembre-se de utilizar uma linguagem neutra! Se alguma parte do texto especificar algum gênero, reescreva-o para que inclua todo mundo.
-   Se houver dúvida em alguma palavra ou tradução, coloque um asterico (\*) que o pessoal que for revisar vai tentar te ajudar.
-   Quando for analisar a submissão de alguma tradução ou revisão, lembre-se de ter a gentileza em primeiro lugar! Todas as pessoas que colaborarem com o projeto querem ter um conteúdo de qualidade na nossa língua e, acima de tudo, querem aprender. Não se esqueça disso :)
-   Leia o [contributing.md](contributing.md) :)

## Submetendo sua tradução/revisão

Após fazer sua tradução para a sua própria branch e subi-la para o GitHub:

-   Vá na aba `Pull Requests` e clique em `New pull request`.

-   Redirecione a sua branch para a branch `master` do projeto `larien/learn-go-with-tests`, se isso já não estiver definido.

-   Na descrição, lembre-se de linkar a Issue referente ao tópico que você traduziu. Por exemplo, se o tópico que você selecionou é a Issue de número #14, digite na descrição do pull request `closes #14`. Isso vai automaticamente mover o card da sua Issue para `Done`.

-   Se for uma revisão, não é necessário linkar nenhuma issue.

-   Clique em `Create pull request` e aguarde a aprovação das outras pessoas envolvidas no projeto. Você receberá a notificação por e-mail quando sua PR for aprovada :)
