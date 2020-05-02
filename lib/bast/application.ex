defmodule Bast.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  def start(_type, _args) do
    children = [
      # Starts a worker by calling: Bast.Worker.start_link(arg)
      # {Bast.Worker, arg}
      
      {
        Plug.Cowboy,
        scheme: :http,
        plug: Bast.Api,
        options: [
          port: 8080
        ]
      },

      Bast.Repo
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Bast.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
