defmodule Bast.Api do
  use Plug.Router

  alias Bast.Api.Status

  plug Plug.Logger
  plug :match
  plug Plug.Parsers, parsers: [:json], json_decoder: Poison
  plug :dispatch

  forward "/admins", to: Bast.Api.Admins
  forward "/locks", to: Bast.Api.Locks
  forward "/systems", to: Bast.Api.Systems

  match _ do
    Status.send_status(Status.not_found, conn)
  end
end
