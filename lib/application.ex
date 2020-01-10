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

    Logger.info("Registering mDNS services...")

    Mdns.Server.start

    Mdns.Server.set_ip({10, 0, 1, 29})

    [
      %Mdns.Server.Service{
        domain: "bast.local",
        data: :ip,
        ttl: 450,
        type: :a
      },

      %Mdns.Server.Service{
        domain: "_services._dns-sd._udp.local",
        data: "_bast._tcp.local",
        ttl: 4500,
        type: :ptr
      },

      %Mdns.Server.Service{
        domain: "_bast._tcp.local",
        data: "Bast._bast._tcp.local",
        ttl: 4500,
        type: :ptr
      }
    ] |> Enum.each(&Mdns.Server.add_service/1)

    Supervisor.start_link(children, options)
  end
end

