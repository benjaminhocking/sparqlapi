# SPARQL-like API
A simple sparql-like api written in Go.
Users can build their own directed graph to explain relationships between things. Two examples are provided: a friends and family dataset, and the ODRL information model.

# How To Use
## Setup
Clone the repo and either:
 - run on your local machine at localhost
 - deploy to an aws instance (t2.micro recommended)
## Functionality
Build your own dataset by creating new triples of the syntax:
> \<subject\> \<predicate\> \<object\> .

## Roadmap:
[x] Generate and Update .ttl file given new triple
[x] SPARQL Querying:
    [x] SELECT
    [ ] CONSTRUCT
    [ ] ASK
    [ ] DESCRIBE
[ ] SPARQL Update:
    [ ] Single Operation
    [ ] Multi-Operation
[x] Build directed graph with node and edge labels
[x] Host on aws
[ ] Daemonize hosting service with file support
[x] Demo Datasets:
    [x] Friends & Family
    [x] ODRL information model
[ ] Improve site design
[ ] Create API endpoints
[ ] Host central site allowing users to spin off their own sparql request handlers

There's still a lot of features to add, but this is where I managed to get in one day.