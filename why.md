# Por qué escribir pruebas unitarias y cómo hacer que trabajen para ti

[Aquí tienes un vídeo sobre mi hablando sobre este tema](https://www.youtube.com/watch?v=Kwtit8ZEK7U)

Si no te van los vídeos, aquí va la versión escrita.

## Software 

La promesa del software es que se puede cambiar. Por eso se llama _soft_ ware: es maleable comparado con el hardware. Un gran equipo de ingeniería ha de ser un activo increíble para una compañía, y desarrollar sistemas que puedan evolucionar con el negocio para seguir aportando valor.

¿Por qué se nos da tan mal entonces? ¿Cuántos proyectos conoces que hayan terminado en fracaso absoluto? O convertidos en "legacy" para ser reescritos completamente (¡y esas nuevas versiones a menudo fracasan también!)

¿Cómo puede fallar un sistema de software después de todo? ¿No se puede simplemente modificarlo hasta que funcione bien? ¡Esa era la promesa!

Mucha gente está eligiendo Go para desarrollar sistemas porque ha hecho una serie de elecciones que se espera lo hagan más resistente a convertirse en "legacy":

- A diferencia de mi vida previa de Scala donde [ya describí cómo tienes suficiente cuerda para ahorcarte tu sólo](http://www.quii.dev/Scala_-_Just_enough_rope_to_hang_yourself), Go cuenta únicamente con 25 palabras reservadas y se pueden construir _un montón_ de sistemas utilizando únicamente la librería estandar y unas pocas librerías más. La idea es que con Go puedas escribir código que al volver 6 meses más tarde siga teniendo sentido.
- Las herramientas relacionadas con testing, métricas de rendimiento, revisión de código y empaquetado son de primera categoría, comparadas con la mayoría de alternativas.
- La librería estándar es brillante.
- El tiempo de compilación es muy rápido, lo que permite bucles de feedback muy ajustados.
- La promesa de compatibilidad hacia atrás. Parece que Go recibirá genéricos y otras funcionalidades en el futuro pero los diseñadores han prometido que incluso el código Go escrito hace 5 años seguirá compilando. He pasado literalmente semanas migrando un proyecto de Scala 2.8 a 2.10

Incluso con todas estas grandes propiedades seguimos pudiendo escribir sistemas horribles, así que deberíamos mirar hacia el pasado para entender las lecciones de ingeniería de software que aplican independientemente de cómo de brillante sea (o no) tu lenguaje.

En 1974 un inteligente ingeniero de software llamado [Manny Lehman](https://en.wikipedia.org/wiki/Manny_Lehman_%28computer_scientist%29) escribió las [leyes de Lehman de la evolución del software](https://en.wikipedia.org/wiki/Lehman%27s_laws_of_software_evolution).

> Las leyes describen un equilibrio entre fuerzas que aceleran los nuevos desarrollos por un lado, y las fuerzas que ralentizan el progreso por el otro.

Parece importante entender estas fuerzas si queremos tener alguna esperanza de no terminar en un ciclo infinito de entrega de sistemas que se convierten en "legacy" y vuelven a ser re-escritos una y otra vez.


## La ley del Cambio Continuo

> Cualquier sistema de software usado en el mundo real debe cambiar, o perderá utilidad en el entorno.

Parece obvio que un sistema _tiene_ que cambiar para no volverse inútil, pero ¿cuántas veces ignoramos esto?

Muchos equipos son presionados para entregar un proyecto en una fecha particular y después son trasladados al siguiente proyecto. Si hay suerte, habrá al menos algún tipo de traspaso a otro conjunto de individuos para mantenerlo, pero por supuesto no serán los autores.

A menundo la gente se preocupa intentando elegir un framework que les ayude a "entregar rádido", sin preocuparse de la longevidad del sistema en términos de cómo necesitará evolucionar.

Aunque seas un ingeniero de software increíble, tu también serás víctima de no saber las necesidades futuras de tu sistema. A medida que el negocio cambie partes de tu briallante código se volverán irrelevantes.

Lehman estaba en racha en los 70 porque nos dio otra ley más para masticar.

## La ley de la Complejidad Incremental

> A medida que un sistema evoluciona, su complejidad aumenta a menos que se trabaje para reducirla.

Lo que nos está diciendo es que no podemos tener equipos de software como meras fábricas de funcionalidad, amontonando más y más funcionalidades en el software con la esperanza de que sobreviva a largo plazo.

**Tenemos** que seguir gestionando la complejidad del sistema a medida que el conocimiento sobre el dominio cambia.

## Refactorizar

Hay _muchas_ facetas de la ingeniería de software que lo mantienen maleable, como:

- Que los programadores tengan capacidad de decisión.
- Código "bueno" en general. Separación de responsabilidades, etc etc.
- Habilidades comunicativas.
- Arquitectura.
- Observabilidad.
- Desplegabilidad.
- Pruebas automatizadas.
- Bucles de feedback.

Me voy a centrar en la refactorización. Frecuentemente decimos "hay que refactorizar esto" a un programador en su primer día, sin más razonamiento.

¿De dónde viene esta frase? ¿Qué diferencia hay entre refactorizar y escribir código?

Sé que tanto yo como muchos otros _pensábamos_ que estábamos refactorizando, pero estábamos equivocados.

[Martin Fowler describe cómo la gente lo malinterpreta](https://martinfowler.com/bliki/RefactoringMalapropism.html)

> Sin embargo el término "refactorizar" se utiliza a menudo de forma inadecuada. Si alguien habla sobre un sistema que queda roto durante un par de días mientras lo están refactorizando, puedes estar seguro de que no están refactorizando.

¿Entonces qué es?

### Refactorización

Cuando estudiabas matemáticas en el colegio probablemente te enseñaron a factorizar. Veamos un ejemplo sencillo:
Calcular `1/2 + 1/4`

Para hacerlo _factorizamos_ los denominadores, convirtiedo la expresión en

`2/4 + 1/4` que finalmente nos da `3/4`

Podemos aprender importantes lecciones aquí.  _Factorizar la expresión_ no **cambia su significado**. Las dos son igualmente `3/4` pero hemos hecho que sea más sencillo de manejar para nosotros; cambiar `1/2` por `2/4` hace que encaje mejor en nuestro "dominio".

Cuando refactorizas tu código, intentas encontrar formas de hacerlo más sencillo de entender y que "encaje" en tu comprensión actual de lo que el sistema debe hacer. Es crucial **no modificar su comportamiento**.


#### Un ejemplo en Go

Aquí tenemos una función que saluda a `name` en un idioma `(language)` particular

    func Hello(name, language string) string {
    
      if language == "es" {
         return "Hola, " + name
      }
    
      if language == "fr" {
         return "Bonjour, " + name
      }
      
      // imagina docenas de idiomas más
    
      return "Hello, " + name
    }

Tener docenas de sentencias `if` no parece buena idea, y tenemos una duplicación en la concatenación del saludo específico del idioma con `, ` y el `name`. Así que voy a refactorizar el código.

    func Hello(name, language string) string {
      	return fmt.Sprintf(
      		"%s, %s",
      		greeting(language),
      		name,
      	)
    }
    
    var greetings = map[string]string {
      "es": "Hola",
      "fr": "Bonjour",
      //etc..
    }
    
    func greeting(language string) string {
      greeting, exists := greetings[language]
      
      if exists {
         return greeting
      }
      
      return "Hello"
    }

La naturaleza de esta refactorización no es importante, sino el hecho de que no he cambiado el comportamiento.

Al refactorizar puedes hacer lo que quieras: añadir interfaces, nuevos tipos, funciones, métodos, etc. La única regla es no cambiar el comportamiento.

### Al refactorizar código no debes modificar el comportamiento

Esto es muy importante. Si cambias el comportamiento al mismo tiempo que la estructura, estás haciendo _dos_ cosas a la vez. Como ingenieros de software aprendemos a partir los sistemas en archivos/paquetes/funciones/etc porque sabemos que intentar entender un muchas cosas a la vez es difícil.

No queremos tener que pensar sobre muchas cosas a la vez porque es entonces cuando cometemos errores. He visto fracasar muchos intentos de refactorización porque los programadores intentaban abarcar demasiado.

Cuando hacía factorizaciones en clase de matemáticas con papel y lápiz tenía que comprobar manualmente que no había cambiado el significado de las expresiones ¿Cómo sabemos que no estamos cambiando el comportamiento al refactorizar cuando trabajamos con código, especialmente en sistemas que no son triviales?

Quienes eligen no escribir tests típicamente confían en hacer pruebas manualmente. Para cualquier proyecto que no sea muy pequeño ésto es una tremenda pérdida de tiempo y no escala a largo plazo.
 
**Para refactorizar con seguridad necesitas tests** porque proporcionan

- Confianza de que puedes modificar el código sin cambiar el comportamiento
- Documentación para humanos sobre cómo debe comportarse el sistema
- Un feedback mucho más fiable y rápido que probar manualmente

#### Un ejemplo en Go

Un test unitario para nuestra función `Hello` tendría este aspecto:

    func TestHello(t *testing.T) {
      got := Hello(“Chris”, es)
      want := "Hola, Chris"
    
      if got != want {
         t.Errorf("got %q want %q", got, want)
      }
    }

Puedo ejecutar `go test` en la línea de comandos y obtener feedback inmediato sobre si mi refactorización ha alterado el comportamiento. En la práctica es mejor que aprender el botón mágico que hay que tocar para ejecutar los tests dentro de tu editor/IDE.

Lo que buscamos es entrar en un estado en el que:

- Hacemos una pequeña refactorización
- Ejecutamos los tests
- Repetimos

Todo en un bucle de feedback muy ajustado, evitando entrar por madrigueras de conejo y cometer errores.

Tener un proyecto en el que todos tus comportamientos tienen pruebas unitarias y te dan feedback en menos de un segundo es una red de seguridad que te habilita para refactorizar siempre que lo necesites. Ésto nos permitirá manejar la complejidad que vendrá, según describe Lehman.

## Si los tests unitarios son tan buenos ¿Por qué a veces hay resistencia a escribirlos?

Por una parte tienes a gente (como yo) diciendo que los tests unitarios son importantes para la salud a largo plazo de tu sistema porque aseguran que puedes seguir refactorizando con confianza.

Por otra, tienes a gente que describe experiencias en las que los tests unitarios han _impedido_ la refactorización

Pregúntate a ti mismo ¿con qué frecuencia tienes que cambiar los tests para refactorizar? A lo largo de los años he estado en muchos proyectos con muy buena coberetura de test en los que sin embargo los ingenieros eran reticentes a refactorizar por el esfuerzo percibido de modificar los tests.

¡Lo contrario a lo que nos prometieron!

### ¿Qué está pasando?

Imagina que te pidieran desarrollar un cuadrado, y que pensáramos que la mejor forma de conseguirlo es juntar dos triángulos.

![Dos triángulos rectángulos para formar un cuadrado](https://i.imgur.com/ela7SVf.jpg)

Escribimos nuestros tests unitarios sobre nuestro cuadrado para asegurarnos de que los lados son iguales y entonces escribimos algunos tests sobre nuestros triángulos. Queremos asegurarnos de que los triángulos se renderizan correctamente, así que introducimos una aserción de que los ángulos suman 180 grados, quizá comprobamos que se construyen 2, etc etc. La cobertura de tests es muy importante y escribir estos tests es fácil, así que ¿Por qué no?

Unas semanas más tarde la Ley del Cambio Continuo golpea nuestro sistema y una programadora nueva hace algunos cambios. A ella le parece que sería mejor si los cuadrados se formaran con dos rectángulos, en lugar de dos triángulos.

![Dos rectángulos para formar un cuadrado](https://i.imgur.com/1G6rYqD.jpg)

Al intentar hacer esta refactorización obtiene señales mezcladas de varios tests que fallan ¿Acaso ha roto algún comportamiento importante? Ahora tiene que revisar los tests de triángulos intentando entender qué está pasando.

_En realidad no es imporante si los cuadrados se forman a partir de triángulos_ pero **nuestros tests han elevado falsamente la importancia de nuestros detalles de implementación**.

## Probar el comportamiento, no la implementación

Cuando oigo a gente quejándose sobre tests unitarios a menudo es porque los tests están en un nivel de abstracción incorrecto. Están probando detalles de implementación, espiando de más en los colaboradores y usando demasiados dobles ("mocks").

En mi opinión ésto se debe a una malinterpretación de qué son los tests, y a perseguir métricas de vanidad (cobertura de tests).

Si estoy diciendo que deberíamos probar únicamente el contenido ¿no deberíamos escribir sólo tests de sistema/caja negra? Este tipo de tests tienen un gran valor para verificar experiencias de usuario, pero suelen ser difíciles de escribir y lentos de ejecutar. Por esa razon no son demasiado útiles para _refactorizar_ porque el bucle de feedback es lento. Además los tests de caja negra tienden a no ser de mucha ayuda para comprender las causas de fallo en comparación con los tesets unitarios.

¿_Cuál_ es el nivel de abstracción correcto entonces?

## Escribir tests unitarios eficaces es un problema de diseño

Olvidándonos de los tests por un momento, es deseable que tu sistema esté formado por "unidades" auto-contenidas y desacopladas, centradas alrededor de los conceptos principales de tu dominio.

Me gusta imaginar estas unidades como simples piezas de Lego que tienen APIs coherentes que puedo combinar con otras piezas para hacer sistemas más grandes. Bajo esas APIs puede haber docenas de cosas (tipos, funciones y demás) colaborando para hacerlas funcionar.

Por ejemplo, si estuvieras escribiendo un banco en Go, probablemente tendrías un paquete "account" (cuenta). Presentaría un API que no expusiera detalles de implementación y que fuese sencilla de integrar.

Si tienes estas unidades que siguen estas propiedades puedes escribir tests unitarios contra sus APIs públicas. _Por definición_ éstos test sólo pueden estar probando comportamiento útil. Por debajo, tengo libertad para refactorizar la implementación cuanto necesite y los tests prácticamente no deberían interferir.

### ¿Estos tests son unitarios?

**SI**. Los tests unitarios prueban "unidades" como las que he descrito. _Nunca_ prueban una única clase/función/lo que sea

## Agrupando conceptos

Hemos visto

- Refactorización
- Tests unitarios
- Diseño unitario

Lo que empezamos a vislumbrar es que éstas facetas del diseño de software se retroalimentan.

### Refactorización

- Nos da señales sobre nuestros tests unitarios. Si tenemos que hacer pruebas manuales, necesitamos más tests. Si los tests están fallando erróneamente entonces nuestros tests están en un nivel de abstracción equivocado (o no tienen valor y habría que borrarlos).
- Nos ayuda a gestionar la complejidad dentro de nuestras unidades, y entre ellas.

### Tests unitarios

- Nos dan una red de seguridad para refactorizar.
- Verifican y documentan el comportamiento de nuestras unidades.

### Unidades (bien diseñadas)

- Permiten escribir tests _relevantes_.
- Son fáciles de refactorizar.

¿Existe un proceso que nos ayude a llegar a un punto en el que poder refactorizar constantemente el código para gestionar la complejidad y mantener nuestros sistemas maleables?

## Por qué hacer Desarrollo Guiado por Pruebas (Test Driven Development, TDD)

Puede que algunos partan de las citas de Lehman sobre cómo el software debe cambiar y creen diseños demasiado elaborados, dedicando un montón de tiempo al principio para crear el sistema extensible "perfecto", y terminen haciéndolo mal y llegando a ninguna parte.

Así era en los malos viejos tiempos del software en los que un equipo de analistas se pasaba 6 meses escribiendo un documento de requisitos y un equipo de arquitectura otros 6 con el diseño para que unos años más tarde el proyecto fracasase.

¡Digo viejos malos tiempos, pero aún pasa!

El Desarrollo Ágil nos enseña que hay que trabajar de forma iterativa, comenzando por lo pequeño y haciendo evolucionar al software de forma que tengamos feedback rápidamente sobre el diseño de nuestro sistema y qué tal funciona para usuarios reales. El TDD refuerza este enfoque.

El TDD aborda las las Leyes de Lehman y otras lecciones aprendidas por las malas a lo largo de la historia, al promover una metodología de refactorización constante y producción iterativa.

### Pequeños pasos

- Escribe un pequeeño test para un pequeño comportamiento.
- Comprueba que el test falla con un error claro (rojo).
- Escribe la mínima cantidad de código para hacer que el test pase (verde)
- Refactoriza.
- Repite.

A medida que cojas práctica, esta forma de trabajar se volverá natural y rápida.

Llegarás a esperar que este bucle de feedback sea rápido y te sentirás incómodo si el sistema no está "verde", porque significa que te puedes haber metido por una madriguera de conejo.

Siempre estarás creando funcionalidades pequeñas y útiles confortablemente respaldadas por el feedback de tus tests.

## Resumiendo

- La fortaleza del software está en que lo podemos cambiar. La _mayoría_ del software requerirá cambios a lo largo del tiempo de formas impredecibles; pero no intentamos sobre-diseñar porque predecir el futuro es demasiado difícil
- En lugar de eso, necesitamos conseguir que nuesetro software siga siendo maleable. Para cambiar el software necesitamos refactorizarlo o se volverá un caos.
- Un buen conjunto de test puede ayudarnos a refactorizar más rápido y con más tranquilidad.
- Escribir buenos tests unitarios es un problema de diseño, así que piensa en organizar el código para que contenga unidades relevantes que puedas unir como piezas de Lego.
- El TDD puede ayudar y obligarte a diseñar código bien organizado iterativamente, respaldado por tests que ayuden con el trabajo futuro.
