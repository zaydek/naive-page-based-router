# create-retro-app

Retro is a friendly development server and static-site generator (SSG) for React apps.

```

  Usage:

    retro create [dir]      Creates a new Retro app at directory dir
    retro watch [...paths]  Starts the dev server and watches paths for changes
    retro build             Builds the production-ready build
    retro serve             Serves the production-ready build

  retro create [dir]

    Creates a new Retro app at directory dir

      --template=[js|ts]    Starter template (defaults to js)

  retro watch [...paths]

    Starts a dev server and watches paths for changes (defaults to pages)

      --cached              Reuse cached props (disabled by default)
      --poll=<duration>     Poll duration (defaults to 250ms)
      --port=<number>       Port number (defaults to 8000)
      --source-map          Add source maps (disabled by default)

  retro build

    Builds the production-ready build

      --cached              Reuse cached props (disabled by default)
      --source-map          Add source maps (disabled by default)

  retro serve

    Serves the production-ready build

      --port=<number>       Port number (defaults to 8000)

  Repository:

    https://github.com/zaydek/retro

```
