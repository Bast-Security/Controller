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
        scheme: :https,
        plug: Bast.Api,
        options: [
          port: 8080,
          otp_app: :bast,
          cipher_suite: :strong,
          keyfile: "priv/pki/bast.key",
          certfile: "priv/pki/bast.crt"
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
