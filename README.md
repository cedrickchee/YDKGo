# You Don't Know Go Yet Book

[![Netlify Status](https://api.netlify.com/api/v1/badges/14232e97-1d7c-4ee4-ae09-da4fcd02d860/deploy-status)](https://app.netlify.com/sites/ydkgo/deploys)

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

#### License

<details>

<summary><b>Expand License</b></summary>

This repository contains a variety of content; some developed by Cedric Chee, and some from third-parties. The third-party content is distributed under the license provided by those parties.

The content of this project itself is licensed under the [Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License](http://creativecommons.org/licenses/by-nc-sa/4.0/), and the underlying source code used to format and display that content is licensed under the [Apache License, Version 2.0](LICENSE).
</details>
