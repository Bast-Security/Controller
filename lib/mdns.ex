defmodule Bast.Controller.Mdns do
  use GenServer

  defp add_services do
    Mdns.Server.start

    {:ok, [ {ip, _a, _b} | _rest ] } = :inet.getif()

    Mdns.Server.set_ip(ip)

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
  end

  def init(state) do
    add_services()
    {:ok, state}
  end

  def start_link(opts) do
    GenServer.start_link(__MODULE__, [], opts)
  end
end
