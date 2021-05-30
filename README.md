# G'Agent
[![Build Status](https://drone.dragonheim.net/api/badges/dragonheim/gagent/status.svg)](https://drone.dragonheim.net/dragonheim/gagent)
[![Go Report Card](https://goreportcard.com/badge/git.dragonheim.net/dragonheim/gagent)](https://goreportcard.com/report/git.dragonheim.net/dragonheim/gagent)

A Golang based mobile agent system loosely inspired by the [Agent Tcl / D'Agents](http://www.cs.dartmouth.edu/~dfk/agents/) system created by Robert S. Gray of Dartmouth college.

## Purpose
As we move close and closer to a true space-age, we need to start thinking about solutions for various space-age issues such as the bi-directional time delay between the surface of Mars and the surface of Earth. At present it takes between 6 and 44 minutes for single round-trip, making most online data services unuseable. G'Agent is a potential solution for data services given the time delay.

Imagine, for a moment, that you are on Mars and need to perform a data search in a specific domain space. You would have to explain it to someone on Earth, and hope they understand enough of the domain space to know where to search and understand you well enough to perform the actual search and then send you the results. With G'Agent, instead you would write a basic script (TCL), hereafter called an agent, providing various hints about the domain space and the search as TCL code. Your client would then send it on to a server, hereafter called a router, on Earth. The router may or may not know anything about the domain space of your search, so the router will use the hints that you provide to attempt to route the agent to known workers or other routers closer to the desired domain space. Eventually your agent will reach a router whose workers can handle your search.  The workers, will take the agent, run the script portion and collect the response(s), returning the reponse(s) to the router(s) for return to your client.

## Example Agent
```tcl
1  : ###################
2  : ### Hello Earth ###
3  : ###################
4  : # HINT START
5  : #  - thermal measurements
6  : #  - gravity measurements
7  : #  - gravity fluctuations
8  : # HINT END
9  : proc hello_earth {} {
10 :   puts "Hello Earth, does localized tempurature variations alter specific gravity?"
11 : }
12 : hello_earth
```
Lines 1 - 3 are simple comments

Line  4 indicates the start of the hints.

Lines 5 - 7 are a list of hints that the router will use to determine which router(s) may have domain specific information.

Line  8 indicates the end of the hints.

Lines 9 - 11 are a tcl procedure that will be executed on the worker before sending the results back to the client.

Line  12 executes the procedure defined above.

## History
More information about Agent TCL / D'Agents can be found in the original [documentation](http://www.cs.dartmouth.edu/~dfk/agents/pub/agents/doc.5.1.ps.gz), and in the project's [wiki](https://git.dragonheim.net/dragonheim/gagent/wiki/_pages).
