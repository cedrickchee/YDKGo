# You Don't Know Go Yet Book

Go programming language study notes turned book.

This book is inspired by ['You Don't Know JS Yet' (YDKJS)](https://github.com/getify/You-Dont-Know-JS) book series. YDKJS helped me understand JavaScript under the hood, after more than 8 years writing software with JS.

You Don't Know Go (YDKGo) book was created based on [Ultimate Go training](https://www.ardanlabs.com/ultimate-go/), which is an intermediate-level class for engineers with some experience with Go trying to dig deeper into the language.

## Development

The site is built using the [Hugo](https://gohugo.io/) static site generator. Check out the [Hugo Quick Start](https://gohugo.io/getting-started/quick-start/) for a quick intro.

### Requirements

- Hugo 0.68 and above
- Hugo extended version (with Sass/SCSS support)

### Serving the site locally

To build and locally serve this site, you need to [install Hugo, extended version](https://gohugo.io/getting-started/installing). Once Hugo is installed:

```sh
make serve
```

### Publishing the site

The website is automatically published by [Netlify](https://netlify.com/). Any time changes are pushed to the `master` branch, the site is rebuilt and redeployed.

### Site content

All of the [Markdown](https://www.markdownguide.org/) content used to build the site's documentation is in the [`content`](content) directory.
