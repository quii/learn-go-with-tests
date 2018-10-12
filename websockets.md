# Websockets (WIP)

In this chapter we'll learn how to use websockets to improve our application. 

## Project recap

We have two applications in our poker codebase

- *Command line app*. Prompts the user to enter the number of players in a game. From then on informs the players of what the "blind bet" value is, which increases over time. At any point a user can enter `"{Playername} wins"` to finish the game and record the victor in a store.
- *Web app*. Allows users to record winners of games and displays a league table. Shares the same store as the command line app. 

## Next steps

The product owner is thrilled with the command line application but would prefer it if we could bring that functionality to the browser. She imagines a web page with a text box that allows the user to enter the number of players and when they submit the form the page displays the blind value and automatically updates it when appropriate. Like the command line application the user can declare the winner and it'll get saved in the database.

On the face of it, it sounds quite simple but as always we must emphasise taking an _iterative_ approach to writing software. Let's break the problem down.

- We can start by rendering a page with a form to enter the number of players. We'll verify that the page starts a game by using a spy. 
- Once we're happy with that we can think about how to wire up the `BlindAlerter` so that it somehow writes it to the web page (we'll use websockets)
- Finally we can add another form (once the game starts) to let the user declare a winner



