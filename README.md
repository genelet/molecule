# Molecule

**Molecule** is a Go library designed to perform complex database operations across multiple tables as a single, atomic unit. It treats the database as a graph of "Atoms" (tables) and executes actions based on defined relationships.

[![GoDoc](https://godoc.org/github.com/genelet/molecule?status.svg)](https://godoc.org/github.com/genelet/molecule)

## Overview

Molecule abstracts database interactions into high-level actions, allowing you to:
- Define relationships between tables (Atoms).
- Execute complex, multi-table transactions (Molecules) via a simple API.
- Support RESTful-like operations (Insert, Update, Delete, Select/Topics) across the entire database graph.

## Packages

The project is organized into three main packages:

### 1. [godbi](./godbi)
The core engine of the library.
- Defines the `Molecule`, `Atom`, and `Action` structures.
- Handles the execution logic for database operations.
- **Key Features**: Stateless execution, `RunOption` configuration, and support for custom actions.

### 2. [rdb](./rdb)
Relational Database adapters and utilities.
- Provides tools to automatically generate a `Molecule` structure from an existing database schema.
- **Supported Databases**: PostgreSQL, MySQL, SQLite.
- Useful for bootstrapping a Molecule configuration from legacy databases.

### 3. [gometa](./gometa)
Metadata and Protobuf integration.
- Defines the `Graph` structure using Protocol Buffers.
- Provides utilities to convert between `Molecule` (Go structs) and `Graph` (Protobuf) representations.
- Enables serialization and transport of database schemas and action definitions.

## Installation

```bash
go get github.com/genelet/molecule
```

## Usage

Please refer to the README files in each subdirectory for detailed usage instructions and examples.
