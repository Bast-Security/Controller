defmodule Bast.Api do
  use Plug.Router

  plug Plug.Logger
  plug :match
  plug Plug.Parsers, parsers: [:json], json_decoder: Poison
  plug :dispatch

  get "/" do
    conn |> send_resp(200, "Hello World!")
  end
end
