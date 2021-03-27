# simpledash

Simple and fast dashboard with a go backend and a mostly HTML+CSS frontend with a small amount of JS.

### Building

Dependencies: 
- sqlite3 (and sqlite3-dev if on debian/ubuntu)
- golang (1.14+)

To build, simply run `go build` in the repo.

---

### Configuration

Configuration is done in the `simpledash.toml` file.

#### Root
The root section is the one not under any header.

It contains the following fields:

- `title`: Title used for whole site
- `theme`: Theme used for site, can be light or dark.
- `loginRequired`: Boolean denoting whether login is required to see public cards
- `allowProxy`: Array containing all sites allowed to be proxied by the integrated HTTP proxy

#### Users

A user can be added under the `[users]` section like so:
```toml
[users]
    [users.admin]
    passwordHash = "$2a$10$w00dzQ1PP6nwXLhuzV2pFOUU6m8bcZXtDX3UVxpOYq3fTSwVMqPge"
    showPublic = true
```
`passwordHash` should be a bcrypt hash of the desired password with a cost of 10 (default)
`showPublic` should be a boolean denoting whether public cards should be displayed while signed in

#### Cards

A card can either be public or belong to a user. A public card should be added under `[users]` like so:
```toml
[users]
    [[users._public_.card]]
        type = "weather"
        title = "Weather"
        data = {"woeid" = "2442047"}
```

A card belonging to a user should be added under that user, like so:
```toml
 [users.admin]
    passwordHash = "$2a$10$w00dzQ1PP6nwXLhuzV2pFOUU6m8bcZXtDX3UVxpOYq3fTSwVMqPge"
    showPublic = true

    [[users.admin.card]]
      type = "status"
      title = "Google"
      icon = "ion:logo-google"
      desc = "Google search engine. Status card example."
      url = "https://www.google.com"
```

The configuration for a card consists of up to six things:

- Type: Type of card to display
- Title: Title to show above card
- Icon: Icon to display on card
- Description: Description for card
- URL: URL to be used inside card
- Data: Extra data for anything not listed

Icons can be anything found on [Iconify](https://iconify.design)

There are currently five types of cards included:

- Simple: Simplest type of card, displays title, icon, description, and URL
- Status: Same as simple but also checks and displays status of URL
- Collection: A card containing multiple links to different sites. Links provided in data field.
- Weather: Display weather using data from Metaweather. Gets location via [WOEIDs](https://nations24.com/world-wide)
- API: Gets JSON data from URL and formats according to `format` field in data. At least part of URL must be inside `allowedProxy` array.

Examples for each card are included in `simpledash-sample.toml`

#### Session

The session cookie name can be set under `[session]` like so:

```toml
[session]
    name = "simpledash-session"
```