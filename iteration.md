# Iteração

**[Você pode encontrar todo o código desse capitulo aqui](https://github.com/quii/learn-go-with-tests/tree/master/for)**

Para fazer coisas repetidamente em Go, você precisara do `for`. Go não possui nenhuma palavra chave do tipo `while`, `do`, `until`, você pode usar apenas `for`, o que é uma coisa boa!

Vamos escrever um teste para uma função que repete um caracter 5 vezes.

Não há nenhuma novidade até aqui, então tente escrever você mesmo para praticar.

## Escreva o teste primeiro

```go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a")
    expected := "aaaaa"

    if repeated != expected {
        t.Errorf("expected '%s' but got '%s'", expected, repeated)
    }
}
```

## Tente e execute o teste

`./repeat_test.go:6:14: undefined: Repeat`

## Escreva a quantidade minima de código para o teste rodar e verifique a saida de errors

_Mantenha a disciplina!_ Você não precisa saber nada de diferente para fazer o teste falhar apropriadamente.

Tudo que você fez até agora é o suficiente para compilar, verifique se você escreveu o teste corretamente.

```go
package iteration

func Repeat(character string) string {
    return ""
}
```

Não é legal saber que você já sabe o bastante em Go para escrever testes para problemas simples? Isso significa que agora você pode mexer no código de produção o quanto quiser sabendo que ele se comportara da maneira que você desejar.

`repeat_test.go:10: expected 'aaaaa' but got ''`

## Escreva código o suficiente para fazer passar

A sintaxe do `for` é muito facil de memorar e segue quase igual as linguagens baseadas em `C`

```go
func Repeat(character string) string {
    var repeated string
    for i := 0; i < 5; i++ {
        repeated = repeated + character
    }
    return repeated
}
```

Ao contrario de outras linguagens como C, Java ou Javascript, não há parenteses ao redor dos três componentes do `for` mas as chaves `{ }` são obrigatórias

Execute o teste e ele devera passar.

Variação adicionais do loop `for` estão descritas [aqui](https://gobyexample.com/for).

## Refatoração

Agora é hora de refatorarmos e apresentarmos outro operador de atribuição o `+=`

```go
const repeatCount = 5

func Repeat(character string) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

O operador adicionar & atribuir `+=` adiciona o valor que esta na direita no valor que esta a esquerda e atribui o resultado ao valor da esquerda, também funciona com outros tipos como por exemplo inteiros (integer).

### Benchmarking

Escrever [benchmarks](https://golang.org/pkg/testing/#hdr-Benchmarks) em Go é outro recurso disponivel diretamente na linguagem e é tão facil quanto escrever tests.

```go
func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a")
    }
}
```

Você ira notar que o código é muito parecido com um teste.

O `testing.B` dara a você acesso a `b.N`

When the benchmark code is executed, it runs `b.N` times and measures how long it takes.

Quando o benchmark é executado, ele executara `b.N` vezes e medira quanto tempo levou.

A quantidade de vezes que o código não deve importar para você, o framework irá determinar qual é um valor "bom" para que você consiga ter resultados decentes.

Para executar o benchmark faça `go test -bench=.` (ou se você estiver executando do PowerShell do Windows `go test-bench="."`)

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           136 ns/op
PASS
```

`136 ns/op` significa que nossa função demora cerca de 136 nano segundos para ser executada \(no meu computador). E isso é ótimo! para chegar a esse resultado ela foi executada 10000000 vezes.

_NOTE_ por padrão o benchmark é executado sequencialmente.

## Exercicios para praticar

* Altere o teste para que a função aceite quantas vezes o caracter deve ser repetido.
* Escreve `ExampleRepeat` para documentar sua função
* Veja também o pacote [strings](https://golang.org/pkg/strings). Encontre funções que você acha serem uteis e experimente elas escrevendo testes como nós fizemos aqui. Investir tempo aprendendo os pacotes nativos irá te recompensar com o tempo.

## Resumindo

* Mais praticas de TDD
* Aprendemos `for`
* Aprendemos como escrever benchmarks
