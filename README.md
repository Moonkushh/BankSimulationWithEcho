# bankSimulationEcho
## Meaning of the project:
The point was to move the functionality of https://github.com/Matterlinkk/goosestasks to the echo web framework
## Project composition
<ol>
  <li>The main structure is Client, identifier, account in the form of a structure consisting of a balance and mutex value, and a channel that collects values of the Transaction structure type</li>
  <li>PerformTransactions function for depositing and withdrawing funds</li>
  <li>ProcessTransactions function for monitoring transactions</li>
  <li>SendFunds function for transferring funds to another user</li>
  <li>Handlers for this functionality, and also routes, have been brought into one function that combines them together</li>
  <li>Docs folder for generating swagger documentation</li>
</ol>
<i>The code is run from the main.go file</i>

# How to
## Run the project
1. Install golang (Go to [official website](https://go.dev) for the help)
2. Install git (Go to [git install instructions](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) for the help)
3. Clone the repository: `git clone https://github.com/Matterlinkk/bankSimulationEcho`
4. In project folder, execute: `go get`
5. Build project, execute: `go build main.go`

<i>This will start the server on localhost:1488, I recommend using swagger's online documentation for easier use</i>
<br>
https://editor.swagger.io/