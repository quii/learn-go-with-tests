# Aprenda Go com Testes

<p align="center">
  <img src="red-green-blue-gophers-smaller.png" />
</p>

[Arte por Denise](https://twitter.com/deniseyu21)

![Build Status](https://travis-ci.org/quii/learn-go-with-tests.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/quii/learn-go-with-tests)](https://goreportcard.com/report/github.com/quii/learn-go-with-tests)

-   Formatos: [Gitbook](https://larien.gitbook.io/aprenda-go-com-testes), [EPUB or PDF](https://github.com/quii/learn-go-with-tests/releases)
-   Versão original: [English](https://quii.gitbook.io/learn-go-with-tests/)
-   Traduções: [中文](https://studygolang.gitbook.io/learn-go-with-tests)

## Motivação

-   Explore a linguagem Go escrevendo testes
-   **Tenha uma base com TDD**. O Go é uma boa linguagem para aprender TDD por ser simples de aprender e ter testes nativamente
-   Tenha confiança de que você será capaz de escrever sistemas robustos e bem testados em Go
-   [Assista a um vídeo ou leia sobre por quê testes unitários e TDD são importantes](why.md)

## Índice

### Primeiros Passos com Go

1. [Instale Go](install-go.md) - Prepare o ambiente para produtividade.
2. [Olá, mundo](hello-world.md) - Declarando variáveis, constantes, declarações `if`/`else`, switch, escreva seu primeiro programa em Go e seu primeiro teste. Sintaxe de subteste e closures.
3. [Inteiros](integers.md) - Mais conteúdo sobre sintaxe de declaração de função e aprenda novas formas de melhorar a documentação do seu código.
4. [Iteração](iteration.md) - Aprenda sobre `for` e benchmarking.
5. [Arrays e Slices](arrays-and-slices.md) - Aprenda sobre arrays, slices, `len`, varargs, `range` e cobertura de testes.
6. [Structs, Métodos e Interfaces](structs-methods-and-interfaces.md) - Aprenda sobre `structs`, métodos, `interface` e testes orientados a tabela (table driven tests).
7. [Ponteiros e Erros](pointers-and-errors.md) - Aprenda sobre ponteiros e erros.
8. [Maps](maps.md) - Aprenda sobre armazenamento de valores na estrutura de dados `map`.
9. [Injeção de Dependência](dependency-injection.md) - Aprenda sobre injeção de dependência, qual sua relação com interfaces e uma introdução a io.
10. [Mocking](mocking.md) - Use injeção de dependência com mocking para testar um código sem nenhum teste.
11. [Concorrência](concurrency.md) - Aprenda como escrever código concorrente para tornar seu software mais rápido.
12. [Select](select.md) - Aprenda a sincronizar processos assíncronos de forma elegante.
13. [Reflection](reflection.md) - Aprenda sobre reflection.
14. [Sync](sync.md) - Conheça algumas funcionalidades do pacote sync, como `WaitGroup` e `Mutex`.
15. [Context](context.md) - Use o pacote context para gerenciar e cancelar processos de longa duração.

### Crie uma aplicação

Agora que você já deu seus _Primeiros Passos com Go_, esperamos que você tenha uma base sólida das principais funcionalidades da linguagem e como TDD funciona.

Essa seção envolve a criação de uma aplicação.

Cada capítulo é uma continuação do antigo, expandindo as funcionalidades da aplicação conforme nosso Product Owner ditar.

Novos conceitos serão apresentados para ajudar a facilitar a escrever código de qualidade, mas a maior parte do material novo terá relação com o que pode ser feito com a biblioteca padrão do Go.

No final desse capítulo, você deverá ter uma boa ideia de como escrever uma aplicação em Go testada.

-   [Servidor HTTP](http-server.md) - Vamos criar uma aplicação que espera por requisições HTTP e responde a elas.
-   [JSON, Routing e Embedding](json.md) - Vamos fazer nossos endpoints retornarem JSON e explorar como trabalhar com rotas.
-   [IO e Classificação](io.md) - Vamos persistir e ler nossos dados do disco e falar sobre classificação de dados.
-   [Linha de Comando e Estrutura do Projeto](command-line.md) - Suportar diversas aplicações em uma base de código e ler entradas da linha de comando.
-   [Tempo](time.md) - Usar o pacote `time` para programar atividades.
-   [Websockets](websockets.md) - Aprenda a escrever e testar um servidor que usa websockets.

### Dúvidas e respostas

Costumo ver perguntas na Interwebs como:

> Como testo minha função incrível que faz x, y e z?

Se tiver esse tipo de pergunta, crie uma Issue no GitHub e vou tentar achar tempo para escrever um pequeno capítulo para resolver o problema. Acho que conteúdo como esse é valioso, já que está resolvendo problemas `reais` envolvendo testes que as pessoas têm.

-   [OS exec](os-exec.md) - Um exemplo de como podemos usar o sistema operacional para executar comandos para buscar dados e manter nossa lógica de negócio testável.
-   [Tipos de erro](error-types.md) - Exemplo de como criar seus próprios tipos de erro para melhorar seus testes e tornar seu código mais fácil de se trabalhar.

## Contribuição

-   _Esse projeto está em desenvolvimento_. Se tiver interesse em contribuir, por favor entre em contato.
-   Leia [contributing.md](https://github.com/quii/learn-go-with-tests/tree/842f4f24d1f1c20ba3bb23cbc376c7ca6f7ca79a/contributing.md) por diretrizes [TODO].
-   Tem ideias? Crie uma issue

## Explicação

Tenho experiência em apresentar Go a times de desenvolvimento e tenho testado abordagens diferentes sobre como evoluir um time de um grupo de pessoas que têm curiosidade sobre Go a criadores extremamente eficazes de sistemas em Go.

### O que não funcionou

#### Ler _o_ livro

Uma abordagem que tentamos foi pegar [o livro azul](https://www.amazon.com.br/Linguagem-Programa%C3%A7%C3%A3o-Go-Alan-Donovan/dp/8575225464) e toda semana discutir um capítulo junto de exercícios.

Amo esse livro, mas ele exige muito comprometimento. O livro é bem detalhado na explicação de conceitos, o que obviamente é ótimo, mas significa que o progresso é lento e uniforme - não é para todo mundo.

Descobri que apenas um pequeno número de pessoas pegaria o capítulo X para ler e faria os exercícios, enquanto que a maioria não.

#### Resolver alguns problemas

Katas são divertidos, mas geralmente se limitam ao escopo de aprender uma linguagem; é improvável que você use goroutines para resolver um kata.

Outro problema é quando você tem níveis diferentes de entusiasmo. Algumas pessoas aprendem mais da linguagem que outras e, quando demonstram o que já fizeram, confundem essas pessoas apresentando funcionalidades que as outras ainda não conhecem.

Isso acaba tornando o aprendizado bem _desestruturado_ e _específico_\*.

### O que funcionou

De longe, a forma mais eficaz foi apresentar os conceitos da linguagem aos poucos lendo o [go by example](https://gobyexample.com/), explorando-o com exemplos e discutindo-o como um grupo. Essa abordagem foi bem mais interativa do que "leia o capítulo X como lição de casa".

Com o tempo, a equipe ganhou uma base sólida da _gramátiica_ da linguagem para que conseguíssemos começar a desenvolver sistemas.

Para mim, isso é semelhante à ideia de praticar escalas quando se tenta aprender a tocar violão.

Não importa quão artístico você seja; é improvável que você crie músicas boas sem entender os fundamentos e praticando os mecanismos.

### O que funcionou para mim

Quando _eu_ aprendo uma nova linguagem de programação, costumo começar brincando em um REPL\*, mas hora ou outra preciso de mais estrutura.

O que eu gosto de fazer é explorar conceitos e então solidificar as ideias com testes. Testes certificam de que o código que escrevi está correto e documentam a funcionalidade que aprendi.

Usando minha experiência de aprendizado em grupo e a minha própria, vou tentar criar algo que seja útil para outras equipes. Aprender os conceitos escrevendo testes pequenos para que você possa usar suas habilidades de desenvolvimento de software e entregar sistemas ótimos.

## Para quem isso foi feito

-   Pessoas que se interessam em aprender Go.
-   Pessoas que já sabem Go, mas querem explorar testes com TDD.

## O que vamos precisar

-   Um computador!
-   [Go instalado](https://golang.org/)
-   Um editor de texto
-   Experiência com programação. Entendimento de conceitos como `if`, variáveis, funções etc.
-   Se sentir confortável com o terminal

## Feedback

-   Crie issues/submita PRs [aqui](https://github.com/quii/learn-go-with-tests) ou [me envie um tweet em @quii](https://twitter.com/quii).
-   Para a versão em português, submita um PR [aqui](https://github.com/larien/learn-go-with-tests) ou entre em contato comigo pelo [meu site](https://larien.dev).

[MIT license](LICENSE.md)

[Logo criado por egonelbre](https://github.com/egonelbre) Que estrela!
