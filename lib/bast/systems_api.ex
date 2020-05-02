defmodule Bast.Api.Systems do
  use Plug.Router

  plug :match
  plug :dispatch

  match _ do
    conn |> send_resp(404, "Not found")
  end
end
