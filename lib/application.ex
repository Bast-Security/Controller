defmodule Bast.Middleman.Application do
  use Application
  require Logger

  def start(_type, _args) do
    children = [
      Bast.Middleman.Repo,

      {
        Plug.Cowboy,
        scheme: :http,
        plug: Bast.Middleman.REST,
        options: [port: 8080]
      }
    ]

    options = [
      strategy: :one_for_one,
      name: Bast.Middleman.Supervisor
    ]

    Logger.info("Starting Middleman.")
    Supervisor.start_link(children, options)
  end
end

