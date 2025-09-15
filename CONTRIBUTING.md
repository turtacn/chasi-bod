# Contributing to Chasi-Bod

First off, thank you for considering contributing to Chasi-Bod! We welcome any and all contributions.

## Development Setup

To get started with developing Chasi-Bod, you'll need the following tools:

*   Go 1.18+
*   Docker
*   `make`

Once you have these tools installed, you can clone the repository and build the project:

```bash
git clone https://github.com/turtacn/chasi-bod.git
cd chasi-bod
make
```

## Running Tests

To run the unit tests, you can use the `test` target in the `Makefile`:

```bash
make test
```

## Pull Request Process

1.  **Fork the repository:** Create your own fork of the Chasi-Bod repository.
2.  **Create a branch:** Create a new branch for your changes.
    ```bash
    git checkout -b my-new-feature
    ```
3.  **Make your changes:** Make your changes to the code, and be sure to add or update tests as appropriate.
4.  **Run the tests:** Make sure that all tests pass before submitting your pull request.
    ```bash
    make test
    ```
5.  **Commit your changes:** Commit your changes with a clear and descriptive commit message.
    ```bash
    git commit -am 'Add some feature'
    ```
6.  **Push to your branch:** Push your changes to your fork.
    ```bash
    git push origin my-new-feature
    ```
7.  **Create a pull request:** Open a pull request from your fork to the `main` branch of the Chasi-Bod repository.

We will review your pull request as soon as possible. Thank you for your contribution!
