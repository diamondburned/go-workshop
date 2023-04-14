# bookstore

```sh
git clone https://github.com/diamondburned/acm-go-present
cd acm-go-present/cmd/bookstore
```

## Running

```sh
go run .
```

## Libraries, not Frameworks

To help us write certain parts of the code, we'll be using the following
libraries:

- [chi](https://github.com/go-chi/chi) for routing. It only implements the
  router and nothing else.
- [assert](https://github.com/alecthomas/assert) for testing. It simply
  provides functions for checking if any values are equal.

When writing Go, we tend to **avoid frameworks**. The difference between
libraries and frameworks is that:

- Libraries are just a set of functions that you import as you need them, while
- Frameworks are a set of functions that come in bulk, most of which you don't
  need.

Because Go has a complete standard library, we rarely need to use frameworks
to make up for the lack of features. Instead, we use libraries to fill in the
gaps.

**Note**: you don't need to manually install these libraries. Simply run the
program normally, and it will automatically download the libraries for you.
