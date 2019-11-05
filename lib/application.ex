defmodule Bast.Controller.Application do
  use Application
  require Logger

  def start(_type, _args) do
    children = [
      Bast.Controller.Repo,

      {
        Plug.Cowboy,
        scheme: :http,
        plug: Bast.Controller.REST,
        options: [port: 8080]
      }
    ]

    options = [
      strategy: :one_for_one,
      name: Bast.Controller.Supervisor
    ]

    Logger.info("Starting Controller.")
    Supervisor.start_link(children, options)
  end
end

