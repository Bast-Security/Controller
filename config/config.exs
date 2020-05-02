import Config

config :bast, Bast.Repo,
  database: "bast",
  username: "bast",
  password: "bast",
  hostname: "localhost"

config :bast, ecto_repos: [Bast.Repo]
