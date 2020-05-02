defmodule Bast.Api.Systems.Roles do
  use Plug.Router

  alias Bast.Api.Status

  plug :match
  plug :dispatch

  match _ do
    Status.send_status(Status.not_found, conn)
  end
end
