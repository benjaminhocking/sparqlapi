# SPARQL-like API
A simple sparql-like api written in Go.
Users can build their own directed graph to explain relationships between things. Two examples are provided: a friends and family dataset, and the ODRL information model.

# How To Use
## Setup
Clone the rep and either:
 - run on your local machine at localhost
 - deploy to an aws instance (t2.micro recommended)
## Functionality
Build your own dataset by creating new records of the syntax:
    <subject> <predicate> <object> .