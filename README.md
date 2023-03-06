# G'Agent
[![Maintained Status](https://img.shields.io/maintenance/yes/2023?style=plastic)](https://github.com/dragonheim/gagent)
[![License](https://img.shields.io/badge/License-MIT-limegreen.svg)](https://github.com/dragonheim/gagent/src/branch/main/LICENSE)

[![Build Status](https://drone.dragonheim.net/api/badges/dragonheim/gagent/status.svg)](https://drone.dragonheim.net/dragonheim/gagent)
[![Go Report Card](https://goreportcard.com/badge/github.com/dragonheim/gagent)](https://goreportcard.com/report/github.com/dragonheim/gagent)
[![Docker Pulls](https://img.shields.io/docker/pulls/dragonheim/gagent)](https://hub.docker.com/r/dragonheim/gagent/tags?page=1&ordering=last_updated)

A Golang based mobile agent system loosely inspired by the [Agent Tcl / D'Agents](http://www.cs.dartmouth.edu/~dfk/agents/) system created by Robert S. Gray of Dartmouth college.

## Purpose
As we move close and closer to a true space-age, we need to start thinking about solutions for various space-age issues such as the bi-directional time delay between the surface of Mars and the surface of Earth. At present it takes between 6 and 44 minutes for single round-trip, making most online data services unuseable. G'Agent is a potential solution for data services given the time delay.

Imagine, for a moment, that you are on Mars and need to perform a data search in a specific domain space. You would have to explain it to someone on Earth, and hope they understand enough of the domain space to know where to search and understand you well enough to perform the actual search and then send you the results. With G'Agent, instead you would write a basic script (TCL), hereafter called an agent, providing various hints about the domain space and the search as TCL code. Your client would then send it on to a server, hereafter called a router, on Earth. The router may or may not know anything about the domain space of your search, so the router will use the hints that you provide to attempt to route the agent to known workers or other routers closer to the desired domain space. Eventually your agent will reach a router whose workers can handle your search.  The workers, will take the agent, run the script portion and collect the response(s), returning the reponse(s) to the router(s) for return to your client.

## Example Agent
```tcl
1  : ###################
2  : ### Hello Earth ###
3  : ###################
4  : set GHINT [split "thermal measurements, gravity measurements, gravity fluctuations" ,]
5  : proc hello_earth {} {
6  :   puts "Hello Earth, does localized tempurature variations alter specific gravity?"
7  : }
8  : hello_earth
```
Lines 1 - 3 are simple comments.

Line  4 sets the hint(s) as an array.

Lines 5 - 7 defines a procedure called hello_earth that will perform the actual search

Line  8 executes the hello_earth procedure defined above.


## History
More information about Agent TCL / D'Agents can be found in the original [documentation](http://www.cs.dartmouth.edu/~dfk/agents/pub/agents/doc.5.1.ps.gz), and in the project's [wiki](https://github.com/dragonheim/gagent/wiki/_pages).
